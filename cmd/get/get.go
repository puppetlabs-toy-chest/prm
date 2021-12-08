package get

import (
	"fmt"

	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/spf13/cobra"
)

func CreateGetCommand(parent *prm.Prm) *cobra.Command {
	tmp := &cobra.Command{
		Use:                   fmt.Sprintf("get <%s|%s>", prm.BackendCmdFlag, prm.PuppetCmdFlag),
		Short:                 "Displays the requested configuration value",
		Long:                  "Displays the requested configuration value",
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{prm.BackendCmdFlag, prm.PuppetCmdFlag},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}
	tmp.AddCommand(createGetPuppetCommand(parent))
	tmp.AddCommand(createGetBackendCommand(parent))

	return tmp
}
