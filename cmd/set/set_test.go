package set_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/prm/cmd/set"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type test struct {
	name               string
	args               []string
	expectedOutput     string
	expectedPuppetVer  string
	expectedBackedType prm.BackendType
	expectError        bool
}

func Test_SetCommand(t *testing.T) {
	tests := []test{
		{
			name:           "Should display help when no subcommand passed to 'set'",
			args:           []string{""},
			expectedOutput: "Sets the specified configuration to the specified value",
			expectError:    true,
		},
	}
	execTests(t, tests)
}

func Test_SetPuppetCommand(t *testing.T) {
	tests := []test{
		{
			name:           "Should display help when invalid subcommand passed to 'set'",
			args:           []string{"foo"},
			expectedOutput: "Error: unknown command \"foo\" for \"set\"",
			expectError:    true,
		},
		{
			name:              "Should keep 'X.Y.Z' ver as-is",
			args:              []string{"puppet", "7.10.1"},
			expectedPuppetVer: "7.10.1",
		},
		{
			name:              "Should normalise 'X' ver to 'X.Y.Z'",
			args:              []string{"puppet", "7"},
			expectedPuppetVer: "7.0.0",
		},
		{
			name:           "Should error when too many args supplied to 'puppet' sub cmd",
			args:           []string{"puppet", "7", "a", "b"},
			expectedOutput: "Error: only a single Puppet version can be set",
			expectError:    true,
		},
		{
			name:           "Should error when no arg supplied to 'puppet' sub cmd",
			args:           []string{"puppet"},
			expectedOutput: "Error: please specify a Puppet version after 'set puppet'",
			expectError:    true,
		},
		{
			name:           "Should error when invalid version supplied to 'puppet' sub cmd",
			args:           []string{"puppet", "foo"},
			expectedOutput: "Error: 'foo' is not a semantic (x.y.z) Puppet version",
			expectError:    true,
		},
	}
	execTests(t, tests)
}

func Test_SetBackendCommand(t *testing.T) {
	tests := []test{
		{
			name:               "Should handle valid backend selection (docker)",
			args:               []string{"backend", "docker"},
			expectedBackedType: prm.DOCKER,
		},
		{
			name:               "Should handle valid backend selection (dOcKeR)",
			args:               []string{"backend", "dOcKeR"},
			expectedBackedType: prm.DOCKER,
		},
		{
			name:           "Should error when too many args supplied to 'backend' sub cmd",
			args:           []string{"backend", "foo", "bar"},
			expectedOutput: fmt.Sprintf("Error: too many args, please specify ONE of the following backend types after 'set backend':\n- %s", prm.DOCKER),
			expectError:    true,
		},
		{
			name:           "Should error when no arg supplied to 'badckend' sub cmd",
			args:           []string{"backend"},
			expectedOutput: fmt.Sprintf("please specify specify one of the following backend types after 'set backend':\n- %s", prm.DOCKER),
			expectError:    true,
		},
		{
			name:           "Should error when invalid backend type supplied to 'badckend' sub cmd",
			args:           []string{"backend", "foo"},
			expectedOutput: fmt.Sprintf("Error: 'foo' is not a valid backend type, please specify one of the following backend types:\n- %s", prm.DOCKER),
			expectError:    true,
		},
	}
	execTests(t, tests)
}

func execTests(t *testing.T, tests []test) {
	// Init a test Viper cfg file
	cfgFile := filepath.Join(t.TempDir(), ".prm.yaml")

	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")

	_, err := os.Stat(cfgFile)
	if os.IsNotExist(err) {
		if _, err := os.Create(cfgFile); err != nil {
			panic(fmt.Sprintf("failed to initialise %s: %s", cfgFile, err))
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			setCmd := set.CreateSetCommand()
			b := bytes.NewBufferString("")
			setCmd.SetOutput(b)
			setCmd.SetArgs(tt.args)

			err := setCmd.Execute()

			fmt.Println(viper.AllKeys())
			fmt.Println(viper.GetViper().ConfigFileUsed())

			if (err != nil) && (!tt.expectError) {
				t.Errorf("Unexpected error message: %s", err)
				return
			}

			if tt.expectedOutput != "" {
				out, _ := ioutil.ReadAll(b)
				assert.Contains(t, string(out), tt.expectedOutput)
				return
			}

			viper.AutomaticEnv()

			if err := viper.ReadInConfig(); err != nil {
				panic(fmt.Errorf("Error loading config from %s: %s", viper.ConfigFileUsed(), err))
			}

			validatePuppetVer(t, tt)
			validateBackendType(t, tt)
		})
	}
}

func validatePuppetVer(t *testing.T, tt test) {
	if tt.expectedPuppetVer != "" {
		puppetVer := viper.GetString(prm.PuppetVerCfgKey)
		if puppetVer != tt.expectedPuppetVer {
			t.Errorf("Normalised Puppet version (%s) did not match expected version (%s)", viper.GetString(prm.PuppetVerCfgKey), tt.expectedPuppetVer)
		}
	}
}

func validateBackendType(t *testing.T, tt test) {
	if tt.expectedBackedType != "" {
		backend := viper.Get(prm.BackendCfgKey).(prm.BackendType)
		if backend != tt.expectedBackedType {
			t.Errorf("Normalised Backend type (%s) did not match expected backend type (%s)", backend, tt.expectedBackedType)
		}
	}
}
