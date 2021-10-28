package backends

import "github.com/puppetlabs/prm/pkg/prm"

type Docker struct {
}

func (*Docker) GetTool(toolName string, prmConfig prm.Config) (prm.Tool, error) {
	// TODO
	return prm.Tool{}, nil
}

func (*Docker) Validate(tool *prm.Tool) (prm.ToolExitCode, error) {
	// TODO
	return prm.FAILURE, nil
}

func (*Docker) Exec(tool *prm.Tool, args []string) (prm.ToolExitCode, error) {
	// TODO
	return prm.FAILURE, nil
}
