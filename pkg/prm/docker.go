package prm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	dockerClient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/rs/zerolog/log"
)

type Docker struct {
	// We need to be able to mock the docker client in testing
	Client     DockerClientI
	OrigClient *dockerClient.Client
	Context    context.Context
}

type DockerClientI interface {
	// All docker client functions must be noted here so they can be mocked
	ServerVersion(context.Context) (types.Version, error)
}

func (d *Docker) GetTool(tool *Tool, prmConfig Config) error {

	// initialise the docker client
	err := d.initClient()
	if err != nil {
		return err
	}

	// what are we looking for?
	toolImageName := d.ImageName(tool, prmConfig)

	// find out if docker knows about our tool
	list, err := d.OrigClient.ImageList(d.Context, types.ImageListOptions{})

	if err != nil {
		log.Debug().Msgf("Error listing containers: %v", err)
		return err
	}

	for _, image := range list {
		for _, tag := range image.RepoTags {
			if tag == toolImageName {
				log.Info().Msgf("Found container: %s", image.ID)
				return nil
			}
		}
	}

	// No image found with that configuration
	// we must create it
	// d.createDockerfile(tool, prmConfig)

	return fmt.Errorf("No image found %s", toolImageName)
}

// Creates a unique name for the image based on the tool and the PRM configuration
func (d *Docker) ImageName(tool *Tool, prmConfig Config) string {
	// build up a name based on the tool and puppet version
	imageName := fmt.Sprintf("pdk:puppet-%s_%s-%s_%s", prmConfig.PuppetVersion.String(), tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id, tool.Cfg.Plugin.Version)
	return imageName
}

func (*Docker) Validate(tool *Tool) (ToolExitCode, error) {
	// TODO
	return FAILURE, nil
}

func (d *Docker) Exec(tool *Tool, args []string, prmConfig Config, paths DirectoryPaths) (ToolExitCode, error) {
	// is Docker up and running?
	status := d.Status()
	if !status.IsAvailable {
		log.Error().Msgf("Docker is not available")
		return FAILURE, fmt.Errorf("%s", status.StatusMsg)
	}

	// clean up paths
	codeDir, _ := filepath.Abs(paths.codeDir)
	log.Info().Msgf("Code path: %s", codeDir)
	cacheDir, _ := filepath.Abs(paths.cacheDir)
	log.Info().Msgf("Cache path: %s", cacheDir)

	// stand up a container
	resp, err := d.OrigClient.ContainerCreate(d.Context, &container.Config{
		Image: d.ImageName(tool, prmConfig),
		Cmd:   args,
		Tty:   false,
	},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: codeDir,
					Target: "/code",
				},
				{
					Type:   mount.TypeBind,
					Source: cacheDir,
					Target: "/cache",
				},
			},
		}, nil, nil, "")

	if err != nil {
		return FAILURE, err
	}
	defer d.OrigClient.ContainerRemove(d.Context, resp.ID, types.ContainerRemoveOptions{})

	if err := d.OrigClient.ContainerStart(d.Context, resp.ID, types.ContainerStartOptions{}); err != nil {
		return FAILURE, err
	}

	statusCh, errCh := d.OrigClient.ContainerWait(d.Context, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return FAILURE, err
		}
	case <-statusCh:
	}

	out, err := d.OrigClient.ContainerLogs(d.Context, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return FAILURE, err
	}

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	if err != nil {
		return FAILURE, err
	}

	return SUCCESS, nil
}

func (d *Docker) initClient() (err error) {
	if d.Client == nil {
		cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv)

        if err != nil {
            return err
        }

		d.Client = cli
		d.OrigClient = cli // TODO: remove this when we know all the functions that need added to the interface
		d.Context = context.Background()
	}
	return nil
}

// Check to see if the Docker runtime is available:
// if so, return true and info about Docker on this node;
// if not, return false and the error message
func (d *Docker) Status() BackendStatus {
	err := d.initClient()
	if err != nil {
		return BackendStatus{
			IsAvailable: false,
			StatusMsg:   fmt.Sprintf("unable to initialize the docker client: %s", err.Error()),
		}
	}
	// The client does not error on creation if the background service is not running,
	// but attempting to list the containers does.
	dockerInfo, err := d.Client.ServerVersion(d.Context)
	if err != nil {
		// message := fmt.Sprintf("%s", err)
		message := err.Error()
		// This is 90% likely the reason this command fails;
		// the alternative error message is lengthy and includes
		// references to pipes and the API which are more likely
		// to confuse than help; so trim it to the most useful info.
		daemonNotRunning := "error during connect: This error may indicate that the docker daemon is not running."
		if strings.Contains(message, daemonNotRunning) {
			message = daemonNotRunning
		}
		return BackendStatus{
			IsAvailable: false,
			StatusMsg:   message,
		}
	}
	status := fmt.Sprintf("\tPlatform: %s\n\tVersion: %s\n\tAPI Version: %s", dockerInfo.Platform.Name, dockerInfo.Version, dockerInfo.APIVersion)
	return BackendStatus{
		IsAvailable: true,
		StatusMsg:   status,
	}
}
