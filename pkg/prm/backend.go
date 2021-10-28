//nolint:structcheck,unused
package prm

type BackendI interface {
	GetTool(toolName string, prmConfig Config) (Tool, error)
	Validate(tool *Tool) (ToolExitCode, error)
	Exec(tool *Tool, args []string) (ToolExitCode, error)
}
