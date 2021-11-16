package get

import (
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createGetBackendCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "backend",
		Short: "Gets the Backend version currently configured",
		Long:  "Gets the Backend version currently configured",
		Run:   getBackend,
	}

	return tmp
}

func getBackend(cmd *cobra.Command, args []string) {
	log.Info().Msgf("Backend is configured to: %s", prm.RunningConfig.Backend)
}
