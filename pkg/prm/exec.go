//nolint:structcheck,unused
package prm

type ExecExitCode int64

const (
	EXEC_SUCCESS ExecExitCode = iota
	EXEC_FAILURE
	EXEC_ERROR
)

// Executes a tool with the given arguments, against the codeDir.
func (*Prm) Exec(toolName string, args []string) (ExecExitCode, error) {
	// TODO
	return EXEC_ERROR, nil
}
