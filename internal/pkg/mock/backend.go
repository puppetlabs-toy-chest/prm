package mock

import (
	"github.com/puppetlabs/prm/pkg/prm"
)

type MockBackend struct {
	StatusIsAvailable   bool
	StatusMessageString string
}

func (m *MockBackend) Status() prm.BackendStatus {
	return prm.BackendStatus{IsAvailable: m.StatusIsAvailable, StatusMsg: m.StatusMessageString}
}

// Implement when needed
func (m *MockBackend) GetTool(toolName string, prmConfig prm.Config) (prm.Tool, error) {
	return prm.Tool{}, nil
}

// Implement when needed
func (m *MockBackend) Validate(tool *prm.Tool) (prm.ToolExitCode, error) {
	return prm.FAILURE, nil
}

// Implement when needed
func (m *MockBackend) Exec(tool *prm.Tool, args []string) (prm.ToolExitCode, error) {
	return prm.FAILURE, nil
}
