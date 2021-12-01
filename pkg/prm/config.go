package prm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	PuppetCmdFlag    string      = "puppet"
	BackendCmdFlag   string      = "backend"
	PuppetVerCfgKey  string      = "puppetversion" // Should match Config struct key.
	BackendCfgKey    string      = "backend"       // Should match Config struct key.
	DefaultPuppetVer string      = "7"
	DefaultBackend   BackendType = DOCKER
	ToolPathCfgKey   string      = "toolpath"
)

type Config struct {
	PuppetVersion *semver.Version
	Backend       BackendType
	ToolPath      string
}

var RunningConfig Config

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
}

func LoadConfig() error {
	// Load Puppet version from config
	pupperSemVer, err := semver.NewVersion(viper.GetString(PuppetVerCfgKey))
	if err != nil {
		return fmt.Errorf("could not load '%s' from config '%s': %s", PuppetVerCfgKey, viper.GetViper().ConfigFileUsed(), err)
	}

	RunningConfig.PuppetVersion = pupperSemVer

	// Load Backend from config
	RunningConfig.Backend = BackendType(viper.GetString(BackendCfgKey))

	// Load ToolPath from config
	RunningConfig.ToolPath = viper.GetString(ToolPathCfgKey)

	return nil
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
