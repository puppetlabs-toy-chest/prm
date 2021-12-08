package get

import (
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createGetPuppetCommand(parent *prm.Prm) *cobra.Command {
	tmp := &cobra.Command{
		Use:   "puppet",
		Short: "Gets the Puppet runtime version currently configured",
		Long:  "Gets the Puppet runtime version currently configured",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msgf("Puppet version is configured to: %s", parent.RunningConfig.PuppetVersion.String())
		},
	}

	return tmp
}
