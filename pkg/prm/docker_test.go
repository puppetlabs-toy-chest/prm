package prm_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/Masterminds/semver"
	"github.com/docker/docker/api/types"
	"github.com/mitchellh/mapstructure"
	"github.com/puppetlabs/pct/pkg/install"
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

func TestDocker_Validate(t *testing.T) {
	logFileOutput := "/file/output"
	type fields struct {
		Client         prm.DockerClientI
		Context        context.Context
		ContextCancel  func()
		ContextTimeout time.Duration
		AFS            *afero.Afero
		IOFS           *afero.IOFS
		AlwaysBuild    bool
	}
	type args struct {
		paths          prm.DirectoryPaths
		author         string
		version        string
		id             string
		puppetVersion  string
		outputSettings prm.OutputSettings
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    prm.ValidateExitCode
		wantErr bool
	}{
		{
			name: "Fails as server version is invalid",
			fields: fields{
				Client: &mock.DockerClient{
					ErrorString: "Invalid server verison",
				},
			},
			want:    prm.VALIDATION_ERROR,
			wantErr: true,
		},
		{
			name: "Tool successfully validates",
			fields: fields{
				Client: &mock.DockerClient{
					ExitCode: 0,
				},
			},
			args: args{
				puppetVersion: "5.0.0",
				author:        "test-user",
				id:            "good-project",
				version:       "0.1.0",
			},
			want: prm.VALIDATION_PASS,
		},
		{
			name: "Tool returns a validation failure with error message",
			fields: fields{
				Client: &mock.DockerClient{
					ExitCode:     1,
					ExitErrorMsg: "Validation Failed",
				},
			},
			args: args{
				puppetVersion: "5.0.0",
				author:        "test-user",
				id:            "good-project",
				version:       "0.1.0",
			},
			want:    prm.VALIDATION_FAILED,
			wantErr: true,
		},
		{
			name: "Tool successfully validates and outputs to file",
			fields: fields{
				Client: &mock.DockerClient{
					ExitCode: 0,
				},
			},
			args: args{
				puppetVersion:  "5.0.0",
				author:         "test-user",
				id:             "good-project",
				version:        "0.1.0",
				outputSettings: prm.OutputSettings{OutputLocation: "file", OutputDir: logFileOutput},
			},
			want: prm.VALIDATION_PASS,
		},
		{
			name: "Tool returns a validation failure with error message and outputs to file",
			fields: fields{
				Client: &mock.DockerClient{
					ExitCode:     1,
					ExitErrorMsg: "Validation Failed",
				},
			},
			args: args{
				puppetVersion: "5.0.0",
				author:        "test-user",
				id:            "good-project",
				version:       "0.1.0",
			},
			want:    prm.VALIDATION_FAILED,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			iofs := &afero.IOFS{Fs: fs}

			d := &prm.Docker{
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
			prmConfig := prm.Config{
				PuppetVersion: puppetVersion,
			}

			tool := &prm.Tool{
				Cfg: prm.ToolConfig{
					Plugin: &prm.PluginConfig{
						ConfigParams: install.ConfigParams{
							Id:      tt.args.id,
							Author:  tt.args.author,
							Version: tt.args.version,
						},
					},
				},
			}

			got, err := d.Validate(tool, prmConfig, tt.args.paths, tt.args.outputSettings)
			if (err != nil) != tt.wantErr {
				t.Errorf("Docker.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Docker.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
