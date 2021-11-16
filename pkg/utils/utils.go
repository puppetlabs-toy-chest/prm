package utils

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func WriteConfig() {
	if err := viper.WriteConfig(); err != nil {
		log.Error().Msgf("could not write config to %s: %s", viper.ConfigFileUsed(), err)
	}
}
