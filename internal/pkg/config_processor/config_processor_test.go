package config_processor_test

import (
	"path/filepath"
	"regexp"
	"testing"

	"github.com/puppetlabs/prm/internal/pkg/config_processor"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type ProcessConfigTest struct {
	name     string
	args     args
	expected expected
	mocks    mocks
}

type args struct {
	targetDir string
	sourceDir string
	force     bool
}

type expected struct {
	errorMsg       string
	namespacedPath string
}

type mocks struct {
	dirs  []string
	files map[string]string
}

func TestProcessConfig(t *testing.T) {
	configDir := "path/to/config"

	tests := []ProcessConfigTest{
		{
			name:     "Config file is present and is correctly constructed",
			args:     args{targetDir: "tools", sourceDir: configDir, force: false},
			expected: expected{errorMsg: "", namespacedPath: filepath.Join("tools", "test-user/test-plugin/0.1.0")},
			mocks: mocks{
				dirs: []string{"tools"},
				files: map[string]string{
					filepath.Join(configDir, "prm-config.yml"): `---
plugin:
  id: test-plugin
  author: test-user
  version: 0.1.0
`,
				},
			},
		},
		{
			name:     "Config file does not exist",
			args:     args{targetDir: "tools", sourceDir: configDir, force: false},
			expected: expected{errorMsg: "Invalid config: "},
		},
		{
			name:     "Config files exists but has invalid yaml",
			args:     args{targetDir: "tools", sourceDir: configDir, force: false},
			expected: expected{errorMsg: "Invalid config: "},
			mocks: mocks{
				dirs: []string{"tools"},
				files: map[string]string{
					filepath.Join(configDir, "prm-config.yml"): `---
		plugin: id: test-plugin author: test-user version: 0.1.0
		`,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}

			for _, path := range tt.mocks.dirs {
				afs.Mkdir(path, 0750) //nolint:gosec,errcheck // this result is not used in a secure application
			}

			for file, content := range tt.mocks.files {
				config, _ := afs.Create(file) //nolint:gosec,errcheck // this result is not used in a secure application
				config.Write([]byte(content)) //nolint:errcheck
			}

			configProcessor := config_processor.ConfigProcessor{AFS: afs}

			returnedPath, err := configProcessor.ProcessConfig(tt.args.sourceDir, tt.args.targetDir, tt.args.force)

			if tt.expected.errorMsg != "" && err != nil {
				assert.Contains(t, err.Error(), tt.expected.errorMsg)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected.namespacedPath, returnedPath)
		})
	}

}

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
  id: test-tool
  author: test-user
  version: 0.1.0
  upstream_project_url: https://github.com/test/test-tool
  display: My Test Tool
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
  id: test-tool
  version: 0.1.0
  upstream_project_url: https://github.com/test/test-tool
  display: My Test Tool
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
  upstream_project_url: https://github.com/test/test-tool
  display: My Test Tool
`,
			errorMsg: `The following attributes are missing in .+:\s+\* id`,
		},
		{
			name:           "When config missing version",
			mockConfigFile: true,
			configFilePath: "my/missing/version/prm-config.yml",

			configFileYaml: `---
plugin:
  id: test-tool
  author: test-user
  upstream_project_url: https://github.com/test/test-tool
  display: My Test Tool
`,
			errorMsg: `The following attributes are missing in .+:\s+\* version`,
		},
		{
			name:           "When config missing upstream project url",
			mockConfigFile: true,
			configFilePath: "my/missing/version/prm-config.yml",

			configFileYaml: `---
plugin:
  id: test-tool
  author: test-user
  version: 0.1.0
  display: My Test Tool
`,
			errorMsg: `The following attributes are missing in .+:\s+\* upstream project url`,
		},
		{
			name:           "When config missing display",
			mockConfigFile: true,
			configFilePath: "my/missing/version/prm-config.yml",

			configFileYaml: `---
plugin:
  id: test-tool
  author: test-user
  version: 0.1.0
  upstream_project_url: https://github.com/test/test-tool
`,
			errorMsg: `The following attributes are missing in .+:\s+\* display name`,
		},
		{
			name:           "Config missing all required keys",
			mockConfigFile: true,
			configFilePath: "my/missing/version/prm-config.yml",

			configFileYaml: `---
foo: bar
`,
			errorMsg: `The following attributes are missing in .+:\s+\* id\s+\* author\s+\* version\s+\* upstream project url\s+\* display name`,
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
