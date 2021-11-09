package get

import (
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createGetPuppetCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "puppet",
		Short: "Gets the Puppet runtime version currently configured",
		Long:  "Gets the Puppet runtime version currently configured",
		Run:   getPuppetVersion,
	}

	return tmp
}

func getPuppetVersion(cmd *cobra.Command, args []string) {
	log.Info().Msgf("Puppet version is configured to: %s", prm.RunningConfig.PuppetVersion.String())
}
