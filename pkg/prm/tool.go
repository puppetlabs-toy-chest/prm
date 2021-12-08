// nolint:structcheck,unused
package prm

import (
	"io"

	"github.com/puppetlabs/pdkgo/pkg/install"
)

type Tool struct {
	Stdout   io.Reader
	Stderr   io.Reader
	ExitCode ToolExitCode
	Cfg      ToolConfig
}

type ToolExitCode int64

const (
	FAILURE ToolExitCode = iota
	SUCCESS
	TOOL_ERROR
	TOOL_NOT_FOUND
)

type ToolConfig struct {
	Path      string
	Plugin    *PluginConfig    `mapstructure:"plugin"`
	Gem       *GemConfig       `mapstructure:"gem"`
	Container *ContainerConfig `mapstructure:"container"`
	Binary    *BinaryConfig    `mapstructure:"binary"`
	Puppet    *PuppetConfig    `mapstructure:"puppet"`
	Common    CommonConfig     `mapstructure:"common"`
}

// ToolConfigInfo is the housing struct for marshaling YAML data
type ToolConfigInfo struct {
	Plugin   PluginConfig `mapstructure:"plugin"`
	Defaults map[string]interface{}
}

type PluginConfig struct {
	install.ConfigParams `mapstructure:",squash"`
	Display              string `mapstructure:"display"`
	UpstreamProjUrl      string `mapstructure:"upstream_project_url"`
}

type BinaryConfig struct {
	Name         string        `mapstructure:"name"`
	InstallSteps *InstallSteps `mapstructure:"install_steps"`
}

type InstallSteps struct {
	Windows string `mapstructure:"windows"`
	Darwin  string `mapstructure:"darwin"`
	Linux   string `mapstructure:"linux"`
}

type ContainerConfig struct {
	Name string `mapstructure:"name"`
	Tag  string `mapstructure:"tag"`
}

type GemConfig struct {
	Name          []string             `mapstructure:"name"`
	Executable    string               `mapstructure:"executable"`
	BuildTools    bool                 `mapstructure:"build_tools"`
	Compatibility map[float32][]string `mapstructure:"compatibility"`
}

type PuppetConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

type CommonConfig struct {
	CanValidate         bool              `mapstructure:"can_validate"`
	NeedsWriteAccess    bool              `mapstructure:"needs_write_access"`
	UseScript           string            `mapstructure:"use_script"`
	RequiresGit         bool              `mapstructure:"requires_git"`
	DefaultArgs         []string          `mapstructure:"default_args"`
	HelpArg             string            `mapstructure:"help_arg"`
	SuccessExitCode     int               `mapstructure:"success_exit_code"`
	InterleaveStdOutErr bool              `mapstructure:"interleave_stdout"`
	OutputMode          *OutputModes      `mapstructure:"output_mode"`
	Env                 map[string]string `mapstructure:"env"`
}

type OutputModes struct {
	Json  string `mapstructure:"json"`
	Yaml  string `mapstructure:"yaml"`
	Junit string `mapstructure:"junit"`
}
