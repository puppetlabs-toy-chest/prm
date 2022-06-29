package mock

import (
	"errors"
	"github.com/puppetlabs/prm/pkg/backend"
	"github.com/puppetlabs/prm/pkg/backend/docker"
	"github.com/puppetlabs/prm/pkg/config"
	"github.com/puppetlabs/prm/pkg/tool"
)

type MockBackend struct {
	StatusIsAvailable   bool
	StatusMessageString string
	ToolAvalible        bool
	ExecReturn          string
	ValidateReturn      string
}

func (m *MockBackend) Status() backend.BackendStatus {
	return backend.BackendStatus{IsAvailable: m.StatusIsAvailable, StatusMessage: m.StatusMessageString}
}

func (m *MockBackend) GetTool(tool *tool.Tool, prmConfig config.Config) error {
	if m.ToolAvalible {
		return nil
	} else {
		return errors.New("Tool Not Found")
	}
}

// Implement when needed
func (m *MockBackend) Validate(toolInfo backend.ToolInfo, prmConfig config.Config, paths backend.DirectoryPaths) (backend.ValidateExitCode, string, error) {
	switch m.ValidateReturn {
	case "PASS":
		return backend.VALIDATION_PASS, "", nil
	case "FAIL":
		return backend.VALIDATION_FAILED, "", errors.New("VALIDATION FAIL")
	case "ERROR":
		return backend.VALIDATION_ERROR, "", errors.New("DOCKER ERROR")
	default:
		return backend.VALIDATION_ERROR, "", errors.New("DOCKER FAIL")
	}
}

func (m *MockBackend) Exec(t *tool.Tool, args []string, prmConfig config.Config, paths backend.DirectoryPaths) (tool.ToolExitCode, error) {
	switch m.ExecReturn {
	case "SUCCESS":
		return tool.SUCCESS, nil
	case "FAILURE":
		return tool.FAILURE, nil
	case "TOOL_ERROR":
		return tool.TOOL_ERROR, nil
	case "TOOL_NOT_FOUND":
		return tool.TOOL_NOT_FOUND, nil
	default:
		return tool.FAILURE, docker.ErrDockerNotRunning
	}
}
