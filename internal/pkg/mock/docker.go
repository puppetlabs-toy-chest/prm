package mock

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

type DockerClient struct {
	Platform    string
	Version     string
	ApiVersion  string
	ErrorString string
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

func (m *DockerClient) ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error) {

	mockReader := strings.NewReader("FAKE LOG MESSAGES!")
	mockReadCloser := io.NopCloser(mockReader)

	return mockReadCloser, nil
}

func (m *DockerClient) ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error {
	return nil
}

func (m *DockerClient) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	return nil
}

func (m *DockerClient) ContainerWait(ctx context.Context, containerID string, condition container.WaitCondition) (<-chan container.ContainerWaitOKBody, <-chan error) {
	waitChan := make(chan container.ContainerWaitOKBody)
	errChan := make(chan error)
	return waitChan, errChan
}

func (m *DockerClient) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	return types.ImageBuildResponse{}, nil
}

func (m *DockerClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	return []types.ImageSummary{}, nil
}
