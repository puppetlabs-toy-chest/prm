package set

import (
	"fmt"
	"github.com/puppetlabs/prm/pkg/config"
	"strings"

	"github.com/spf13/cobra"
)

var SelectedBackend config.BackendType

func (sc *SetCommand) createSetBackendCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:       "backend <BACKEND>",
		Short:     "Sets the backend exec environment to the specified type",
		Long:      `Sets the backend exec environment to the specified type`,
		PreRunE:   sc.setBackendPreRunE,
		RunE:      sc.setBackendType,
		ValidArgs: []string{string(config.DOCKER)},
	}

	return tmp
}

func (sc *SetCommand) setBackendPreRunE(cmd *cobra.Command, args []string) (err error) {
	if len(args) > 1 {
		return fmt.Errorf("too many args, please specify ONE of the following backend types after 'set backend':\n- %s", config.DOCKER)
	}

	if len(args) < 1 {
		return fmt.Errorf("please specify specify one of the following backend types after 'set backend':\n- %s", config.DOCKER)
	}

	switch strings.ToLower(args[0]) {
	case string(config.DOCKER):
		SelectedBackend = config.DOCKER
	default:
		return fmt.Errorf("'%s' is not a valid backend type, please specify one of the following backend types:\n- %s", args[0], config.DOCKER)
	}

	return nil
}

func (sc *SetCommand) setBackendType(cmd *cobra.Command, args []string) error {
	return config.SetAndWriteConfig(config.BackendCfgKey, string(SelectedBackend))
}
