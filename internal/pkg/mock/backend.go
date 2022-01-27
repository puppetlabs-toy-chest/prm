package mock

import (
	"errors"

	"github.com/puppetlabs/prm/pkg/prm"
)

type MockBackend struct {
	StatusIsAvailable   bool
	StatusMessageString string
	ToolAvalible        bool
	ExecReturn          string
}

func (m *MockBackend) Status() prm.BackendStatus {
	return prm.BackendStatus{IsAvailable: m.StatusIsAvailable, StatusMsg: m.StatusMessageString}
}

func (m *MockBackend) GetTool(tool *prm.Tool, prmConfig prm.Config) error {
	if m.ToolAvalible {
		return nil
	} else {
		return errors.New("Tool Not Found")
	}
}

// Implement when needed
func (m *MockBackend) Validate(tool *prm.Tool) (prm.ToolExitCode, error) {
	return prm.FAILURE, nil
}

func (m *MockBackend) Exec(tool *prm.Tool, args []string, prmConfig prm.Config, paths prm.DirectoryPaths) (prm.ToolExitCode, error) {
	switch m.ExecReturn {
	case "SUCCESS":
		return prm.SUCCESS, nil
	case "FAILURE":
		return prm.FAILURE, nil
	case "TOOL_ERROR":
		return prm.TOOL_ERROR, nil
	case "TOOL_NOT_FOUND":
		return prm.TOOL_NOT_FOUND, nil
	default:
		return prm.FAILURE, errors.New("DOCKER FAILURE")
	}
}
