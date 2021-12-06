package exec

import (
	"fmt"
	"strings"

	"github.com/puppetlabs/prm/internal/pkg/utils"

	"github.com/puppetlabs/pdkgo/pkg/telemetry"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	localToolPath       string
	format              string
	selectedTool        string
	selectedToolDirPath string
	// selectedToolInfo    string
	listTools bool
	prmApi    *prm.Prm
)

func CreateCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:               "exec <tool> [overrides|flags]",
		Short:             "Executes a given tool against some Puppet Content",
		Long:              `Executes a given tool against some Puppet Content`,
		Args:              validateArgCount,
		ValidArgsFunction: flagCompletion,
		PreRunE:           preExecute,
		RunE:              execute,
	}

	// Configure PRM
	fs := afero.NewOsFs() // configure afero to use real filesystem
	prmApi = &prm.Prm{
		AFS:  &afero.Afero{Fs: fs},
		IOFS: &afero.IOFS{Fs: fs},
	}

	tmp.Flags().SortFlags = false

	tmp.Flags().BoolVarP(&listTools, "list", "l", false, "list tools")
	err := tmp.RegisterFlagCompletionFunc("list", flagCompletion)
	cobra.CheckErr(err)

	tmp.Flags().StringVar(&format, "format", "table", "display output in table or json format")
	err = tmp.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var formats = []string{"table", "json"}
		return utils.Find(formats, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})
	cobra.CheckErr(err)

	tmp.Flags().StringVar(&localToolPath, "toolpath", "", "location of installed tools")
	err = viper.BindPFlag("toolpath", tmp.Flags().Lookup("toolpath"))
	cobra.CheckErr(err)

	tmp.Flags().StringVar(&prmApi.CodeDir, "codedir", "", "location of code to execute against")
	err = viper.BindPFlag("codedir", tmp.Flags().Lookup("codedir"))
	cobra.CheckErr(err)

	return tmp
}

func preExecute(cmd *cobra.Command, args []string) error {
	if localToolPath == "" {
		localToolPath = prm.RunningConfig.ToolPath
	}
	prmApi.List(localToolPath, "")
	return nil
}

func validateArgCount(cmd *cobra.Command, args []string) error {
	if len(args) >= 1 {
		if len(strings.Split(args[0], "/")) != 2 {
			return fmt.Errorf("Selected tool must be in AUTHOR/ID format")
		}
		selectedTool = args[0]
	}

	return nil
}

func flagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if localToolPath == "" {
		err := preExecute(cmd, args)
		if err != nil {
			log.Error().Msgf("Unable to set tool path: %s", err.Error())
			return nil, cobra.ShellCompDirectiveError
		}
	}
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	localToolPath = viper.GetString(prm.ToolPathCfgKey)

	return completeName(localToolPath, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

func completeName(cache string, match string) []string {
	var names []string
	for toolName, tool := range prmApi.Cache {
		if strings.HasPrefix(toolName, match) {
			m := toolName + "\t" + tool.Cfg.Plugin.Display
			names = append(names, m)
		}
	}
	return names
}

func execute(cmd *cobra.Command, args []string) error {
	span := telemetry.GetSpanFromContext(cmd.Context())
	// Add tool to span if needed
	if len(args) == 1 {
		telemetry.AddStringSpanAttribute(span, "tool", args[0])
	}

	log.Trace().Msg("Run")
	log.Trace().Msgf("Tool path: %v", localToolPath)
	log.Trace().Msgf("Selected Tool: %v", selectedTool)

	if listTools {
		formattedTemplates, err := prmApi.FormatTools(prmApi.Cache, format)
		if err != nil {
			return err
		}
		fmt.Print(formattedTemplates)

		return nil
	}

	if selectedTool != "" {
		// get the tool from the cache
		cachedTool, ok := prmApi.IsToolAvailable(selectedTool)
		if !ok {
		 return fmt.Errorf("Tool %s not found in cache", selectedTool)
		}
		// execute!
		err := prmApi.Exec(cachedTool, args[1:])
		if err != nil {
			return err
		}
	} else {
		// No tool specified, so check if their code contains a validate.yml, which returns the list of tools
		// Their code is expected to be in the directory where the executable is run from
		toolList, err := prmApi.CheckLocalConfig()
		if err != nil {
			return err
		}

		log.Info().Msgf("Found tools: %v ", toolList)

	}

	return nil
}
