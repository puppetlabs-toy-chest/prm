package prm_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pdkgo/pkg/install"
	"github.com/puppetlabs/prm/pkg/prm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// hide logging output
	log.Logger = zerolog.New(ioutil.Discard).With().Timestamp().Logger()
	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	type args struct {
		toolDirPath string
		setup       bool
		toolConfig  string
	}
	tests := []struct {
		name        string
		args        args
		want        prm.Tool
		wantErr     bool
		expectedErr string
	}{
		{
			name: "returns error for non-existent tool",
			args: args{
				toolDirPath: "tools/author/i-dont-exist/0.1.0",
				setup:       false,
			},
			wantErr:     true,
			expectedErr: "Couldn't find an installed tool at 'tools/author/i-dont-exist/0.1.0'",
		},
		{
			name: "returns config for existent tool",
			args: args{
				toolDirPath: "tools/author/jeans/0.1.0",
				setup:       true,
				toolConfig: `---
plugin:
  id: jeans
  author: JoeBloggs
  display: Jeans
  version: 0.1.0
  upstream_project_url: https://github.com/joebloggs/prm-jeans
`,
			},
			want: prm.Tool{
				Cfg: prm.ToolConfig{
					Plugin: &prm.PluginConfig{
						ConfigParams: install.ConfigParams{
							Id:      "jeans",
							Author:  "JoeBloggs",
							Version: "0.1.0",
						},
						Display:         "Jeans",
						UpstreamProjUrl: "https://github.com/joebloggs/prm-jeans",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "returns config for existent tool with a malformed config",
			args: args{
				toolDirPath: "tools/author/dud/0.1.0",
				setup:       true,
				toolConfig:  `<xml></xml>`,
			},
			want: prm.Tool{
				Cfg: prm.ToolConfig{},
			},
			wantErr:     true,
			expectedErr: fmt.Sprintf("Couldn't parse tool config at '%s'", filepath.Join("tools/author/dud/0.1.0", "prm-config.yml")),
		},
		{
			name: "returns config for existent GEM tool",
			args: args{
				toolDirPath: "tools/author/jeans/0.1.0",
				setup:       true,
				toolConfig: `---
plugin:
  id: jeans
  author: JoeBloggs
  display: Jeans
  version: 0.1.0
  upstream_project_url: https://github.com/joebloggs/prm-jeans

gem:
  name: ['prmjeans', 'jeans-belt']
  executable: jeans
  compatibility:
    - 2.4: ['~> 0.1.0']
    - 2.5: ['>= 1.3.2', '<= 1.5.7']
`,
			},
			want: prm.Tool{
				Cfg: prm.ToolConfig{
					Plugin: &prm.PluginConfig{
						ConfigParams: install.ConfigParams{
							Id:      "jeans",
							Author:  "JoeBloggs",
							Version: "0.1.0",
						},
						Display:         "Jeans",
						UpstreamProjUrl: "https://github.com/joebloggs/prm-jeans",
					},
					Gem: &prm.GemConfig{
						Name:       []string{"prmjeans", "jeans-belt"},
						Executable: "jeans",
						BuildTools: false,
						Compatibility: map[float32][]string{
							2.4: {"~> 0.1.0"},
							2.5: {">= 1.3.2", "<= 1.5.7"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "returns config for existent CONTAINER tool",
			args: args{
				toolDirPath: "tools/author/jeans/0.1.0",
				setup:       true,
				toolConfig: `---
plugin:
  id: jeans
  author: JoeBloggs
  display: Jeans
  version: 0.1.0
  upstream_project_url: https://github.com/joebloggs/prm-jeans

container:
  name: 'prmjeans'
  tag: 'latest'
`,
			},
			want: prm.Tool{
				Cfg: prm.ToolConfig{
					Plugin: &prm.PluginConfig{
						ConfigParams: install.ConfigParams{
							Id:      "jeans",
							Author:  "JoeBloggs",
							Version: "0.1.0",
						},
						Display:         "Jeans",
						UpstreamProjUrl: "https://github.com/joebloggs/prm-jeans",
					},
					Container: &prm.ContainerConfig{
						Name: "prmjeans",
						Tag:  "latest",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "returns config for existent BINARY tool",
			args: args{
				toolDirPath: "tools/author/jeans/0.1.0",
				setup:       true,
				toolConfig: `---
plugin:
  id: jeans
  author: JoeBloggs
  display: Jeans
  version: 0.1.0
  upstream_project_url: https://github.com/joebloggs/prm-jeans

binary:
  name: 'prmjeans'
  install_steps:
    windows: |
      curl http://github.com/joebloggs/prm-jeans/raw/master/bin/windows/prmjeans.exe -o prmjeans.exe
      ./prmjeans.exe
    linux: |
      curl http://github.com/joebloggs/prm-jeans/raw/master/bin/linux/prmjeans -o prmjeans
      ./prmjeans
    darwin: |
      curl http://github.com/joebloggs/prm-jeans/raw/master/bin/darwin/prmjeans -o prmjeans
      ./prmjeans
`,
			},
			want: prm.Tool{
				Cfg: prm.ToolConfig{
					Plugin: &prm.PluginConfig{
						ConfigParams: install.ConfigParams{
							Id:      "jeans",
							Author:  "JoeBloggs",
							Version: "0.1.0",
						},
						Display:         "Jeans",
						UpstreamProjUrl: "https://github.com/joebloggs/prm-jeans",
					},
					Binary: &prm.BinaryConfig{
						Name: "prmjeans",
						InstallSteps: &prm.InstallSteps{
							Windows: "curl http://github.com/joebloggs/prm-jeans/raw/master/bin/windows/prmjeans.exe -o prmjeans.exe\n./prmjeans.exe\n",
							Linux:   "curl http://github.com/joebloggs/prm-jeans/raw/master/bin/linux/prmjeans -o prmjeans\n./prmjeans\n",
							Darwin:  "curl http://github.com/joebloggs/prm-jeans/raw/master/bin/darwin/prmjeans -o prmjeans\n./prmjeans\n",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "returns config for existent PUPPET tool",
			args: args{
				toolDirPath: "tools/author/jeans/0.1.0",
				setup:       true,
				toolConfig: `---
plugin:
  id: jeans
  author: JoeBloggs
  display: Jeans
  version: 0.1.0
  upstream_project_url: https://github.com/joebloggs/prm-jeans

puppet:
  enabled: true
`,
			},
			want: prm.Tool{
				Cfg: prm.ToolConfig{
					Plugin: &prm.PluginConfig{
						ConfigParams: install.ConfigParams{
							Id:      "jeans",
							Author:  "JoeBloggs",
							Version: "0.1.0",
						},
						Display:         "Jeans",
						UpstreamProjUrl: "https://github.com/joebloggs/prm-jeans",
					},
					Puppet: &prm.PuppetConfig{
						Enabled: true,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "returns config for existent tool with common config items",
			args: args{
				toolDirPath: "tools/author/jeans/0.1.0",
				setup:       true,
				toolConfig: `---
plugin:
  id: jeans
  author: JoeBloggs
  display: Jeans
  version: 0.1.0
  upstream_project_url: https://github.com/joebloggs/prm-jeans

puppet:
  enabled: true

common:
  can_validate: true
  needs_write_access: false
  use_script: "entrypoint"
  requires_git: true
  default_args: ["--verbose"]
  help_arg: "--help"
  success_exit_code: 0
  interleave_stdout: false
  output_mode:
    json: "--output json"
    yaml: "-y"
    junit: "--output xml"
`,
			},
			want: prm.Tool{
				Cfg: prm.ToolConfig{
					Plugin: &prm.PluginConfig{
						ConfigParams: install.ConfigParams{
							Id:      "jeans",
							Author:  "JoeBloggs",
							Version: "0.1.0",
						},
						Display:         "Jeans",
						UpstreamProjUrl: "https://github.com/joebloggs/prm-jeans",
					},
					Puppet: &prm.PuppetConfig{
						Enabled: true,
					},
					Common: prm.CommonConfig{
						CanValidate:         true,
						NeedsWriteAccess:    false,
						UseScript:           "entrypoint",
						RequiresGit:         true,
						DefaultArgs:         []string{"--verbose"},
						HelpArg:             "--help",
						SuccessExitCode:     0,
						InterleaveStdOutErr: false,
						OutputMode: &prm.OutputModes{
							Json:  "--output json",
							Yaml:  "-y",
							Junit: "--output xml",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			iofs := &afero.IOFS{Fs: fs}

			if tt.args.setup {
				// Create tool config
				config, _ := afs.Create(filepath.Join(tt.args.toolDirPath, "prm-config.yml"))
				config.Write([]byte(tt.args.toolConfig)) //nolint:errcheck
			}

			p := &prm.Prm{
				AFS:  afs,
				IOFS: iofs,
			}

			got, err := p.Get(tt.args.toolDirPath)

			if tt.wantErr {
				assert.Equal(t, tt.expectedErr, err.Error())
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormatTools(t *testing.T) {
	type args struct {
		tools      []prm.ToolConfig
		jsonOutput string
	}
	tests := []struct {
		name    string
		args    args
		matches []string
		wantErr bool
	}{
		{
			name: "When no tools are passed",
			args: args{
				tools:      []prm.ToolConfig{},
				jsonOutput: "table",
			},
			matches: []string{},
		},
		{
			name: "When only one tool is passed",
			args: args{
				tools: []prm.ToolConfig{
					{
						Plugin: &prm.PluginConfig{
							ConfigParams: install.ConfigParams{
								Id:      "foo",
								Author:  "bar",
								Version: "0.1.0",
							},
							Display:         "Foo Item",
							UpstreamProjUrl: "https://github.com/bar/pct-foo",
						},
					},
				},
				jsonOutput: "table",
			},
			matches: []string{
				`DisplayName:\s+Foo Item`,
				`Author:\s+bar`,
				`Name:\s+foo`,
				`Project_URL:\s+https://github.com/bar/pct-foo`,
				`Version:\s+0\.1\.0`,
			},
		},
		{
			name: "When more than one tool is passed",
			args: args{
				tools: []prm.ToolConfig{
					{
						Plugin: &prm.PluginConfig{
							ConfigParams: install.ConfigParams{
								Id:      "foo",
								Author:  "baz",
								Version: "0.1.0",
							},
							Display:         "Foo Item",
							UpstreamProjUrl: "https://github.com/baz/pct-foo",
						},
					},
					{
						Plugin: &prm.PluginConfig{
							ConfigParams: install.ConfigParams{
								Id:      "bar",
								Author:  "baz",
								Version: "0.1.0",
							},
							Display:         "Bar Item",
							UpstreamProjUrl: "https://github.com/baz/pct-bar",
						},
					},
				},
				jsonOutput: "table",
			},
			matches: []string{
				`\s+DISPLAYNAME\s+\|\s+AUTHOR\s+\|\s+NAME\s+\|\s+PROJECT URL\s+\|\s+VERSION\s+`,
				`Foo Item\s+\|\sbaz\s+\|\sfoo\s+\|\shttps:\/\/github.com\/baz\/pct-foo\s+|\s0.1.0`,
				`Bar Item\s+\|\sbaz\s+\|\sbar\s+\|\shttps:\/\/github.com\/baz\/pct-bar\s+|\s0.1.0`,
			},
		},
		{
			name: "When format is specified as json",
			args: args{
				tools: []prm.ToolConfig{
					{
						Plugin: &prm.PluginConfig{
							ConfigParams: install.ConfigParams{
								Id:      "foo",
								Author:  "baz",
								Version: "0.1.0",
							},
							Display:         "Foo Item",
							UpstreamProjUrl: "https://github.com/baz/pct-foo",
						},
					},
					{
						Plugin: &prm.PluginConfig{
							ConfigParams: install.ConfigParams{
								Id:      "bar",
								Author:  "baz",
								Version: "0.1.0",
							},
							Display:         "Bar Item",
							UpstreamProjUrl: "https://github.com/baz/pct-bar",
						},
					},
				},
				jsonOutput: "json",
			},
			matches: []string{
				`\"Id\": \"foo\"`,
				`\"Id\": \"bar\"`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			iofs := &afero.IOFS{Fs: fs}

			p := &prm.Prm{
				AFS:  afs,
				IOFS: iofs,
			}

			output, err := p.FormatTools(tt.args.tools, tt.args.jsonOutput)
			if tt.wantErr {
				assert.Error(t, err)
			}
			for _, m := range tt.matches {
				assert.Regexp(t, m, output)
			}
		})
	}
}

func TestList(t *testing.T) {
	type stubbedConfig struct {
		relativeConfigPath string
		configContent      string
	}
	type args struct {
		toolPath       string
		toolName       string
		stubbedConfigs []stubbedConfig
	}
	tests := []struct {
		name string
		args args
		want []prm.ToolConfig
	}{
		{
			name: "when no tools are found",
			args: args{
				toolPath: "stubbed/tools/none",
			},
		},
		{
			name: "when an invalid tool config is found",
			args: args{
				toolPath: "stubbed/tools/invalid",
				stubbedConfigs: []stubbedConfig{
					{
						relativeConfigPath: "some_author/bad-tool/0.1.0",
						configContent:      "I am WILDLY INVALID",
					},
				},
			},
		},
		{
			name: "when valid tool configs are found",
			args: args{
				toolPath: "stubbed/tools/valid",
				stubbedConfigs: []stubbedConfig{
					{
						relativeConfigPath: "some_author/first/0.1.0",
						configContent: `---
plugin:
  author: some_author
  id: first
  display: First Tool
  version: 0.1.0
  upstream_project_url: https://github.com/some_author/pct-first-tool
`,
					},
					{
						relativeConfigPath: "some_author/second/0.1.0",
						configContent: `---
plugin:
  author: some_author
  id: second
  display: Second Tool
  version: 0.1.0
  upstream_project_url: https://github.com/some_author/pct-second-tool
`,
					},
				},
			},
			want: []prm.ToolConfig{
				{
					Plugin: &prm.PluginConfig{
						ConfigParams: install.ConfigParams{
							Author:  "some_author",
							Id:      "first",
							Version: "0.1.0",
						},
						Display:         "First Tool",
						UpstreamProjUrl: "https://github.com/some_author/pct-first-tool",
					},
				},
				{
					Plugin: &prm.PluginConfig{
						ConfigParams: install.ConfigParams{
							Author:  "some_author",
							Id:      "second",
							Version: "0.1.0",
						},
						Display:         "Second Tool",
						UpstreamProjUrl: "https://github.com/some_author/pct-second-tool",
					},
				},
			},
		},
		{
			name: "when tools are found with the same author/id and different versions",
			args: args{
				toolPath: "stubbed/tools/multiversion",
				stubbedConfigs: []stubbedConfig{
					{
						relativeConfigPath: "some_author/first/0.1.0",
						configContent: `---
plugin:
  author: some_author
  id: first
  display: First Tool
  version: 0.1.0
  upstream_project_url: https://github.com/some_author/pct-first-tool
`,
					},
					{
						relativeConfigPath: "some_author/first/0.2.0",
						configContent: `---
plugin:
  author: some_author
  id: first
  display: First Tool
  version: 0.2.0
  upstream_project_url: https://github.com/some_author/pct-first-tool
`,
					},
				},
			},
			want: []prm.ToolConfig{
				{
					Plugin: &prm.PluginConfig{
						ConfigParams: install.ConfigParams{
							Author:  "some_author",
							Id:      "first",
							Version: "0.2.0",
						},
						Display:         "First Tool",
						UpstreamProjUrl: "https://github.com/some_author/pct-first-tool",
					},
				},
			},
		},
		{
			name: "when tool name is specified",
			args: args{
				toolPath: "stubbed/tools/named",
				toolName: "first",
				stubbedConfigs: []stubbedConfig{
					{
						relativeConfigPath: "some_author/first/0.1.0",
						configContent: `---
plugin:
  author: some_author
  id: first
  display: First Tool
  version: 0.1.0
  upstream_project_url: https://github.com/some_author/pct-first-tool
`,
					},
					{
						relativeConfigPath: "some_author/second/0.1.0",
						configContent: `---
plugin:
  author: some_author
  id: second
  display: Second Tool
  version: 0.1.0
  upstream_project_url: https://github.com/some_author/pct-second-tool
`,
					},
				},
			},
			want: []prm.ToolConfig{
				{
					Plugin: &prm.PluginConfig{
						ConfigParams: install.ConfigParams{
							Author:  "some_author",
							Id:      "first",
							Version: "0.1.0",
						},
						Display:         "First Tool",
						UpstreamProjUrl: "https://github.com/some_author/pct-first-tool",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: fs}
			iofs := &afero.IOFS{Fs: fs}

			for _, st := range tt.args.stubbedConfigs {
				toolDir := filepath.Join(tt.args.toolPath, st.relativeConfigPath)
				afs.MkdirAll(toolDir, 0750) //nolint:errcheck
				// Create tool config
				config, _ := afs.Create(filepath.Join(toolDir, "prm-config.yml"))
				config.Write([]byte(st.configContent)) //nolint:errcheck
			}

			p := &prm.Prm{
				AFS:  afs,
				IOFS: iofs,
			}

			got := p.List(tt.args.toolPath, tt.args.toolName)
			assert.Equal(t, tt.want, got)
		})
	}
}
