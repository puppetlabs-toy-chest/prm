package exec_test

import (
	"github.com/puppetlabs/prm/pkg/backend/docker"
	"github.com/puppetlabs/prm/pkg/config"
	"github.com/puppetlabs/prm/pkg/tool"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/mitchellh/mapstructure"
	"github.com/puppetlabs/prm/internal/pkg/mock"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/stretchr/testify/assert"
)

func TestPrm_Exec(t *testing.T) {
	tests := []struct {
		name           string
		expectError    bool
		expectedErrMsg string
		p              *prm.Prm
		tool           *tool.Tool
		args           []string
		toolId         string
		toolAuthor     string
		toolVersion    string
	}{
		{
			name: "Tool is unavailible",
			p: &prm.Prm{
				RunningConfig: config.Config{
					PuppetVersion: semver.MustParse("7.15.0"),
					Backend:       config.DOCKER,
				},
				Backend: &mock.MockBackend{
					ToolAvalible:      false,
					StatusIsAvailable: true,
				},
			},
			expectedErrMsg: "Tool Not Found",
		},
		{
			name: "Tool is availible and reports Success",
			p: &prm.Prm{
				RunningConfig: config.Config{
					PuppetVersion: semver.MustParse("7.15.0"),
					Backend:       config.DOCKER,
				},
				Backend: &mock.MockBackend{
					ToolAvalible:      true,
					ExecReturn:        "SUCCESS",
					StatusIsAvailable: true,
				},
			},
			args:           []string{"Example"},
			toolId:         "test",
			toolAuthor:     "user",
			toolVersion:    "0.1.0",
			expectedErrMsg: "",
		},
		{
			name: "Tool is availible and reports Failure",
			p: &prm.Prm{
				RunningConfig: config.Config{
					PuppetVersion: semver.MustParse("7.15.0"),
					Backend:       config.DOCKER,
				},
				Backend: &mock.MockBackend{
					ToolAvalible:      true,
					ExecReturn:        "FAILURE",
					StatusIsAvailable: true,
				},
			},
			args:           []string{"Example"},
			toolId:         "test",
			toolAuthor:     "user",
			toolVersion:    "0.1.0",
			expectedErrMsg: "", // Tool has reported a failure
		},
		{
			name: "Tool is availible and reports Tool Error",
			p: &prm.Prm{
				RunningConfig: config.Config{
					PuppetVersion: semver.MustParse("7.15.0"),
					Backend:       config.DOCKER,
				},
				Backend: &mock.MockBackend{
					ToolAvalible:      true,
					ExecReturn:        "TOOL_ERROR",
					StatusIsAvailable: true,
				},
			},
			args:           []string{"Example"},
			toolId:         "test",
			toolAuthor:     "user",
			toolVersion:    "0.1.0",
			expectedErrMsg: "", // Tool has reported an error
		},
		{
			name: "Tool is availible and reports Tool Not Found",
			p: &prm.Prm{
				RunningConfig: config.Config{
					PuppetVersion: semver.MustParse("7.15.0"),
					Backend:       config.DOCKER,
				},
				Backend: &mock.MockBackend{
					ToolAvalible:      true,
					ExecReturn:        "TOOL_NOT_FOUND",
					StatusIsAvailable: true,
				},
			},
			args:           []string{"Example"},
			toolId:         "test",
			toolAuthor:     "user",
			toolVersion:    "0.1.0",
			expectedErrMsg: "", // Tool canot not be found
		},
		{
			name: "Error executing tool",
			p: &prm.Prm{
				RunningConfig: config.Config{
					PuppetVersion: semver.MustParse("7.15.0"),
					Backend:       config.DOCKER,
				},
				Backend: &mock.MockBackend{
					ToolAvalible:      true,
					ExecReturn:        "",
					StatusIsAvailable: false,
				},
			},
			args:           []string{"Example"},
			toolId:         "test",
			toolAuthor:     "user",
			toolVersion:    "0.1.0",
			expectedErrMsg: docker.ErrDockerNotRunning.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toolinfo := map[string]interface{}{
				"id":      tt.toolId,
				"author":  tt.toolAuthor,
				"version": tt.toolVersion,
			}
			var tool tool.Tool
			_ = mapstructure.Decode(toolinfo, &tool.Cfg.Plugin)
			tt.tool = &tool

			err := tt.p.Exec(tt.tool, tt.args)
			// If an error is expected and returned
			if tt.expectedErrMsg != "" && err != nil {
				assert.Contains(t, tt.expectedErrMsg, err.Error())
				return
			}

			// If no error is expected but one is returned
			if tt.expectedErrMsg == "" && err != nil {
				t.Errorf("LoadConfig() Unexpected error: %s", err)
				return
			}

			// If an error is expected but none is returned
			if tt.expectedErrMsg != "" && err == nil {
				t.Errorf("LoadConfig() Expected error not found: %s", tt.expectedErrMsg)
				return
			}
		})
	}
}
