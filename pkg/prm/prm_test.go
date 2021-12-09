package prm_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

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
		tools      map[string]*prm.Tool
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
				tools:      map[string]*prm.Tool{},
				jsonOutput: "table",
			},
			matches: []string{},
		},
		{
			name: "When only one tool is passed",
			args: args{
				tools: map[string]*prm.Tool{
					"bar/foo": {
						Cfg: prm.ToolConfig{
							Plugin: &prm.PluginConfig{
								Id:              "foo",
								Author:          "bar",
								Display:         "Foo Item",
								Version:         "0.1.0",
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
				tools: map[string]*prm.Tool{
					"baz/foo": {
						Cfg: prm.ToolConfig{
							Plugin: &prm.PluginConfig{
								Id:              "foo",
								Author:          "baz",
								Display:         "Foo Item",
								Version:         "0.1.0",
								UpstreamProjUrl: "https://github.com/baz/pct-foo",
							},
						},
					},
					"baz/bar": {
						Cfg: prm.ToolConfig{
							Plugin: &prm.PluginConfig{
								Id:              "bar",
								Author:          "baz",
								Display:         "Bar Item",
								Version:         "0.1.0",
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
				tools: map[string]*prm.Tool{
					"baz/foo": {
						Cfg: prm.ToolConfig{
							Plugin: &prm.PluginConfig{
								Id:              "foo",
								Author:          "baz",
								Display:         "Foo Item",
								Version:         "0.1.0",
								UpstreamProjUrl: "https://github.com/baz/pct-foo",
							},
						},
					},
					"baz/bar": {
						Cfg: prm.ToolConfig{
							Plugin: &prm.PluginConfig{
								Id:              "bar",
								Author:          "baz",
								Display:         "Bar Item",
								Version:         "0.1.0",
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
		stubbedConfigs []stubbedConfig
	}
	tests := []struct {
		name string
		args args
		want map[string]*prm.Tool
	}{
		{
			name: "when no tools are found",
			args: args{
				toolPath: "stubbed/tools/none",
			},
			want: map[string]*prm.Tool{},
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
			want: map[string]*prm.Tool{},
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
			want: map[string]*prm.Tool{
				"some_author/first": {
					Cfg: prm.ToolConfig{
						Path: filepath.Join("stubbed/tools/valid/some_author/first/0.1.0"),
						Plugin: &prm.PluginConfig{
							Author:          "some_author",
							Id:              "first",
							Display:         "First Tool",
							Version:         "0.1.0",
							UpstreamProjUrl: "https://github.com/some_author/pct-first-tool",
						},
					},
				},
				"some_author/second": {
					Cfg: prm.ToolConfig{
						Path: filepath.Join("stubbed/tools/valid/some_author/second/0.1.0"),
						Plugin: &prm.PluginConfig{
							Author:          "some_author",
							Id:              "second",
							Display:         "Second Tool",
							Version:         "0.1.0",
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
			want: map[string]*prm.Tool{
				"some_author/first": {
					Cfg: prm.ToolConfig{
						Path: filepath.Join("stubbed/tools/multiversion/some_author/first/0.2.0"),
						Plugin: &prm.PluginConfig{
							Author:          "some_author",
							Id:              "first",
							Display:         "First Tool",
							Version:         "0.2.0",
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
			want: map[string]*prm.Tool{
				"some_author/first": {
					Cfg: prm.ToolConfig{
						Path: filepath.Join("stubbed/tools/named/some_author/first/0.1.0"),
						Plugin: &prm.PluginConfig{
							Author:          "some_author",
							Id:              "first",
							Display:         "First Tool",
							Version:         "0.1.0",
							UpstreamProjUrl: "https://github.com/some_author/pct-first-tool",
						},
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

			p.List(tt.args.toolPath, tt.args.toolName)
			assert.Equal(t, tt.want, p.Cache)
		})
	}
}
