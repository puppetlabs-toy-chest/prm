package utils

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type UtilsI interface {
	SetAndWriteConfig(string, string) error
}

type Utils struct{}

func (u *Utils) SetAndWriteConfig(k, v string) (err error) {
	log.Trace().Msgf("Setting and saving config '%s' to '%s' in %s", k, v, viper.ConfigFileUsed())

	viper.Set(k, v)

	if err = viper.WriteConfig(); err != nil {
		log.Error().Msgf("could not write config to %s: %s", viper.ConfigFileUsed(), err)
	}
	return err
}
