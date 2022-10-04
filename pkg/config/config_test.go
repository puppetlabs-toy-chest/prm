package config_test

import (
	"fmt"
	"github.com/puppetlabs/prm/pkg/config"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGenerateDefaultCfg(t *testing.T) {
	tests := []struct {
		name                  string
		expectedPuppetVersion string
		expectedBackend       string
		expectedToolPath      string
	}{
		{
			name:                  "Should generate default Puppet and Backend cfgs",
			expectedPuppetVersion: "7.15.0",
			expectedBackend:       string(config.DOCKER),
			expectedToolPath:      "tools",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.GenerateDefaultCfg()
			assert.Equal(t, tt.expectedPuppetVersion, viper.GetString(config.PuppetVerCfgKey))
			assert.Equal(t, tt.expectedBackend, viper.Get(config.BackendCfgKey))
		})
	}
}

// To test unlikely error condition that a garbage or nil version has made it
// in as the configured Puppet version
func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name                string
		expectedErrMsg      string
		configuredPuppetVer string
	}{
		{
			name:           "Should error when nil returned for Puppet ver",
			expectedErrMsg: fmt.Sprintf("could not load '%s' from config '%s': Invalid Semantic Version", config.PuppetVerCfgKey, viper.GetViper().ConfigFileUsed()),
		},
		{
			name:                "Should error when invalid semver returned for Puppet ver",
			expectedErrMsg:      fmt.Sprintf("could not load '%s' from config '%s': Invalid Semantic Version", config.PuppetVerCfgKey, viper.GetViper().ConfigFileUsed()),
			configuredPuppetVer: "foo.bar",
		},
		{
			name:                "Should not error when valid semver returned for Puppet ver",
			configuredPuppetVer: "7.15.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.SetDefault(config.PuppetVerCfgKey, tt.configuredPuppetVer)

			_, err := config.LoadConfig()
			if tt.expectedErrMsg != "" && err != nil {
				assert.Contains(t, tt.expectedErrMsg, err.Error())
				return
			}

			if tt.expectedErrMsg == "" && err != nil {
				t.Errorf("LoadConfig() Unexpected error: %s", err)
				return
			}
		})
	}
}
