package docker_test

import (
	"context"
	"github.com/puppetlabs/prm/pkg/backend"
	"github.com/puppetlabs/prm/pkg/backend/docker"
	"github.com/puppetlabs/prm/pkg/config"
	"github.com/puppetlabs/prm/pkg/tool"
	"github.com/puppetlabs/prm/pkg/validate"
	"reflect"
	"testing"
	"time"

	"github.com/Masterminds/semver"
	"github.com/docker/docker/api/types"
	"github.com/mitchellh/mapstructure"
	"github.com/puppetlabs/pct/pkg/install"
	"github.com/puppetlabs/prm/internal/pkg/mock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestDocker_Status(t *testing.T) {
	tests := []struct {
		name       string
		mockClient mock.DockerClient
		want       backend.BackendStatus
	}{
		{
			name: "When connection unavailable",
			mockClient: mock.DockerClient{
				ErrorString: "error during connect: This error may indicate that the docker daemon is not running.: Get \"http://%2F%2F.%2Fpipe%2Fdocker_engine/v1.41/version\": open //./pipe/docker_engine: The system cannot find the file specified.",
			},
			want: backend.BackendStatus{
				IsAvailable: false,
				StatusMsg:   "error during connect: This error may indicate that the docker daemon is not running.",
			},
		},
		{
			name: "When an edge case failure occurs",
			mockClient: mock.DockerClient{
				ErrorString: "Something has gone terribly wrong!",
			},
			want: backend.BackendStatus{
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
			want: backend.BackendStatus{
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
			d := &docker.Docker{Client: &tt.mockClient}
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
		tool        tool.Tool
		config      config.Config
		toolInfo    ToolInfo
		errorMsg    string
		alwaysBuild bool
	}{
		{
			name:   "Image not found and create new image",
			config: config.Config{PuppetVersion: &semver.Version{}},
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
			config: config.Config{PuppetVersion: &semver.Version{}},
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
			config:      config.Config{PuppetVersion: &semver.Version{}},
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
			d := &docker.Docker{Client: &tt.mockClient, AFS: afs, AlwaysBuild: tt.alwaysBuild}
			err := d.GetTool(&tt.tool, tt.config)
			if err != nil {
				assert.Contains(t, err.Error(), tt.errorMsg)
			}
		})
	}
}

func TestDocker_Validate(t *testing.T) {
	defaultStdoutText := "This is stdout"
	type fields struct {
		Client         *mock.DockerClient
		Context        context.Context
		ContextCancel  func()
		ContextTimeout time.Duration
		AFS            *afero.Afero
		IOFS           *afero.IOFS
		AlwaysBuild    bool
	}
	type args struct {
		paths         backend.DirectoryPaths
		author        string
		version       string
		id            string
		puppetVersion string
		toolArgs      []string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       validate.ValidateExitCode
		wantErr    bool
		wantStdout string
	}{
		{
			name: "Fails as server version is invalid",
			fields: fields{
				Client: &mock.DockerClient{
					ErrorString: "Invalid server verison",
				},
			},
			want:    validate.VALIDATION_ERROR,
			wantErr: true,
		},
		{
			name: "Tool successfully validates",
			fields: fields{
				Client: &mock.DockerClient{
					ExitCode: 0,
					Stdout:   defaultStdoutText,
				},
			},
			args: args{
				puppetVersion: "5.0.0",
				author:        "test-user",
				id:            "good-project",
				version:       "0.1.0",
			},
			want:       validate.VALIDATION_PASS,
			wantStdout: defaultStdoutText,
		},
		{
			name: "Tool successfully validates with tool args",
			fields: fields{
				Client: &mock.DockerClient{
					ExitCode: 0,
					Stdout:   defaultStdoutText,
				},
			},
			args: args{
				puppetVersion: "5.0.0",
				author:        "test-user",
				id:            "good-project",
				version:       "0.1.0",
				toolArgs:      []string{"-l", "-v"},
			},
			want:       validate.VALIDATION_PASS,
			wantStdout: defaultStdoutText,
		},
		{
			name: "Tool returns a validation failure with error message",
			fields: fields{
				Client: &mock.DockerClient{
					ExitCode: 1,
					Stderr:   "Tool found 1 validation error",
					Stdout:   defaultStdoutText,
				},
			},
			args: args{
				puppetVersion: "5.0.0",
				author:        "test-user",
				id:            "good-project",
				version:       "0.1.0",
			},
			want:       validate.VALIDATION_FAILED,
			wantErr:    true,
			wantStdout: defaultStdoutText,
		},
		{
			name: "Error occurs while trying to validate with a tool",
			fields: fields{
				Client: &mock.DockerClient{
					ExitCode:     1,
					ExitErrorMsg: "Validation Failed",
					WantChanErr:  true,
				},
			},
			args: args{
				puppetVersion: "5.0.0",
				author:        "test-user",
				id:            "good-project",
				version:       "0.1.0",
			},
			want:    validate.VALIDATION_ERROR,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			iofs := &afero.IOFS{Fs: fs}

			d := &docker.Docker{
				Client:         tt.fields.Client,
				Context:        tt.fields.Context,
				ContextCancel:  tt.fields.ContextCancel,
				ContextTimeout: tt.fields.ContextTimeout,
				AFS:            afs,
				IOFS:           iofs,
				AlwaysBuild:    tt.fields.AlwaysBuild,
			}

			if tt.args.puppetVersion == "" {
				tt.args.puppetVersion = "5.0.0"
			}
			puppetVersion, err := semver.NewVersion(tt.args.puppetVersion)
			if err != nil {
				t.Errorf("Invalid Puppet Version %s", tt.args.puppetVersion)
				return
			}
			prmConfig := config.Config{
				PuppetVersion: puppetVersion,
			}

			toolInfo := CreateToolInfo(tt.args.id, tt.args.author, tt.args.version, tt.args.toolArgs)

			got, stdout, err := d.Validate(toolInfo, prmConfig, tt.args.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("Docker.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Docker.Validate() = %v, want %v", got, tt.want)
			}
			if stdout != tt.wantStdout {
				t.Errorf("Docker.Validate() = %v, want %v", stdout, tt.wantStdout)
			}
		})
	}
}

func CreateToolInfo(id, author, version string, args []string) backend.ToolInfo {
	tool := &tool.Tool{
		Cfg: tool.ToolConfig{
			Plugin: &tool.PluginConfig{
				ConfigParams: install.ConfigParams{
					Id:      id,
					Author:  author,
					Version: version,
				},
			},
		},
	}

	return backend.ToolInfo{
		Tool: tool,
		Args: args,
	}
}
