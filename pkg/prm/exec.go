//nolint:structcheck,unused
package prm

import "fmt"

type ExecExitCode int64

const (
	EXEC_SUCCESS ExecExitCode = iota
	EXEC_FAILURE
	EXEC_ERROR
)

// Executes a tool with the given arguments, against the codeDir.
func (*Prm) Exec(tool *Tool, args []string) error {

	if tool.Cfg.Gem != nil {
		fmt.Printf("GEM")
	}

	if tool.Cfg.Puppet != nil {
		fmt.Printf("PUPPET")
	}

	if tool.Cfg.Binary != nil {
		fmt.Printf("BINARY")
	}

	if tool.Cfg.Container != nil {
		fmt.Printf("CONTAINER")
	}

	return nil
}
