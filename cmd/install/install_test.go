package install_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/spf13/afero"

	"github.com/puppetlabs/prm/cmd/install"
	"github.com/puppetlabs/prm/internal/pkg/mock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCreateinstallCommand(t *testing.T) {
	tests := []struct {
		name                string
		args                []string
		expectError         bool
		expectedToolPkgPath string
		expectedTargetDir   string
		viperToolPath       string
		expectedOutput      string
		expectedGitUri      string
	}{
		{
			name:           "Should error when no args provided",
			args:           []string{},
			expectError:    true,
			expectedOutput: "Path to tool package (tar.gz) should be first argument",
		},
		{
			name:           "Should error when > 1 arg provided",
			args:           []string{"first/arg", "second/undeed/arg"},
			expectError:    true,
			expectedOutput: "Incorrect number of arguments; path to tool package (tar.gz) should be first argument",
		},
		{
			name:                "Sets ToolPkgPath to passed arg and InstallPath to default tool dir",
			args:                []string{"/path/to/my-cool-tool.tar.gz"},
			expectError:         false,
			expectedToolPkgPath: "/path/to/my-cool-tool.tar.gz",
			expectedTargetDir:   "/the/default/location/for/tools",
			viperToolPath:       "/the/default/location/for/tools",
		},
		{
			name:                "Sets ToolPkgPath and InstallPath to passed args",
			args:                []string{"/path/to/my-cool-tool.tar.gz", "--toolpath", "/a/new/place/for/tools"},
			expectError:         false,
			expectedToolPkgPath: "/path/to/my-cool-tool.tar.gz",
			expectedTargetDir:   "/a/new/place/for/tools",
			viperToolPath:       "/the/default/location/for/tools",
		},
		{
			name:              "Sets GitUri to passed arg and InstallPath to default tool dir",
			args:              []string{"--git-uri", "https://github.com/puppetlabs/pct-test-tool-01.git"},
			viperToolPath:     "/the/default/location/for/tools",
			expectError:       false,
			expectedTargetDir: "/the/default/location/for/tools",
			expectedGitUri:    "https://github.com/puppetlabs/pct-test-tool-01.git",
		},
		{
			name:              "Sets GitUri and InstallPath to passed args",
			args:              []string{"--git-uri", "https://github.com/puppetlabs/pct-test-tool-01.git", "--toolpath", "/a/new/place/for/tools"},
			viperToolPath:     "/the/default/location/for/tools",
			expectError:       false,
			expectedTargetDir: "/a/new/place/for/tools",
			expectedGitUri:    "https://github.com/puppetlabs/pct-test-tool-01.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			viper.SetDefault("toolpath", tt.viperToolPath)
			cmd := install.InstallCommand{
				PrmInstaller: &mock.PctInstaller{
					ExpectedToolPkg:   tt.expectedToolPkgPath,
					ExpectedTargetDir: tt.expectedTargetDir,
					ExpectedGitUri:    tt.expectedGitUri,
				},
				AFS: &afero.Afero{Fs: fs},
			}
			installCmd := cmd.CreateCommand()

			b := bytes.NewBufferString("")
			installCmd.SetOutput(b)

			installCmd.SetArgs(tt.args)
			err := installCmd.Execute()

			if (err != nil) != tt.expectError {
				t.Errorf("executeTestUnit() error = %v, wantErr %v", err, tt.expectError)
				return
			}

			if (err != nil) && tt.expectError {
				out, _ := ioutil.ReadAll(b)
				assert.Contains(t, string(out), tt.expectedOutput)
			}

		})
	}
}
