//nolint:structcheck,unused
package backend

import (
	"github.com/puppetlabs/prm/pkg/backend/docker"
	"github.com/puppetlabs/prm/pkg/config"
	"github.com/puppetlabs/prm/pkg/tool"
)

const (
	VALIDATION_PASS ValidateExitCode = iota
	VALIDATION_FAILED
	VALIDATION_ERROR
)

type BackendI interface {
	GetTool(tool *tool.Tool, prmConfig config.Config) error
	Validate(toolInfo ToolInfo, prmConfig config.Config, paths DirectoryPaths) (ValidateExitCode, string, error)
	Exec(tool *tool.Tool, args []string, prmConfig config.Config, paths DirectoryPaths) (tool.ToolExitCode, error)
	Status() BackendStatus
}

func GetBackend(backendType config.BackendType) BackendI {
	switch backendType {
	case config.DOCKER:
		return &docker.Docker{
			Client:         nil,
			Context:        nil,
			ContextCancel:  nil,
			ContextTimeout: 0,
			AFS:            nil,
			IOFS:           nil,
			AlwaysBuild:    false,
		}
	default:
		return nil
	}
}

// The BackendStatus must report whether the backend is available
// and any useful status information; in the case of the backend
// being unavailable, report the error message to the user.
type BackendStatus struct {
	IsAvailable   bool
	StatusMessage string
}

type DirectoryPaths struct {
	CodeDir  string
	CacheDir string
}

type OutputSettings struct {
	ResultsView string // Either "terminal" or "file"
	OutputDir   string // Directory to write log file to
}

type ToolInfo struct {
	Tool *tool.Tool
	Args []string
}

type ContainerOutput struct {
	Stdout string
	Stderr string
}

type ValidateExitCode int64

type ValidationOutput struct {
	Err      error
	ExitCode ValidateExitCode
	Stdout   string
}
