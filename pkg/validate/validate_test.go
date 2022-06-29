package validate_test

import (
	"fmt"
	"github.com/puppetlabs/pct/pkg/install"
	"github.com/puppetlabs/prm/pkg/backend"
	"github.com/puppetlabs/prm/pkg/backend/docker"
	"github.com/puppetlabs/prm/pkg/config"
	"github.com/puppetlabs/prm/pkg/tool"
	"github.com/puppetlabs/prm/pkg/validate"
	"testing"

	"github.com/puppetlabs/prm/internal/pkg/mock"
	"github.com/spf13/afero"
)

func TestPrm_Validate(t *testing.T) {
	pathToLogs := "path/to/tools"
	codeDirPath := "path/to/code"
	type fields struct {
		RunningConfig config.Config
		CodeDir       string
		CacheDir      string
		Cache         map[string]*tool.Tool
		Backend       backend.BackendI
	}
	type args struct {
		id                   string
		validateReturn       string
		expectedErrMsg       string
		toolNotAvailable     bool
		outputSettings       backend.OutputSettings
		toolArgs             []string
		workerCount          int
		extraTools           int
		statusIsNotAvailable bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Validation passes with one tool and results are returned to the terminal",
			args: args{
				id:             "my-tool",
				validateReturn: "PASS",
				outputSettings: backend.OutputSettings{
					ResultsView: "terminal",
					OutputDir:   pathToLogs,
				},
				workerCount: 1,
			},
		},
		{
			name: "Validation passes with multiple tools and results are returned to the terminal",
			args: args{
				id:             "my-tool",
				validateReturn: "PASS",
				outputSettings: backend.OutputSettings{
					ResultsView: "terminal",
					OutputDir:   pathToLogs,
				},
				extraTools:  10,
				workerCount: 10,
			},
		},
		{
			name: "Validation passes with one tool and results are returned to a file",
			args: args{
				id:             "my-tool",
				validateReturn: "PASS",
				outputSettings: backend.OutputSettings{
					ResultsView: "file",
					OutputDir:   pathToLogs,
				},
				workerCount: 1,
			},
		},
		{
			name: "Validation passes with multiple tools and results are returned to a file",
			args: args{
				id:             "my-tool",
				validateReturn: "PASS",
				outputSettings: backend.OutputSettings{
					ResultsView: "file",
					OutputDir:   pathToLogs,
				},
				extraTools:  10,
				workerCount: 10,
			},
		},
		{
			name: "Validation fails with a single tool and results are returned to the terminal",
			args: args{
				id:             "fail",
				validateReturn: "FAIL",
				expectedErrMsg: "Validation returned 1 error",
				outputSettings: backend.OutputSettings{
					ResultsView: "terminal",
					OutputDir:   pathToLogs,
				},
				workerCount: 1,
			},
			wantErr: true,
		},
		{
			name: "Validation fails with multiple tool and results are returned to the terminal",
			args: args{
				id:             "fail",
				validateReturn: "FAIL",
				expectedErrMsg: "Validation returned 11 errors",
				outputSettings: backend.OutputSettings{
					ResultsView: "terminal",
					OutputDir:   pathToLogs,
				},
				workerCount: 10,
				extraTools:  10,
			},
			wantErr: true,
		},
		{
			name: "Validation fails with a single tool and results are returned to the file",
			args: args{
				id:             "fail",
				validateReturn: "FAIL",
				expectedErrMsg: "Validation returned 1 error",
				outputSettings: backend.OutputSettings{
					ResultsView: "file",
					OutputDir:   pathToLogs,
				},
				workerCount: 1,
			},
			wantErr: true,
		},
		{
			name: "Validation fails with multiple tool and results are returned to the file",
			args: args{
				id:             "fail",
				validateReturn: "FAIL",
				expectedErrMsg: "Validation returned 11 errors",
				outputSettings: backend.OutputSettings{
					ResultsView: "file",
					OutputDir:   pathToLogs,
				},
				workerCount: 10,
				extraTools:  10,
			},
			wantErr: true,
		},
		{
			name: "Validation error caused by a tool not being available",
			args: args{
				id:             "fail",
				validateReturn: "FAIL",
				expectedErrMsg: "Validation returned 1 error",
				outputSettings: backend.OutputSettings{
					ResultsView: "terminal",
					OutputDir:   pathToLogs,
				},
				toolNotAvailable: true,
				workerCount:      1,
			},
			wantErr: true,
		},
		{
			name: "Validation error caused by an invalid `--resultsView flag being specified`",
			args: args{
				id:             "error",
				validateReturn: "FAIL",
				outputSettings: backend.OutputSettings{
					ResultsView: "invalid flag",
					OutputDir:   pathToLogs,
				},
				workerCount:    1,
				expectedErrMsg: "invalid --resultsView flag specified",
			},
			wantErr: true,
		},
		{
			name: "Validation error caused by backend not being available",
			args: args{
				id:             "error",
				validateReturn: "FAIL",
				outputSettings: backend.OutputSettings{
					ResultsView: "terminal",
					OutputDir:   pathToLogs,
				},
				workerCount:          1,
				expectedErrMsg:       docker.ErrDockerNotRunning.Error(),
				statusIsNotAvailable: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			//iofs := &afero.IOFS{Fs: fs}

			if tt.args.outputSettings.ResultsView == "file" {
				afs.MkdirAll(tt.args.outputSettings.OutputDir, 0755) //nolint:gosec,errcheck
			}
			afs.MkdirAll(pathToLogs, 0755)  //nolint:gosec,errcheck
			afs.MkdirAll(codeDirPath, 0755) //nolint:gosec,errcheck

			var tools []backend.ToolInfo
			for i := 0; i < tt.args.extraTools+1; i++ {
				id, author, version := fmt.Sprint(tt.args.id, i), "puppetlabs", "0.1.0"
				tools = append(tools, CreateToolInfo(id, author, version, tt.args.toolArgs))
			}

			validator := validate.Validator{
				Backend: &mock.MockBackend{
					StatusIsAvailable: !tt.args.statusIsNotAvailable,
					ToolAvalible:      !tt.args.toolNotAvailable,
					ValidateReturn:    tt.args.validateReturn,
				},
				AFS:            afs,
				DirectoryPaths: backend.DirectoryPaths{CodeDir: codeDirPath, CacheDir: tt.fields.CacheDir},
				RunningConfig:  tt.fields.RunningConfig,
			}

			if err := validator.Validate(tools, tt.args.workerCount, tt.args.outputSettings); err != nil && err.Error() != tt.args.expectedErrMsg {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
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
