package prm_test

import (
	"fmt"
	"testing"

	"github.com/puppetlabs/prm/pkg/prm"
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
			expectedPuppetVersion: "7.0.0",
			expectedBackend:       string(prm.DOCKER),
			expectedToolPath:      "tools",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prmObj := &prm.Prm{}
			prmObj.GenerateDefaultCfg()
			assert.Equal(t, tt.expectedPuppetVersion, viper.GetString(prm.PuppetVerCfgKey))
			assert.Equal(t, tt.expectedBackend, viper.Get(prm.BackendCfgKey))
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
			expectedErrMsg: fmt.Sprintf("could not load '%s' from config '%s': Invalid Semantic Version", prm.PuppetVerCfgKey, viper.GetViper().ConfigFileUsed()),
		},
		{
			name:                "Should error when invalid semver returned for Puppet ver",
			expectedErrMsg:      fmt.Sprintf("could not load '%s' from config '%s': Invalid Semantic Version", prm.PuppetVerCfgKey, viper.GetViper().ConfigFileUsed()),
			configuredPuppetVer: "foo.bar",
		},
		{
			name:                "Should not error when valid semver returned for Puppet ver",
			configuredPuppetVer: "7.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.SetDefault(prm.PuppetVerCfgKey, tt.configuredPuppetVer)

			prmObj := &prm.Prm{}
			err := prmObj.LoadConfig()

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
