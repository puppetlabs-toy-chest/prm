//nolint:structcheck,unused
package prm

type BackendI interface {
	GetTool(toolName string, prmConfig Config) (Tool, error)
	Validate(tool *Tool) (ToolExitCode, error)
	Exec(tool *Tool, args []string) (ToolExitCode, error)
}

// The BackendStatus must report whether the backend is available
// and any useful status information; in the case of the backend
// being unavailable, report the error message to the user.
type BackendStatus struct {
	IsAvailable bool
	StatusMsg   string
}
