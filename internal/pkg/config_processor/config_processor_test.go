package config_processor_test

import (
	"path/filepath"
	"reflect"
	"regexp"
	"testing"

	cfg_iface "github.com/puppetlabs/pct/pkg/config_processor"
	"github.com/puppetlabs/prm/internal/pkg/config_processor"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type CheckConfigTest struct {
	name           string
	mockConfigFile bool
	configFilePath string
	configFileYaml string
	errorMsg       string
}

func TestPrmConfigProcessor_CheckConfig(t *testing.T) {
	tests := []CheckConfigTest{
		{
			name:           "When config not found",
			mockConfigFile: false,
			configFilePath: "my/missing/prm-config.yml",
			errorMsg:       "file does not exist",
		},
		{
			name:           "When config valid",
			mockConfigFile: true,
			configFilePath: "my/valid/prm-config.yml",

			configFileYaml: `---
plugin:
  id: test-plugin
  author: test-user
  version: 0.1.0
`,
			errorMsg: "",
		},
		{
			name:           "When config invalid",
			mockConfigFile: true,
			configFilePath: "my/invalid/prm-config.yml",
			// This is invalid because it starts with tabs which the parses errors on
			configFileYaml: `---
			foo: bar
			`,
			errorMsg: "parsing config: yaml",
		},
		{
			name:           "When config missing author",
			mockConfigFile: true,
			configFilePath: "my/missing/author/prm-config.yml",

			configFileYaml: `---
plugin:
  id: test-plugin
  version: 0.1.0
`,
			errorMsg: `The following attributes are missing in .+:\s+\* author`,
		},
		{
			name:           "When config missing id",
			mockConfigFile: true,
			configFilePath: "my/missing/id/prm-config.yml",

			configFileYaml: `---
plugin:
  author: test-user
  version: 0.1.0
`,
			errorMsg: `The following attributes are missing in .+:\s+\* id`,
		},
		{
			name:           "When config missing version",
			mockConfigFile: true,
			configFilePath: "my/missing/version/prm-config.yml",

			configFileYaml: `---
plugin:
  author: test-user
  id: test-plugin
`,
			errorMsg: `The following attributes are missing in .+:\s+\* version`,
		},
		{
			name:           "When config missing author, id, and version",
			mockConfigFile: true,
			configFilePath: "my/missing/version/prm-config.yml",

			configFileYaml: `---
plugin:
`,
			errorMsg: `The following attributes are missing in .+:\s+\* id\s+\* author\s+\* version`,
		},
		{
			name:           "When config missing plugin key",
			mockConfigFile: true,
			configFilePath: "my/missing/version/prm-config.yml",

			configFileYaml: `---
foo: bar
`,
			errorMsg: `The following attributes are missing in .+:\s+\* id\s+\* author\s+\* version`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}

			if tt.mockConfigFile {
				dir := filepath.Dir(tt.configFilePath)
				afs.Mkdir(dir, 0750)                       //nolint:gosec,errcheck // this result is not used in a secure application
				config, _ := afs.Create(tt.configFilePath) //nolint:gosec,errcheck // this result is not used in a secure application
				config.Write([]byte(tt.configFileYaml))    //nolint:errcheck
			}

			configProcessor := config_processor.ConfigProcessor{AFS: afs}

			err := configProcessor.CheckConfig(tt.configFilePath)

			if tt.errorMsg != "" && err != nil {
				assert.Regexp(t, regexp.MustCompile(tt.errorMsg), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPrmConfigProcessor_GetConfigMetadata(t *testing.T) {
	type args struct {
		configFile string
	}
	configParentPath := "path/to/extract/to/"

	tests := []struct {
		name         string
		args         args
		wantMetadata cfg_iface.ConfigMetadata
		wantErr      bool
		pluginConfig string // Leave blank for config file not to be created
	}{
		{
			name: "Successfully gets config metadata",
			args: args{
				configFile: filepath.Join(configParentPath, "prm-config.yml"),
			},
			wantMetadata: cfg_iface.ConfigMetadata{Author: "test-user", Id: "full-project", Version: "0.1.0"},
			pluginConfig: `---
plugin:
  id: full-project
  author: test-user
  version: 0.1.0
`,
		},
		{
			name: "Missing vital metadata from prm-config.yml (id omitted)",
			args: args{
				configFile: filepath.Join(configParentPath, "prm-config.yml"),
			},
			wantErr:      true,
			wantMetadata: cfg_iface.ConfigMetadata{},
			pluginConfig: `---
plugin:
  author: test-user
  version: 0.1.0
`,
		},
		{
			name: "Malformed prm-config (extra indentation)",
			args: args{
				configFile: filepath.Join(configParentPath, "prm-config.yml"),
			},
			wantErr:      true,
			wantMetadata: cfg_iface.ConfigMetadata{},
			pluginConfig: `---
	plugin:
		id: full-project
  	author: test-user
  	version: 0.1.0
`, // Contains an erroneous extra indent
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Instantiate afs
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			p := &config_processor.ConfigProcessor{
				AFS: afs,
			}

			// Create all useful directories
			afs.MkdirAll(configParentPath, 0750) //nolint:gosec,errcheck
			if tt.pluginConfig != "" {
				config, _ := afs.Create(tt.args.configFile)
				config.Write([]byte(tt.pluginConfig)) //nolint:errcheck
			}

			gotMetadata, err := p.GetConfigMetadata(tt.args.configFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfigMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMetadata, tt.wantMetadata) {
				t.Errorf("GetConfigMetadata() gotMetadata = %v, want %v", gotMetadata, tt.wantMetadata)
			}
		})
	}
}
