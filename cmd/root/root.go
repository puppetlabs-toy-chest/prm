package root

import (
	"errors"
	"os"
	"strings"

	"github.com/puppetlabs/prm/cmd/exec"
	"github.com/puppetlabs/prm/cmd/explain"
	"github.com/puppetlabs/prm/cmd/get"
	"github.com/puppetlabs/prm/cmd/set"
	"github.com/puppetlabs/prm/cmd/status"
	"github.com/puppetlabs/prm/cmd/validate"
	"github.com/puppetlabs/prm/cmd/version"
	"github.com/puppetlabs/prm/pkg/config"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/puppetlabs/prm/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	LogLevel           string
	LocalTemplateCache string
	prmApi             *prm.Prm
	debug              bool
	errSilent          = errors.New("ErrSilent")
	currentVersion     = "dev"
	commit             = "none"
	date               = "unknown"

	//	format             string
)

func CreateRootCommand(prmApi *prm.Prm) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "prm",
		Short:            "prm - Puppet Runtime Manager",
		Long:             `Puppet Runtime Manager (PRM) - Execute commands and validate against Puppet content`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		SilenceUsage:     true,
		SilenceErrors:    true,
	}

	cmd.PersistentFlags().StringVar(&config.CfgFile, "config", "", "config file (default is $HOME/.config/.prm.yaml)")
	cmd.PersistentFlags().StringVar(&LogLevel, "log-level", zerolog.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")
	cmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug output")

	err := registerLogLevelFlagCompleteion(cmd)
	cobra.CheckErr(err)

	// version command
	v := utils.FormatVersion(currentVersion, date, commit)
	cmd.Version = v
	cmd.SetVersionTemplate(v)
	cmd.AddCommand(version.CreateVersionCommand(currentVersion, date, commit))

	// set command
	sc := set.SetCommand{Utils: &utils.Utils{}}
	cmd.AddCommand(sc.CreateSetCommand())

	// get command
	cmd.AddCommand(get.CreateGetCommand(prmApi))

	// exec command
	cmd.AddCommand(exec.CreateCommand(prmApi))

	// validate command
	cmd.AddCommand(validate.CreateCommand(prmApi))

	// status command
	cmd.AddCommand(status.CreateStatusCommand(prmApi))

	// explain
	cmd.AddCommand(explain.CreateCommand())

	// build
	// buildCmd := cmd_build.BuildCommand{
	// 	ProjectType: "tool",
	// 	Builder: &build.Builder{
	// 		Tar:  &tar.Tar{AFS: prmApi.AFS},
	// 		Gzip: &gzip.Gzip{AFS: prmApi.AFS},
	// 		AFS:  prmApi.AFS,
	// 		ConfigProcessor: &config_processor.ConfigProcessor{
	// 			AFS: prmApi.AFS,
	// 		},
	// 		ConfigFile: "prm-config.yml",
	// 	},
	// }
	// cmd.AddCommand(buildCmd.CreateCommand())
	//
	// // install command
	// installCmd := cmd_install.InstallCommand{
	// 	PrmInstaller: &install.Installer{
	// 		Tar:        &tar.Tar{AFS: prmApi.AFS},
	// 		Gunzip:     &gzip.Gunzip{AFS: prmApi.AFS},
	// 		AFS:        prmApi.AFS,
	// 		IOFS:       prmApi.IOFS,
	// 		HTTPClient: &http.Client{},
	// 		Exec:       &exec_runner.Exec{},
	// 		ConfigProcessor: &config_processor.ConfigProcessor{
	// 			AFS: prmApi.AFS,
	// 		},
	// 		ConfigFileName: "prm-config.yml",
	// 	},
	// 	AFS: prmApi.AFS,
	// }
	// cmd.AddCommand(installCmd.CreateCommand())
	//
	// tmp.PersistentFlags().StringVarP(&format, "format", "f", "junit", "formating (default is junit)")

	return cmd
}

func registerLogLevelFlagCompleteion(cmd *cobra.Command) error {
	return cmd.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var levels = []string{"debug", "info", "warn", "error", "fatal", "panic"}
		return find(levels, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})
}

// Returns the cobra command called, e.g. new or install
// and also the fully formatted command as passed with arguments/flags.
// Idea borrowed from carolynvs/porter:
// https://github.com/carolynvs/porter/blob/ccca10a63627e328616c1006600153da8411a438/cmd/porter/main.go
func GetCalledCommand(cmd *cobra.Command) (string, string) {
	if len(os.Args) < 2 {
		return "", ""
	}

	calledCommandStr := os.Args[1]

	// Also figure out the full called command from the CLI
	// Is there anything sensitive you could leak here? 🤔
	calledCommandArgs := strings.Join(os.Args[1:], " ")

	return calledCommandStr, calledCommandArgs
}

// Both contains and find are copied from the pdkgo repo because they lived in an internal pkg:
// github.com/puppetlabs/pct/internal/pkg/utils
// To use these directly, the utils pkg would need to be public
// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// finds a string present in a slice
func find(source []string, match string) []string {
	var matches []string
	if contains(source, match) {
		matches = append(matches, match)
	}
	return matches
}

//func Execute() {
//	cmd := CreateRootCommand()
//	if err := cmd.Execute(); err != nil {
//		if err != errSilent {
//			log.Fatal().Err(err).Msg("Failed to execute command")
//		}
//	}
//}
