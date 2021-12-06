//nolint:structcheck,unused
package prm

type ExecExitCode int64

const (
	EXEC_SUCCESS ExecExitCode = iota
	EXEC_FAILURE
	EXEC_ERROR
)

// Executes a tool with the given arguments, against the codeDir.
func (p *Prm) Exec(tool *Tool, args []string) error {

	var toolList []string

	// perform a check for validate.yml

	// flatten the tool list
	p.flattenToolList(&toolList)

	var backend BackendI

	switch RunningConfig.Backend {
	case DOCKER:
		backend = &Docker{}
	default:
	}

	for _, toolName := range toolList {
		tool, err := backend.GetTool(toolName, RunningConfig)
	}

	return nil
}
