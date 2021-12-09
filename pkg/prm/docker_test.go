package prm_test

import (
	"reflect"
	"testing"

	"github.com/puppetlabs/prm/internal/pkg/mock"
	"github.com/puppetlabs/prm/pkg/prm"
)

func TestDocker_Status(t *testing.T) {
	tests := []struct {
		name       string
		mockClient mock.DockerClient
		want       prm.BackendStatus
	}{
		{
			name: "When connection unavailable",
			mockClient: mock.DockerClient{
				ErrorString: "error during connect: This error may indicate that the docker daemon is not running.: Get \"http://%2F%2F.%2Fpipe%2Fdocker_engine/v1.41/version\": open //./pipe/docker_engine: The system cannot find the file specified.",
			},
			want: prm.BackendStatus{
				IsAvailable: false,
				StatusMsg:   "error during connect: This error may indicate that the docker daemon is not running.",
			},
		},
		{
			name: "When an edge case failure occurs",
			mockClient: mock.DockerClient{
				ErrorString: "Something has gone terribly wrong!",
			},
			want: prm.BackendStatus{
				IsAvailable: false,
				StatusMsg:   "Something has gone terribly wrong!",
			},
		},
		{
			name: "When everything is working",
			mockClient: mock.DockerClient{
				Platform:   "Docker",
				Version:    "1.2.3",
				ApiVersion: "3.2.1",
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
			d := &prm.Docker{Client: &tt.mockClient}
			if got := d.Status(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Docker.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}
