//nolint:structcheck,unused
package prm

type ValidateExitCode int64

const (
	VALIDATION_PASS ValidateExitCode = iota
	VALIDATION_FAILED
	VALIDATION_ERROR
)

// Validate allows a lits of tool names to be executed against
// the codeDir.
//
// Tools can be empty, in which case we expect that a local
// configuration file (validate.yml) will contain a list of
// tools to run.
func (*Prm) Validate(tools []string) (ValidateExitCode, error) {
	// TODO
	return VALIDATION_ERROR, nil
}
