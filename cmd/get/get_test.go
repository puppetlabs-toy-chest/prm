package get_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/puppetlabs/prm/cmd/get"
	"github.com/stretchr/testify/assert"
)

type test struct {
	name           string
	args           []string
	expectedOutput string
	expectError    bool
}

func Test_GetCommand(t *testing.T) {
	tests := []test{
		{
			name:           "Should display help when no subcommand passed to 'get'",
			args:           []string{""},
			expectedOutput: "Displays the requested configuration value",
			expectError:    true,
		},
		{
			name:           "Should display help when invalid subcommand passed to 'get'",
			args:           []string{"foo"},
			expectedOutput: "Error: unknown command \"foo\" for \"get\"",
			expectError:    true,
		},
	}
	execTests(t, tests)
}

func execTests(t *testing.T, tests []test) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			getCmd := get.CreateGetCommand()
			b := bytes.NewBufferString("")
			getCmd.SetOutput(b)
			getCmd.SetArgs(tt.args)

			err := getCmd.Execute()

			if (err != nil) && (!tt.expectError) {
				t.Errorf("Unexpected error message: %s", err)
				return
			}

			out, _ := ioutil.ReadAll(b)

			if tt.expectedOutput != "" {
				assert.Contains(t, string(out), tt.expectedOutput)
			}
		})
	}
}
