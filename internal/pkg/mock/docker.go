package mock

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/puppetlabs/prm/pkg/prm"
)

type DockerClient struct {
	Platform     string
	Version      string
	ApiVersion   string
	ErrorString  string
	ImagesSlice  []types.ImageSummary
	ExitCode     int64
	ExitErrorMsg string
}

type ReadClose struct{}

func (re *ReadClose) Read(r []byte) (n int, err error) {
	return 0, nil
}

func (re *ReadClose) Close() error {
	return nil
}

func (m *DockerClient) ServerVersion(ctx context.Context) (types.Version, error) {
	if m.ErrorString != "" {
		return types.Version{}, fmt.Errorf(m.ErrorString)
	}
	versionInfo := &types.Version{
		Platform:   struct{ Name string }{m.Platform},
		Version:    m.Version,
		APIVersion: m.ApiVersion,
	}
	return *versionInfo, nil
}

func (m *DockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
	return container.ContainerCreateCreatedBody{}, nil
}

func getSrcBuffer(stdOutBytes, stdErrBytes []byte) (buffer *bytes.Buffer, err error) {
	buffer = new(bytes.Buffer)
	dstOut := stdcopy.NewStdWriter(buffer, stdcopy.Stdout)
	_, err = dstOut.Write(stdOutBytes)
	if err != nil {
		return
	}
	dstErr := stdcopy.NewStdWriter(buffer, stdcopy.Stderr)
	_, err = dstErr.Write(stdErrBytes)
	return
}

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() (err error) {
	//we don't actually have to do anything here, since the buffer is just some data in memory
	//and the error is initialized to no-error
	return
}

func (m *DockerClient) ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	stdOutBytes := []byte("This is a test")
	stdErrBytes := []byte("")
	buffer, err := getSrcBuffer(stdOutBytes, stdErrBytes)
	if err != nil {
		return nil, err
	}

	closingBuffer := &ClosingBuffer{buffer}

	return closingBuffer, nil
}

func (m *DockerClient) ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error {
	return nil
}

func (m *DockerClient) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	return nil
}

func (m *DockerClient) ContainerWait(ctx context.Context, containerID string, condition container.WaitCondition) (<-chan container.ContainerWaitOKBody, <-chan error) {
	waitChan := make(chan container.ContainerWaitOKBody)
	go func() {
		waitChan <- container.ContainerWaitOKBody{StatusCode: m.ExitCode, Error: &container.ContainerWaitOKBodyError{Message: m.ExitErrorMsg}}
	}()
	errChan := make(chan error)
	return waitChan, errChan
}

func (m *DockerClient) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	return types.ImageBuildResponse{Body: &ReadClose{}}, nil
}

func (m *DockerClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	return m.ImagesSlice, nil
}

func (m *DockerClient) ImageRemove(ctx context.Context, imageID string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	return []types.ImageDeleteResponseItem{{Deleted: "test_id"}}, nil
}

func (m *DockerClient) ImageName(tool *prm.Tool, prmConfig prm.Config) string {
	// build up a name based on the tool and puppet version
	imageName := fmt.Sprintf("pdk:puppet-%s_%s-%s_%s", prmConfig.PuppetVersion.String(), tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id, tool.Cfg.Plugin.Version)
	return imageName
}
