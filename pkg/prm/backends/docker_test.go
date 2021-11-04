package backends_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/puppetlabs/prm/pkg/prm/backends"
)

type MockDockerClient struct {
	platform    string
	version     string
	apiVersion  string
	errorString string
}

func (m *MockDockerClient) ServerVersion(ctx context.Context) (types.Version, error) {
	if m.errorString != "" {
		return types.Version{}, fmt.Errorf(m.errorString)
	}
	versionInfo := &types.Version{
		Platform:   struct{ Name string }{m.platform},
		Version:    m.version,
		APIVersion: m.apiVersion,
	}
	return *versionInfo, nil
}

func TestDocker_Status(t *testing.T) {
	tests := []struct {
		name       string
		mockClient MockDockerClient
		want       prm.BackendStatus
	}{
		{
			name: "When connection unavailable",
			mockClient: MockDockerClient{
				errorString: "error during connect: This error may indicate that the docker daemon is not running.: Get \"http://%2F%2F.%2Fpipe%2Fdocker_engine/v1.41/version\": open //./pipe/docker_engine: The system cannot find the file specified.",
			},
			want: prm.BackendStatus{
				IsAvailable: false,
				StatusMsg:   "error during connect: This error may indicate that the docker daemon is not running.",
			},
		},
		{
			name: "When an edge case failure occurs",
			mockClient: MockDockerClient{
				errorString: "Something has gone terribly wrong!",
			},
			want: prm.BackendStatus{
				IsAvailable: false,
				StatusMsg:   "Something has gone terribly wrong!",
			},
		},
		{
			name: "When everything is working",
			mockClient: MockDockerClient{
				platform:   "Docker",
				version:    "1.2.3",
				apiVersion: "3.2.1",
			},
			want: prm.BackendStatus{
				IsAvailable: true,
				StatusMsg:   "\tPlatform: Docker\n\tVersion: 1.2.3\n\tAPI Version: 3.2.1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Uncomment to run unmocked
			// cli, _ := dockerClient.NewClientWithOpts(dockerClient.FromEnv)
			// d := &Docker{Client: cli}
			d := &backends.Docker{Client: &tt.mockClient}
			if got := d.Status(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Docker.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}
