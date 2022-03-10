package validate

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/puppetlabs/prm/internal/pkg/utils"

	"github.com/puppetlabs/pdkgo/pkg/telemetry"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	localToolPath string
	format        string
	selectedTool  string
	// selectedToolInfo    string
	listTools   bool
	prmApi      *prm.Prm
	toolArgs    string
	alwaysBuild bool
	toolTimeout int
	resultsView string
)

func CreateCommand(parent *prm.Prm) *cobra.Command {

	prmApi = parent

	tmp := &cobra.Command{
		Use:               "validate <tool> [overrides|flags]",
		Short:             "Validates Puppet Content with a given tool",
		Long:              `Validates Puppet Content with a given tool`,
		Args:              validateArgCount,
		ValidArgsFunction: flagCompletion,
		PreRunE:           preExecute,
		RunE:              execute,
	}

	tmp.Flags().SortFlags = false

	tmp.Flags().BoolVarP(&listTools, "list", "l", false, "list tools")
	err := tmp.RegisterFlagCompletionFunc("list", flagCompletion)
	cobra.CheckErr(err)

	tmp.Flags().StringVar(&format, "format", "table", "display output in human-readable or json format")
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

	tmp.Flags().StringVar(&prmApi.CacheDir, "cachedir", "", "location of cache used by PRM")
	err = viper.BindPFlag("cachedir", tmp.Flags().Lookup("cachedir"))
	cobra.CheckErr(err)

	tmp.Flags().StringVar(&toolArgs, "toolArgs", "", "Additional arguments to pass to the tool")
	err = viper.BindPFlag("toolArgs", tmp.Flags().Lookup("toolArgs"))
	cobra.CheckErr(err)

	tmp.Flags().BoolVarP(&alwaysBuild, "alwaysBuild", "a", false, "Rebuild the docker image for each tool execution, even if it already exists")
	err = viper.BindPFlag("alwaysBuild", tmp.Flags().Lookup("alwaysBuild"))
	cobra.CheckErr(err)

	tmp.Flags().IntVar(&toolTimeout, "toolTimeout", 1800, "Time in seconds to wait for a response before exiting; defaults to 1800 (i.e. 30 minutes)")
	err = viper.BindPFlag("toolTimeout", tmp.Flags().Lookup("toolTimeout"))
	cobra.CheckErr(err)

	tmp.Flags().StringVar(&resultsView, "resultsView", "terminal", "Controls where results are outputted to, either 'terminal' or 'file' (Defaults: single tool = 'terminal', multiple tools = 'file')")
	err = viper.BindPFlag("resultsView", tmp.Flags().Lookup("resultsView"))
	cobra.CheckErr(err)

	return tmp
}

func preExecute(cmd *cobra.Command, args []string) error {
	if localToolPath == "" {
		localToolPath = prmApi.RunningConfig.ToolPath
	}

	switch prmApi.RunningConfig.Backend {
	case prm.DOCKER:
		prmApi.Backend = &prm.Docker{AFS: prmApi.AFS, IOFS: prmApi.IOFS, AlwaysBuild: alwaysBuild, ContextTimeout: prmApi.RunningConfig.Timeout}
	default:
		prmApi.Backend = &prm.Docker{AFS: prmApi.AFS, IOFS: prmApi.IOFS, AlwaysBuild: alwaysBuild, ContextTimeout: prmApi.RunningConfig.Timeout}
	}

	// handle the default cachepath
	if prmApi.CacheDir == "" {
		usr, _ := user.Current()
		dir := usr.HomeDir
		prmApi.CacheDir = filepath.Join(dir, ".pdk/prm/cache")
	}

	prmApi.List(localToolPath, "", true)
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

	log.Trace().Msg("Validate")
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

		workingDir, err := os.Getwd()
		if err != nil {
			return err
		}

		log.Debug().Msgf("Working Directory: %v", workingDir)
		err = prmApi.Validate(cachedTool, prm.OutputSettings{OutputLocation: resultsView, OutputDir: path.Join(workingDir, ".prm-validate")})
		if err != nil {
			return err
		}
	}
	// Uncomment when implementing validate.yml
	// else {
	// 	// No tool specified, so check if their code contains a validate.yml, which returns the list of tools
	// 	// Their code is expected to be in the directory where the executable is run from
	// 	toolList, err := prmApi.CheckLocalConfig()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	log.Info().Msgf("Found tools: %v ", toolList)

	// 	for _, tool := range toolList {
	// 		cachedTool, ok := prmApi.IsToolAvailable(tool.Name)
	// 		if !ok {
	// 			return fmt.Errorf("Tool %s not found in cache", tool)
	// 		}

	// 		err := prmApi.Exec(cachedTool, tool.Args)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}
