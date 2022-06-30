package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"

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

func FormatVersion(version, buildDate string, commit string) string {
	version = strings.TrimSpace(strings.TrimPrefix(version, "v"))

	var dateStr string
	if buildDate != "" {
		t, _ := time.Parse(time.RFC3339, buildDate)
		dateStr = t.Format("2006/01/02")
	}

	if commit != "" && len(commit) > 7 {
		length := len(commit) - 7
		commit = strings.TrimSpace(commit[:len(commit)-length])
	}

	return fmt.Sprintf("prm %s %s %s\n\n%s",
		version, commit, dateStr, changelogURL(version))
}

func changelogURL(version string) string {
	path := "https://github.com/puppetlabs/prm"
	r := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[\w.]+)?$`)
	if !r.MatchString(version) {
		return fmt.Sprintf("%s/releases/latest", path)
	}

	url := fmt.Sprintf("%s/releases/tag/%s", path, strings.TrimPrefix(version, "v"))
	return url
}
