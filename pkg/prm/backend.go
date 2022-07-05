//nolint:structcheck,unused
package prm

type BackendI interface {
	GetTool(tool *Tool) error
	Validate(toolInfo ToolInfo, paths DirectoryPaths) (ValidateExitCode, string, error)
	Exec(tool *Tool, args []string, paths DirectoryPaths) (ToolExitCode, error)
	Status() BackendStatus
}

// The BackendStatus must report whether the backend is available
// and any useful status information; in the case of the backend
// being unavailable, report the error message to the user.
type BackendStatus struct {
	IsAvailable bool
	StatusMsg   string
}

type DirectoryPaths struct {
	codeDir  string
	cacheDir string
}

type OutputSettings struct {
	ResultsView string // Either "terminal" or "file"
	OutputDir   string // Directory to write log file to
}

type ToolInfo struct {
	Tool *Tool
	Args []string
}

type ContainerOutput struct {
	stdout string
	stderr string
}
