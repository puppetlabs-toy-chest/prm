package set

import (
	"fmt"

	"github.com/puppetlabs/prm/pkg/config"
	"github.com/puppetlabs/prm/pkg/utils"
	"github.com/spf13/cobra"
)

type SetCommand struct {
	Utils utils.UtilsI
}

func (sc *SetCommand) CreateSetCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:                   fmt.Sprintf("set <%s|%s> value", config.BackendCmdFlag, config.BackendCmdFlag),
		Short:                 "Sets the specified configuration to the specified value",
		Long:                  "Sets the specified configuration to the specified value",
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{config.BackendCmdFlag, config.BackendCmdFlag},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	tmp.AddCommand(sc.createSetPuppetCommand())
	tmp.AddCommand(sc.createSetBackendCommand())

	return tmp
}
