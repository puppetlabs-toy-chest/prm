package prm_test

import (
	"reflect"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/docker/docker/api/types"
	"github.com/mitchellh/mapstructure"
	"github.com/puppetlabs/prm/internal/pkg/mock"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
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

func TestDocker_GetTool(t *testing.T) {
	type ToolInfo struct {
		id      string
		author  string
		version string
	}
	tests := []struct {
		name        string
		mockClient  mock.DockerClient
		tool        prm.Tool
		config      prm.Config
		toolInfo    ToolInfo
		errorMsg    string
		alwaysBuild bool
	}{
		{
			name:   "Image not found and create new image",
			config: prm.Config{PuppetVersion: &semver.Version{}},
		},
		{
			name:     "Image found and alwaysBuild set to false",
			toolInfo: ToolInfo{id: "test", author: "user", version: "0.1.0"},
			mockClient: mock.DockerClient{
				ImagesSlice: []types.ImageSummary{
					{
						RepoTags: []string{"pdk:puppet-0.0.0_user-test_0.1.0"},
						ID:       "foo",
					},
				},
			},
			config: prm.Config{PuppetVersion: &semver.Version{}},
		},
		{
			name:     "Image found and alwaysBuild set to true",
			toolInfo: ToolInfo{id: "test", author: "user", version: "0.1.0"},
			mockClient: mock.DockerClient{
				ImagesSlice: []types.ImageSummary{
					{
						RepoTags: []string{"pdk:puppet-0.0.0_user-test_0.1.0"},
						ID:       "foo",
					},
				},
			},
			alwaysBuild: true,
			config:      prm.Config{PuppetVersion: &semver.Version{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Map fields to the Plugin config which are then squashed into its struct
			toolinfo := map[string]interface{}{
				"id":      tt.toolInfo.id,
				"author":  tt.toolInfo.author,
				"version": tt.toolInfo.version,
			}
			_ = mapstructure.Decode(toolinfo, &tt.tool.Cfg.Plugin)

			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			d := &prm.Docker{Client: &tt.mockClient, AFS: afs, AlwaysBuild: tt.alwaysBuild}
			err := d.GetTool(&tt.tool, tt.config)
			if err != nil {
				assert.Contains(t, err.Error(), tt.errorMsg)
			}
		})
	}
}
