package prm

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
)

type Docker struct {
	// We need to be able to mock the docker client in testing
	Client DockerClientI
}

type DockerClientI interface {
	// All docker client functions must be noted here so they can be mocked
	ServerVersion(context.Context) (types.Version, error)
}

func (*Docker) GetTool(toolName string, prmConfig Config) (Tool, error) {
	// TODO
	return Tool{}, nil
}

func (*Docker) Validate(tool *Tool) (ToolExitCode, error) {
	// TODO
	return FAILURE, nil
}

func (d *Docker) Exec(tool *Tool, args []string) (ToolExitCode, error) {

	d.initClient()

	log.Info().Msgf("Executing docker exec command")
	log.Info().Msgf("Tool: %v", tool.Cfg.Plugin)

	if tool.Cfg.Gem != nil {
		log.Info().Msgf("GEM")
	}

	if tool.Cfg.Puppet != nil {
		log.Info().Msgf("PUPPET")
	}

	if tool.Cfg.Binary != nil {
		log.Info().Msgf("BINARY")
	}

	if tool.Cfg.Container != nil {
		log.Info().Msgf("CONTAINER")
	}

	// TODO
	return FAILURE, nil
}

func (d *Docker) initClient() (err error) {
	if d.Client == nil {
		cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv)

		d.Client = cli
		return err
	}
	return err
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
	dockerInfo, err := d.Client.ServerVersion(context.Background())
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
