package status

import (
	"fmt"

	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/spf13/cobra"
)

var (
	format string
	prmApi *prm.Prm
)

func CreateStatusCommand(parent *prm.Prm) *cobra.Command {
	prmApi = parent

	tmp := &cobra.Command{
		Use:     "status",
		Short:   "Returns runtime status",
		Long:    "Returns the status of the puppet manager runtime as currently configured, including both puppet version and backend",
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
	status, err := prm.FormatStatus(prmApi.GetStatus(), format)
	if err != nil {
		return err
	}
	fmt.Print(status)

	return nil
}
