package set

import (
	"fmt"

	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/spf13/cobra"
)

func CreateSetCommand() *cobra.Command {
	tmp := &cobra.Command{
		// TODO: Add Backend as valid arg
		Use:                   fmt.Sprintf("set %s", prm.PuppetCmdFlag),
		Short:                 "Sets the specified configuration to the specified value",
		Long:                  "Sets the specified configuration to the specified value",
		DisableFlagsInUseLine: true,
		// TODO: Add Backend as valid arg		
		ValidArgs:             []string{prm.PuppetCmdFlag},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	tmp.AddCommand(createSetPuppetCommand())
	tmp.AddCommand(createSetBackendCommand())

	return tmp
}
