package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Masterminds/semver"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	PuppetCmdFlag      string = "puppet"
	BackendCmdFlag     string = "backend"
	PuppetVerCfgKey    string = "puppetversion" // Should match Config struct key.
	BackendCfgKey      string = "backend"       // Should match Config struct key.
	DefaultPuppetVer   string = "7.15.0"
	DefaultBackend     string = "docker"
	ToolPathCfgKey     string = "toolpath"
	ToolTimeoutCfgKey  string = "toolTimeout"
	DefaultToolTimeout int    = 1800 // 30 minutes
)

var Config config

type config struct {
	PuppetVersion *semver.Version
	Backend       string
	ToolPath      string
	Timeout       time.Duration
}

func InitConfig() {
	var cfgFile string
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, _ := os.UserHomeDir()

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
	}

	generateDefaultCfg()

	if err := viper.ReadInConfig(); err != nil {
		err := viper.SafeWriteConfig()
		if err != nil {
			log.Error().Msgf("failed to write config: %s", err)
		}
	}

	viper.AutomaticEnv()
	err := viper.Unmarshal(&Config)
	if err != nil {
		log.Error().Msgf("failed to unmarshal config: %s", err)
	}
}

func generateDefaultCfg() {
	// Generate default configuration
	puppetVer, err := semver.NewVersion(DefaultPuppetVer)
	if err != nil {
		panic(fmt.Sprintf("Unable to generate default cfg value for 'puppet': %s", err))
	}

	log.Trace().Msgf("Setting default config (%s: %s)", PuppetVerCfgKey, puppetVer.String())
	viper.SetDefault(PuppetVerCfgKey, puppetVer)
	log.Trace().Msgf("Setting default config (%s: %s)", BackendCfgKey, DefaultBackend)
	viper.SetDefault(BackendCfgKey, string(DefaultBackend))

	defaultToolPath, err := getDefaultToolPath()
	if err != nil {
		panic(fmt.Sprintf("Unable to generate default cfg value for 'toolpath': %s", err))
	}
	log.Trace().Msgf("Setting default toolpath (%s: %s)", ToolPathCfgKey, defaultToolPath)
	viper.SetDefault(ToolPathCfgKey, defaultToolPath)

	log.Trace().Msgf("Setting default Timeout (%s: %d)", ToolTimeoutCfgKey, DefaultToolTimeout)
	viper.SetDefault(ToolTimeoutCfgKey, DefaultToolTimeout)
}

func loadConfig() error {
	// If the scenario where any other config value has been set AND the Puppet version is unset, a '{}' is written
	// to the config file on disk. This causes issues when attempting to call semver.NewVersion.
	// puppetVer := viper.GetString(PuppetVerCfgKey)
	// if puppetVer == "" {
	// 	puppetVer = DefaultPuppetVer
	// }
	// pupperSemVer, err := semver.NewVersion(puppetVer)
	// if err != nil {
	// 	return fmt.Errorf("could not load '%s' from config '%s': %s", PuppetVerCfgKey, viper.GetViper().ConfigFileUsed(), err)
	// }
	//
	// p.RunningConfig.PuppetVersion = pupperSemVer

	// Load Backend from config
	// p.RunningConfig.Backend = BackendType(viper.GetString(BackendCfgKey))

	// Load ToolPath from config
	// p.RunningConfig.ToolPath = viper.GetString(ToolPathCfgKey)

	// Load Timeout from config
	// p.RunningConfig.Timeout = viper.GetDuration(ToolTimeoutCfgKey) * time.Second

	return nil
}

func getDefaultToolPath() (string, error) {
	execDir, err := os.Executable()
	if err != nil {
		return "", err
	}

	defaultToolPath := filepath.Join(filepath.Dir(execDir), "tools")
	log.Trace().Msgf("Default tool config path: %v", defaultToolPath)
	return defaultToolPath, nil
}
