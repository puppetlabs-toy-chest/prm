package prm

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/spf13/viper"
)

const (
	PuppetCmdFlag   string = "puppet"
	BackendCmdFlag  string = "backend"
	PuppetVerCfgKey string = "puppetversion" // Should match Config struct key.
	BackendCfgKey   string = "backend"       // Should match Config struct key.
)

type Config struct {
	PuppetVersion *semver.Version
	Backend       BackendType
}

var RunningConfig Config

func GenerateDefaultCfg() {
	// Generate default configuration
	puppetVer, err := semver.NewVersion("7")
	if err != nil {
		panic(fmt.Sprintf("Unable to generate default cfg value for 'puppet': %s", err))
	}

	viper.SetDefault(PuppetVerCfgKey, puppetVer)
	viper.SetDefault(BackendCfgKey, DOCKER)
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

	return nil
}
