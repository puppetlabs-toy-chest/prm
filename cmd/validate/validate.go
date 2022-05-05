package validate

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
	"github.com/puppetlabs/prm/internal/pkg/utils"

	"github.com/puppetlabs/pct/pkg/telemetry"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	localToolPath string
	format        string
	selectedTool  string
	listTools     bool
	prmApi        *prm.Prm
	toolArgs      string
	alwaysBuild   bool
	toolTimeout   int
	resultsView   string
	isSerial      bool
	workerCount   int
	selectedGroup string
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

	tmp.Flags().IntVar(&toolTimeout, "toolTimeout", 1800, "Time in seconds to wait for a response before exiting")
	err = viper.BindPFlag("toolTimeout", tmp.Flags().Lookup("toolTimeout"))
	cobra.CheckErr(err)

	tmp.Flags().StringVar(&resultsView, "resultsView", "", "Controls where results are outputted to, either 'terminal' or 'file' (Defaults: single tool = 'terminal', multiple tools = 'file')")
	err = viper.BindPFlag("resultsView", tmp.Flags().Lookup("resultsView"))
	cobra.CheckErr(err)

	tmp.Flags().BoolVar(&isSerial, "serial", false, "Runs validation one tool at a time instead of in parallel")
	err = viper.BindPFlag("serial", tmp.Flags().Lookup("serial"))
	cobra.CheckErr(err)

	tmp.Flags().IntVar(&workerCount, "workerCount", 10, "Worker count for running validation tools in parallel")
	err = viper.BindPFlag("workerCount", tmp.Flags().Lookup("workerCount"))
	cobra.CheckErr(err)

	tmp.Flags().StringVar(&selectedGroup, "group", "", "Select which tool group to use for multi-tool validation. Groups are defined inside of the validate.yml file.")
	err = viper.BindPFlag("group", tmp.Flags().Lookup("group"))
	cobra.CheckErr(err)

	return tmp
}

func preExecute(cmd *cobra.Command, args []string) error {
	if localToolPath == "" {
		localToolPath = prmApi.RunningConfig.ToolPath
	}

	if resultsView != "terminal" && resultsView != "file" && resultsView != "" {
		return fmt.Errorf("the --resultsView flag must be set to either [terminal|file]")
	}

	if prmApi.CodeDir == "" {
		workingDirectory, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("unable to set working directory as default codedir: %s", err)
		}
		prmApi.CodeDir = workingDirectory
	}

	if toolTimeout < 1 {
		return fmt.Errorf("the --toolTimeout flag must be set to a value greater than 1")
	}

	switch prmApi.RunningConfig.Backend {
	case prm.DOCKER:
		prmApi.Backend = &prm.Docker{AFS: prmApi.AFS, IOFS: prmApi.IOFS, AlwaysBuild: alwaysBuild, ContextTimeout: prmApi.RunningConfig.Timeout}
	default:
		prmApi.Backend = &prm.Docker{AFS: prmApi.AFS, IOFS: prmApi.IOFS, AlwaysBuild: alwaysBuild, ContextTimeout: prmApi.RunningConfig.Timeout}
	}

	if !listTools {
		doesExist, err := prmApi.AFS.DirExists(prmApi.CodeDir)
		if !doesExist {
			return fmt.Errorf("the --codedir flag must be set to a valid directory")
		}
		if err != nil {
			return err
		}
	}

	// handle the default cachepath
	if prmApi.CacheDir == "" {
		usr, _ := user.Current()
		dir := usr.HomeDir
		prmApi.CacheDir = filepath.Join(dir, ".pdk/prm/cache")
		err := prmApi.EnsureCacheDirExists()
		if err != nil {
			return err
		}
	}

	return prmApi.List(localToolPath, "", true)
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

		// Default resultsView for single tool validation is "file"
		if !cmd.Flags().Changed("resultsView") {
			resultsView = "terminal"
		}

		var additionalToolArgs []string
		if toolArgs != "" {
			additionalToolArgs, _ = shlex.Split(toolArgs)
		}

		toolInfo := prm.ToolInfo{
			Tool: cachedTool,
			Args: additionalToolArgs,
		}
		settings := prm.OutputSettings{
			ResultsView: resultsView,
			OutputDir:   path.Join(prmApi.CodeDir, ".prm-validate"),
		}

		err := prmApi.Validate([]prm.ToolInfo{toolInfo}, 1, settings)
		if err != nil {
			return err
		}
	} else {
		// Default resultsView for multitool validation is "file"
		if !cmd.Flags().Changed("resultsView") {
			resultsView = "file"
		}
		// No tool specified, so check if their code contains a validate.yml, which returns the list of tools
		// Their code is expected to be in the directory where the executable is run from
		toolGroup, err := prmApi.GetValidationGroupFromFile(selectedGroup)
		if err != nil {
			return err
		}

		outputDir := path.Join(prmApi.CodeDir, ".prm-validate")
		if toolGroup.ID != "" {
			outputDir = path.Join(outputDir, toolGroup.ID)
		}

		// Gather a list of tools
		var toolList []prm.ToolInfo
		for _, tool := range toolGroup.Tools {
			cachedTool, ok := prmApi.IsToolAvailable(tool.Name)
			if !ok {
				return fmt.Errorf("Tool %s not found in cache", tool)
			}

			info := prm.ToolInfo{Tool: cachedTool, Args: tool.Args}
			toolList = append(toolList, info)
		}

		if isSerial && cmd.Flags().Changed("workerCount") {
			log.Warn().Msgf("The --workerCount flag has no affect when used with the --serial flag")
		}
		if isSerial || workerCount < 1 {
			workerCount = 1
		}
		err = prmApi.Validate(toolList, workerCount, prm.OutputSettings{ResultsView: resultsView, OutputDir: outputDir})
		if err != nil {
			return err
		}
	}

	return nil
}
