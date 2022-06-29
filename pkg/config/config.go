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
	PuppetCmdFlag      string      = "puppet"
	BackendCmdFlag     string      = "backend"
	PuppetVerCfgKey    string      = "puppetversion" // Should match Config struct key.
	BackendCfgKey      string      = "backend"       // Should match Config struct key.
	DefaultPuppetVer   string      = "7.15.0"
	DefaultBackend     BackendType = DOCKER
	ToolPathCfgKey     string      = "toolpath"
	ToolTimeoutCfgKey  string      = "toolTimeout"
	DefaultToolTimeout int         = 1800 // 30 minutes
)

const (
	DOCKER BackendType = "docker"
)

type BackendType string

type Config struct {
	PuppetVersion *semver.Version
	Backend       BackendType
	ToolPath      string
	Timeout       time.Duration
}

func GenerateDefaultCfg() {
	// Generate default configuration
	puppetVer, err := semver.NewVersion(DefaultPuppetVer)
	if err != nil {
		panic(fmt.Sprintf("Unable to generate default cfg value for 'puppet': %s", err))
	}

	log.Trace().Msgf("Setting default config (%s: %s)", PuppetVerCfgKey, puppetVer.String())
	viper.SetDefault(PuppetVerCfgKey, puppetVer)
	log.Trace().Msgf("Setting default config (%s: %s)", BackendCfgKey, DefaultBackend)
	viper.SetDefault(BackendCfgKey, string(DefaultBackend))

	defaultToolPath, err := GetDefaultToolPath()
	if err != nil {
		panic(fmt.Sprintf("Unable to generate default cfg value for 'toolpath': %s", err))
	}
	log.Trace().Msgf("Setting default toolpath (%s: %s)", ToolPathCfgKey, defaultToolPath)
	viper.SetDefault(ToolPathCfgKey, defaultToolPath)

	log.Trace().Msgf("Setting default Timeout (%s: %d)", ToolTimeoutCfgKey, DefaultToolTimeout)
	viper.SetDefault(ToolTimeoutCfgKey, DefaultToolTimeout)
}

func LoadConfig() (Config, error) {
	// If the scenario where any other config value has been set AND the Puppet version is unset, a '{}' is written
	// to the config file on disk. This causes issues when attempting to call semver.NewVersion.
	puppetVer := viper.GetString(PuppetVerCfgKey)
	if puppetVer == "" {
		puppetVer = DefaultPuppetVer
	}
	puppetSemVer, err := semver.NewVersion(puppetVer)
	if err != nil {
		return Config{}, fmt.Errorf("could not load '%s' from config '%s': %s", PuppetVerCfgKey, viper.GetViper().ConfigFileUsed(), err)
	}

	config := Config{
		PuppetVersion: puppetSemVer,
		Backend:       BackendType(viper.GetString(BackendCfgKey)),
		ToolPath:      viper.GetString(ToolPathCfgKey),
		Timeout:       viper.GetDuration(ToolTimeoutCfgKey) * time.Second,
	}

	return config, nil
}

func GetDefaultToolPath() (string, error) {
	execDir, err := os.Executable()
	if err != nil {
		return "", err
	}

	defaultToolPath := filepath.Join(filepath.Dir(execDir), "tools")
	log.Trace().Msgf("Default tool config path: %v", defaultToolPath)
	return defaultToolPath, nil
}

func SetAndWriteConfig(k, v string) (err error) {
	log.Trace().Msgf("Setting and saving config '%s' to '%s' in %s", k, v, viper.ConfigFileUsed())

	viper.Set(k, v)

	if err = viper.WriteConfig(); err != nil {
		log.Error().Msgf("could not write config to %s: %s", viper.ConfigFileUsed(), err)
	}
	return err
}
