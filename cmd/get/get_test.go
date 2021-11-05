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
	expectedErrMsg string
}

func Test_GetCommand(t *testing.T) {
	tests := []test{
		{
			name:           "Should display help when no subcommand passed to 'get'",
			args:           []string{""},
			expectedErrMsg: "Displays the requested configuration value",
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

			if (err != nil) && (tt.expectedErrMsg == "") {
				t.Errorf("Unexpected error message: %s", err)
				return
			}

			if tt.expectedErrMsg != "" {
				out, _ := ioutil.ReadAll(b)
				assert.Contains(t, string(out), tt.expectedErrMsg)
				return
			}
		})
	}
}
