package prm

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	PuppetCmdFlag    string = "puppet"
	PuppetVerCfgKey  string = "puppet.version"
	DefaultPuppetVer string = "7"
	BackendCfgKey   string = "backend.type"
)

type Config struct {
	PuppetVersion *semver.Version
	Backend       BackendI
}

var RunningConfig Config

func LoadConfig() error {
	puppetVer := viper.GetString(PuppetVerCfgKey)

	// Set a default Puppet version if it's unset in config
	if puppetVer == "" {
		log.Debug().Msgf("'%s' unset in %s, setting default value: %s", PuppetVerCfgKey, viper.GetViper().ConfigFileUsed(), DefaultPuppetVer)
		puppetVer = DefaultPuppetVer
	}

	puppetSemVer, err := semver.NewVersion(puppetVer)

	if err != nil {
		return fmt.Errorf("Value for '%s' in config is not a valid Puppet semver: %s", PuppetVerCfgKey, err)
	}

	RunningConfig.PuppetVersion = puppetSemVer

	return nil
}
