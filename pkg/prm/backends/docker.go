package backends

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/puppetlabs/prm/pkg/prm"
)

type Docker struct {
	// We need to be able to mock the docker client in testing
	Client DockerClient
}

type DockerClient interface {
	// All docker client functions must be noted here so they can be mocked
	ServerVersion(context.Context) (types.Version, error)
}

func (*Docker) GetTool(toolName string, prmConfig prm.Config) (prm.Tool, error) {
	// TODO
	return prm.Tool{}, nil
}

func (*Docker) Validate(tool *prm.Tool) (prm.ToolExitCode, error) {
	// TODO
	return prm.FAILURE, nil
}

func (*Docker) Exec(tool *prm.Tool, args []string) (prm.ToolExitCode, error) {
	// TODO
	return prm.FAILURE, nil
}

// Check to see if the Docker runtime is available:
// if so, return true and info about Docker on this node;
// if not, return false and the error message
func (d *Docker) Status() prm.BackendStatus {
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
		return prm.BackendStatus{
			IsAvailable: false,
			StatusMsg:   message,
		}
	}
	status := fmt.Sprintf("\tPlatform: %s\n\tVersion: %s\n\tAPI Version: %s", dockerInfo.Platform.Name, dockerInfo.Version, dockerInfo.APIVersion)
	return prm.BackendStatus{
		IsAvailable: true,
		StatusMsg:   status,
	}
}
