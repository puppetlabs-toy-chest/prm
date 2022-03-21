package validate_test

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/puppetlabs/prm/cmd/validate"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func nullFunction(cmd *cobra.Command, args []string) error {
	return nil
}

func TestCreateCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		returnCode int
		out        string
		wantCmd    *cobra.Command
		wantErr    bool
		createDirs []string
		f          func(cmd *cobra.Command, args []string) error
	}{
		{
			name:    "executes without error",
			f:       nullFunction,
			out:     "",
			wantErr: false,
		},
		{
			name:    "executes without error for valid flag",
			args:    []string{"author/templateId"},
			f:       nullFunction,
			out:     "",
			wantErr: false,
		},
		{
			name:    "executes with error when tool provided in the wrong format",
			args:    []string{"foo-bar"},
			f:       nullFunction,
			out:     "Selected tool must be in AUTHOR/ID format",
			wantErr: true,
		},
		{
			name:    "executes with error for invalid flag",
			args:    []string{"--foo"},
			f:       nullFunction,
			out:     "unknown flag: --foo",
			wantErr: true,
		},
		{
			name: "executes without error for valid arg for resultsView flag",
			args: []string{"--resultsView", "terminal"},
			f:    nullFunction,
		},
		{
			name:    "executes with error for invalid arg for resultsView flag",
			args:    []string{"--resultsView", "foo"},
			f:       nullFunction,
			out:     "the --resultsView flag must be set to either [terminal|file]",
			wantErr: true,
		},
		{
			name:    "executes with error for invalid toolTimeout flag",
			args:    []string{"--toolTimeout", "-1"},
			f:       nullFunction,
			out:     "the --toolTimeout flag must be set to a value greater than 1",
			wantErr: true,
		},
		{
			name:    "executes with error for invalid codedir",
			args:    []string{"--codedir", "random/dir"},
			f:       nullFunction,
			out:     "the --codedir flag must be set to a valid directory",
			wantErr: true,
		},
		{
			name:       "executes without error for valid codedir",
			args:       []string{"--codedir", "a/valid/codedir"},
			f:          nullFunction,
			wantErr:    false,
			createDirs: []string{"a/valid/codedir"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			// Create illusion of a valid tool dir
			toolDir := "path/to/tools"
			toolConfigPath := path.Join(toolDir, "puppetlabs/foo-bar/0.1.0/")
			fs.MkdirAll(toolConfigPath, 0755)                                 //nolint:gosec,errcheck
			file, _ := fs.Create(path.Join(toolConfigPath, "prm-config.yml")) //nolint:gosec,errcheck
			fileText := `---
plugin:
  author: puppetlabs
  id: foo-bar
  display: foo-bar
  version: 0.1.0
  upstream_project_url: https://github.com/puppetlabs/foo-bar/

common:
  can_validate: true`
			file.WriteString(fileText) //nolint:gosec,errcheck
			file.Close()               //nolint:gosec,errcheck

			for _, dir := range tt.createDirs {
				fs.MkdirAll(dir, 0755) //nolint:gosec,errcheck
			}

			prmObj := &prm.Prm{
				AFS:  &afero.Afero{Fs: fs},
				IOFS: &afero.IOFS{Fs: fs},
				RunningConfig: prm.Config{
					ToolPath: toolDir,
				},
			}
			cmd := validate.CreateCommand(prmObj)
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetErr(b)
			cmd.SetArgs(tt.args)
			cmd.RunE = tt.f

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("executeTestUnit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			out, err := ioutil.ReadAll(b)
			if err != nil {
				t.Errorf("Failed to read stdout: %v", err)
				return
			}

			assert.Contains(t, string(out), tt.out)
		})
	}
}
