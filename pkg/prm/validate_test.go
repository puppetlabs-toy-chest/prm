package prm_test

import (
	"fmt"
	"testing"

	"github.com/puppetlabs/prm/internal/pkg/mock"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/spf13/afero"
)

func TestPrm_Validate(t *testing.T) {
	pathToLogs := "path/to/tools"
	codeDirPath := "path/to/code"
	type fields struct {
		RunningConfig prm.Config
		CodeDir       string
		CacheDir      string
		Cache         map[string]*prm.Tool
		Backend       prm.BackendI
	}
	type args struct {
		id               string
		validateReturn   string
		expectedErrMsg   string
		toolNotAvailable bool
		outputSettings   prm.OutputSettings
		toolArgs         []string
		workerCount      int
		extraTools       int
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
				outputSettings: prm.OutputSettings{
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
				outputSettings: prm.OutputSettings{
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
				outputSettings: prm.OutputSettings{
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
				outputSettings: prm.OutputSettings{
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
				outputSettings: prm.OutputSettings{
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
				outputSettings: prm.OutputSettings{
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
				outputSettings: prm.OutputSettings{
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
				outputSettings: prm.OutputSettings{
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
				outputSettings: prm.OutputSettings{
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
				outputSettings: prm.OutputSettings{
					ResultsView: "invalid flag",
					OutputDir:   pathToLogs,
				},
				workerCount:    1,
				expectedErrMsg: "invalid --resultsView flag specified",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			iofs := &afero.IOFS{Fs: fs}

			if tt.args.outputSettings.ResultsView == "file" {
				afs.MkdirAll(tt.args.outputSettings.OutputDir, 0755) //nolint:gosec,errcheck
			}
			afs.MkdirAll(pathToLogs, 0755)  //nolint:gosec,errcheck
			afs.MkdirAll(codeDirPath, 0755) //nolint:gosec,errcheck

			var tools []prm.ToolInfo
			for i := 0; i < tt.args.extraTools+1; i++ {
				id, author, version := fmt.Sprint(tt.args.id, i), "puppetlabs", "0.1.0"
				tools = append(tools, CreateToolInfo(id, author, version, tt.args.toolArgs))
			}

			p := &prm.Prm{
				AFS:           afs,
				IOFS:          iofs,
				RunningConfig: tt.fields.RunningConfig,
				CodeDir:       codeDirPath,
				CacheDir:      tt.fields.CacheDir,
				Cache:         tt.fields.Cache,
				Backend: &mock.MockBackend{
					ToolAvalible:   !tt.args.toolNotAvailable,
					ValidateReturn: tt.args.validateReturn,
				},
			}

			if err := p.Validate(tools, tt.args.workerCount, tt.args.outputSettings); err != nil && err.Error() != tt.args.expectedErrMsg {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
