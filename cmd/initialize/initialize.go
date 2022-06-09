package initialize

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/spf13/cobra"
)

var (
	format string
	prmApi *prm.Prm
)

func CreateCommand(parent *prm.Prm) *cobra.Command {
	prmApi = parent

	tmp := &cobra.Command{
		Use:     "initialize",
		Short:   "Initiates a directory with a `validate.yml` file",
		Long:    "Initiates a directory with a `validate.yml` file, for multi-tool validation",
		PreRunE: preExecute,
		RunE:    execute,
	}

	tmp.Flags().SortFlags = false
	tmp.Flags().StringVarP(&format, "format", "f", "human", "display output in human or json format")
	err := tmp.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"human", "json"}, cobra.ShellCompDirectiveNoFileComp
	})
	cobra.CheckErr(err)

	return tmp
}

func preExecute(cmd *cobra.Command, args []string) error {
	switch prmApi.RunningConfig.Backend {
	default:
		prmApi.Backend = &prm.Docker{AFS: prmApi.AFS, IOFS: prmApi.IOFS, ContextTimeout: prmApi.RunningConfig.Timeout}
	}
	return nil
}

func execute(cmd *cobra.Command, args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath := filepath.Join(wd, "validate.yml")

	_, err = prmApi.AFS.Stat(filePath)
	if err == nil {
		return fmt.Errorf("content has already been initialized")
	}

	file, err := prmApi.AFS.Create(filePath)
	if err != nil {
		return err
	}

	groups := `groups:
  - id: default
    tools:
      - name: puppetlabs/epp
      - name: puppetlabs/puppet-syntax
      - name: puppetlabs/metadata-json-lint
`

	_, err = file.WriteString(groups)
	if err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		log.Error().Msgf("Error closing file: %s", err)
	}

	log.Info().Msgf("PRM content initialized successfully")
	return nil
}
