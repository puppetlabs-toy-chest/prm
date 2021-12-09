package prm

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	dockerClient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
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
		log.Debug().Msgf("Error listing images: %v", err)
		return err
	}

	for _, image := range list {
		for _, tag := range image.RepoTags {
			if tag == toolImageName {
				log.Info().Msgf("Found image: %s", image.ID)
				return nil
			}
		}
	}

	log.Info().Msg("Creating new image. Please wait...")

	// No image found with that configuration
	// we must create it
	fileString := d.createDockerfile(tool, prmConfig)
	log.Debug().Msgf("Creating Dockerfile\n--------------------\n%s--------------------\n", fileString)
	reader := strings.NewReader(fileString)

	// write the contents of fileString to a Dockerfile stored in the
	// tool path
	filePath := filepath.Join(tool.Cfg.Path, "generated.Dockerfile")
	file, err := os.Create(filePath)
	if err != nil {
		log.Error().Msgf("Error creating Dockerfile: %v", err)
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Msgf("Error closing file: %s", err)
		}
	}()

	_, err = io.Copy(file, reader)
	if err != nil {
		log.Error().Msgf("Error copying Dockerfile: %v", err)
		return err
	}

	// create a tar of the tool directory *shrug*
	tar, err := archive.TarWithOptions(tool.Cfg.Path, &archive.TarOptions{})
	if err != nil {
		return err
	}

	// build the image
	imageBuildResponse, err := d.OrigClient.ImageBuild(
		d.Context,
		tar,
		types.ImageBuildOptions{
			Dockerfile: "generated.Dockerfile",
			Tags:       []string{toolImageName},
			Remove:     true,
		})

	if err != nil {
		log.Error().Msgf("Unable to build docker image")
		return err
	}

	defer imageBuildResponse.Body.Close()

	// Parse the output from Docker, cleaning up where possible
	scanner := bufio.NewScanner(imageBuildResponse.Body)
	for scanner.Scan() {
		var line map[string]string
		_ = json.Unmarshal(scanner.Bytes(), &line) // nolint:errcheck // we don't care about the error here
		printLine := strings.TrimSuffix(line["stream"], "\n")
		if printLine != "" {
			log.Debug().Msgf("%s", printLine)
		}
	}

	return nil
}

func (d *Docker) createDockerfile(tool *Tool, prmConfig Config) string {
	// create a dockerfile from the Tool and prmConfig
	dockerfile := strings.Builder{}
	dockerfile.WriteString(fmt.Sprintf("FROM puppet/puppet-agent:%s\n", prmConfig.PuppetVersion.String()))

	if tool.Cfg.Common.RequiresGit || (tool.Cfg.Gem != nil && tool.Cfg.Gem.BuildTools) {
		dockerfile.WriteString("RUN apt update\n")
	}

	if tool.Cfg.Common.RequiresGit {
		dockerfile.WriteString("RUN apt install git -y\n")
	}

	if tool.Cfg.Gem != nil {
		if tool.Cfg.Gem.BuildTools {
			dockerfile.WriteString("RUN apt install build-essential -y\n")
		}

		dockerfile.WriteString("RUN /opt/puppetlabs/puppet/bin/gem install bundler --no-document\n")

		for _, gem := range tool.Cfg.Gem.Name {
			dockerfile.WriteString(fmt.Sprintf("RUN /opt/puppetlabs/puppet/bin/gem install %s -f --conservative --minimal-deps --no-document\n", gem))
		}
	}

	for key, val := range tool.Cfg.Common.Env {
		dockerfile.WriteString(fmt.Sprintf("ENV %s=%s\n", key, val))
	}

	// Copy the tools content into the image
	// contentPath := filepath.Join(tool.Cfg.Path, "content", "*")
	// dockerfile.WriteString(fmt.Sprintf("COPY %s /tmp/ \n", contentPath))
	if _, err := os.Stat(filepath.Join(tool.Cfg.Path, "/content")); err == nil {
		dockerfile.WriteString("COPY ./content/* /tmp/ \n")
	}

	dockerfile.WriteString("VOLUME [ /code, /cache ]\n")
	dockerfile.WriteString("WORKDIR /code\n")

	if tool.Cfg.Common.UseScript != "" {
		// todo: handle ps1 scripts
		dockerfile.WriteString(fmt.Sprintf("ENTRYPOINT [\"/tmp/%s.sh\"]\n", tool.Cfg.Common.UseScript))
	} else {
		if tool.Cfg.Gem != nil {
			dockerfile.WriteString(fmt.Sprintf("ENTRYPOINT [ \"/opt/puppetlabs/puppet/bin/%s\"]\n", tool.Cfg.Gem.Executable))
		}
	}

	if len(tool.Cfg.Common.DefaultArgs) > 0 {
		dockerfile.WriteString(fmt.Sprintf("CMD %q\n", tool.Cfg.Common.DefaultArgs))
	}

	return dockerfile.String()
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

	log.Info().Msgf("Additional Args: %v", args)

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
	defer func() {
		err := d.OrigClient.ContainerRemove(d.Context, resp.ID, types.ContainerRemoveOptions{})
		if err != nil {
			log.Error().Msgf("Error removing container: %s", err)
		}
	}()

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

	out, err := d.OrigClient.ContainerLogs(d.Context, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
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
