package exec

import (
	"fmt"
	"path/filepath"
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
	listTools   bool
	prmApi      *prm.Prm
	cachedTools []prm.ToolConfig
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

	// tmp.Flags().StringVarP(&selectedToolInfo, "info", "i", "", "display the selected template's configuration and default values")
	// err = tmp.RegisterFlagCompletionFunc("info", flagCompletion)
	// cobra.CheckErr(err)

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

	return tmp
}

func preExecute(cmd *cobra.Command, args []string) error {
	if localToolPath == "" {
		localToolPath = prm.RunningConfig.ToolPath
	}
	cachedTools = prmApi.List(localToolPath, "")
	return nil
}

func validateArgCount(cmd *cobra.Command, args []string) error {
	// show available tools if user runs `prm exec`
	if len(args) == 0 && !listTools {
		listTools = true
	}

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
	for _, tool := range cachedTools {
		namespacedTemplate := fmt.Sprintf("%s/%s", tool.Plugin.Author, tool.Plugin.Id)
		if strings.HasPrefix(namespacedTemplate, match) {
			m := namespacedTemplate + "\t" + tool.Plugin.Display
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
		formattedTemplates, err := prmApi.FormatTools(cachedTools, format)
		if err != nil {
			return err
		}
		fmt.Print(formattedTemplates)

		return nil
	}

	matchingTools := prmApi.FilterFiles(cachedTools, func(f prm.ToolConfig) bool {
		return fmt.Sprintf("%s/%s", f.Plugin.Author, f.Plugin.Id) == selectedTool
	})

	if len(matchingTools) == 1 {
		matchingTool := matchingTools[0]
		selectedToolDirPath = filepath.Join(localToolPath, matchingTool.Plugin.Author, matchingTool.Plugin.Id, matchingTool.Plugin.Version)
		tool, err := prmApi.Get(selectedToolDirPath)
		if err != nil {
			return err
		}

		return prmApi.Exec(&tool, args[1:])

	} else {
		return fmt.Errorf("Couldn't find an installed tool that matches '%s'", selectedTool)
	}
}
