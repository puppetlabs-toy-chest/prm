package prm_test

import (
	"testing"

	"github.com/puppetlabs/pdkgo/pkg/install"
	"github.com/puppetlabs/prm/internal/pkg/mock"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/spf13/afero"
)

func TestPrm_Validate(t *testing.T) {
	type fields struct {
		RunningConfig prm.Config
		CodeDir       string
		CacheDir      string
		Cache         map[string]*prm.Tool
		Backend       prm.BackendI
	}
	type args struct {
		id               string
		validateReturn   string // Could this not just be the actual enum?
		expectedErrMsg   string
		toolNotAvailable bool
		outputSettings   prm.OutputSettings
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Content passes validation",
			args: args{
				id:             "pass",
				validateReturn: "PASS",
			},
		},
		{
			name: "Content fails validation",
			args: args{
				id:             "fail",
				validateReturn: "FAIL",
				expectedErrMsg: "VALIDATION FAIL",
			},
		},
		{
			name: "Tool errors during validation",
			args: args{
				id:             "error",
				validateReturn: "ERROR",
				expectedErrMsg: "DOCKER ERROR",
			},
		},
		{
			name: "Fails to get tool",
			args: args{
				id:               "tool-get-fail",
				expectedErrMsg:   "Tool Not Found",
				toolNotAvailable: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			iofs := &afero.IOFS{Fs: fs}

			tool := &prm.Tool{
				Cfg: prm.ToolConfig{
					Plugin: &prm.PluginConfig{
						ConfigParams: install.ConfigParams{
							Id:      tt.args.id,
							Author:  "puppetlabs",
							Version: "0.1.0",
						},
					},
				},
			}

			p := &prm.Prm{
				AFS:           afs,
				IOFS:          iofs,
				RunningConfig: tt.fields.RunningConfig,
				CodeDir:       tt.fields.CodeDir,
				CacheDir:      tt.fields.CacheDir,
				Cache:         tt.fields.Cache,
				Backend: &mock.MockBackend{
					ToolAvalible: !tt.args.toolNotAvailable,
					ExecReturn:   tt.args.validateReturn,
				},
			}

			if err := p.Validate(tool, tt.args.outputSettings); err != nil && err.Error() != tt.args.expectedErrMsg {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
