package set_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/puppetlabs/prm/cmd/set"
	"github.com/stretchr/testify/assert"
)

type test struct {
	name              string
	args              []string
	expectedOutput    string
	expectedPuppetVer string
	expectError       bool
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

func execTests(t *testing.T, tests []test) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			setCmd := set.CreateSetCommand()
			b := bytes.NewBufferString("")
			setCmd.SetOutput(b)
			setCmd.SetArgs(tt.args)

			err := setCmd.Execute()

			if (err != nil) && (!tt.expectError) {
				t.Errorf("Unexpected error message: %s", err)
				return
			}

			if tt.expectedOutput != "" {
				out, _ := ioutil.ReadAll(b)
				assert.Contains(t, string(out), tt.expectedOutput)
				return
			}

			if tt.expectedPuppetVer != "" {
				if set.PuppetSemVer.String() != tt.expectedPuppetVer {
					t.Errorf("Normalised Puppet version (%s) did not match expected version (%s)", set.PuppetSemVer.String(), tt.expectedPuppetVer)
				}
			}
		})
	}
}
