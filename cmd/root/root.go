package root

import (
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile            string
	LogLevel           string
	LocalTemplateCache string
	prmApi             *prm.Prm

	debug bool
	// format string
)

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	lvl, err := zerolog.ParseLevel(LogLevel)
	if err != nil {
		panic("Could not initialize zerolog")
	}

	zerolog.SetGlobalLevel(lvl)

	if lvl == zerolog.InfoLevel {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()
	}

	log.Trace().Msg("Initialized zerolog")
}

func CreateRootCommand(parent *prm.Prm) *cobra.Command {
	prmApi = parent

	tmp := &cobra.Command{
		Use:   "prm",
		Short: "prm - Puppet Runtime Manager",
		Long:  `Puppet Runtime Manager (PRM) - Execute commands and validate against Puppet content`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	tmp.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/.prm.yaml)")

	tmp.PersistentFlags().StringVar(&LogLevel, "log-level", zerolog.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")
	err := tmp.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var levels = []string{"debug", "info", "warn", "error", "fatal", "panic"}
		return find(levels, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})
	cobra.CheckErr(err)

	tmp.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug output")
	// tmp.PersistentFlags().StringVarP(&format, "format", "f", "junit", "formating (default is junit)")

	return tmp
}

func InitConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, _ := homedir.Dir()
		cfgFile = ".prm.yaml"
		viper.SetConfigName(cfgFile)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
		cfgPath := filepath.Join(home, ".config")
		viper.AddConfigPath(cfgPath)

		if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
			log.Trace().Msgf("%s does not exist, creating", cfgPath)
			if err := os.MkdirAll(cfgPath, 0750); err != nil {
				log.Error().Msgf("failed to create dir %s: %s", cfgPath, err)
			}
		}

		cfgFilePath := filepath.Join(cfgPath, cfgFile)

		if _, err := os.Stat(cfgFilePath); os.IsNotExist(err) {
			_, err := os.Create(filepath.Clean(cfgFilePath))
			if err != nil {
				log.Error().Msgf("failed to initialise %s: %s", cfgFilePath, err)
			}
		}
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Trace().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}

	prmApi.GenerateDefaultCfg()

	if err := prmApi.LoadConfig(); err != nil {
		log.Warn().Msgf("Error setting running config: %s", err)
	}
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
// github.com/puppetlabs/pdkgo/internal/pkg/utils
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
