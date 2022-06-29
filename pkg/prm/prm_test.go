package prm_test

import (
	"github.com/puppetlabs/prm/pkg/tool"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pct/pkg/install"
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

func TestFormatTools(t *testing.T) {
	type args struct {
		tools      map[string]*tool.Tool
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
				tools:      map[string]*tool.Tool{},
				jsonOutput: "table",
			},
			matches: []string{},
		},
		{
			name: "When only one tool is passed",
			args: args{
				tools: map[string]*tool.Tool{
					"bar/foo": {
						Cfg: tool.ToolConfig{
							Plugin: &tool.PluginConfig{
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
				tools: map[string]*tool.Tool{
					"baz/foo": {
						Cfg: tool.ToolConfig{
							Plugin: &tool.PluginConfig{
								ConfigParams: install.ConfigParams{
									Id:      "foo",
									Author:  "baz",
									Version: "0.1.0",
								},
								Display:         "Foo Item",
								UpstreamProjUrl: "https://github.com/baz/pct-foo",
							},
						},
					},
					"baz/bar": {
						Cfg: tool.ToolConfig{
							Plugin: &tool.PluginConfig{
								ConfigParams: install.ConfigParams{
									Id:      "bar",
									Version: "0.1.0",
									Author:  "baz",
								},
								Display:         "Bar Item",
								UpstreamProjUrl: "https://github.com/baz/pct-bar",
							},
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
				tools: map[string]*tool.Tool{
					"baz/foo": {
						Cfg: tool.ToolConfig{
							Plugin: &tool.PluginConfig{
								ConfigParams: install.ConfigParams{
									Id:      "foo",
									Version: "0.1.0",
									Author:  "baz",
								},
								Display:         "Foo Item",
								UpstreamProjUrl: "https://github.com/baz/pct-foo",
							},
						},
					},
					"baz/bar": {
						Cfg: tool.ToolConfig{
							Plugin: &tool.PluginConfig{
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
		validateOnly   bool
		stubbedConfigs []stubbedConfig
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*tool.Tool
		wantErr bool
	}{
		{
			name: "when no tools are found",
			args: args{
				toolPath: "stubbed/tools/none",
			},
			wantErr: true,
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
			wantErr: true,
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
			want: map[string]*tool.Tool{
				"some_author/first": {
					Cfg: tool.ToolConfig{
						Path: filepath.Join("stubbed/tools/valid/some_author/first/0.1.0"),
						Plugin: &tool.PluginConfig{
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
				"some_author/second": {
					Cfg: tool.ToolConfig{
						Path: filepath.Join("stubbed/tools/valid/some_author/second/0.1.0"),
						Plugin: &tool.PluginConfig{
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
			want: map[string]*tool.Tool{
				"some_author/first": {
					Cfg: tool.ToolConfig{
						Path: filepath.Join("stubbed/tools/multiversion/some_author/first/0.2.0"),
						Plugin: &tool.PluginConfig{
							ConfigParams: install.ConfigParams{
								Author:  "some_author",
								Version: "0.2.0",
								Id:      "first",
							},
							Display:         "First Tool",
							UpstreamProjUrl: "https://github.com/some_author/pct-first-tool",
						},
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
			want: map[string]*tool.Tool{
				"some_author/first": {
					Cfg: tool.ToolConfig{
						Path: filepath.Join("stubbed/tools/named/some_author/first/0.1.0"),
						Plugin: &tool.PluginConfig{
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
		},
		{
			name: "when validate only is specified",
			args: args{
				toolPath:     "stubbed/tools/named",
				toolName:     "first",
				validateOnly: true,
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

common:
  can_validate: true
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

common:
  can_validate: false
`,
					},
				},
			},
			want: map[string]*tool.Tool{
				"some_author/first": {
					Cfg: tool.ToolConfig{
						Path: filepath.Join("stubbed/tools/named/some_author/first/0.1.0"),
						Plugin: &tool.PluginConfig{
							ConfigParams: install.ConfigParams{
								Author:  "some_author",
								Id:      "first",
								Version: "0.1.0",
							},
							Display:         "First Tool",
							UpstreamProjUrl: "https://github.com/some_author/pct-first-tool",
						},
						Common: tool.CommonConfig{CanValidate: true},
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

			err := p.List(tt.args.toolPath, tt.args.toolName, tt.args.validateOnly)
			if (err != nil) != tt.wantErr {
				t.Errorf("Prm.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, p.Cache)
		})
	}
}
