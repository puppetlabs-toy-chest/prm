package set

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	PUPPET string = "puppet"
)

func CreateSetCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:       fmt.Sprintf("set %s [args]", PUPPET),
		Short:     "Sets the specified configuration to the specified value",
		Long:      "Sets the specified configuration to the specified value",
		ValidArgs: []string{PUPPET},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}
	return tmp
}
