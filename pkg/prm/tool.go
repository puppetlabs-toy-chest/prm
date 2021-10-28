//nolint:structcheck,unused
package prm

import (
	"io"

	"github.com/Masterminds/semver"
)

type Tool struct {
	stdout   io.Reader
	stderr   io.Reader
	exitCode ToolExitCode
	cfg      ToolConfig
}

type ToolExitCode int64

const (
	SUCCESS ToolExitCode = iota
	FAILURE
	TOOL_ERROR
	TOOL_NOT_FOUND
)

type ToolConfig struct {
	version         semver.Version
	author          string
	id              string
	name            string
	display         string
	upstreamProjUrl string
	binaryConfig    BinaryConfig
	commonConfig    CommonConfig
	containerConfig ContainerConfig
	gemConfig       GemConfig
	puppetConfig    PuppetConfig
}

type BinaryConfig struct {
	name    string
	windows string
	darwin  string
	linux   string
}

type CommonConfig struct {
	canValidate         bool
	needsWriteAccess    bool
	defaultArgs         []string
	helpArg             string
	successExitCode     int
	interleaveStdOutErr bool
	jsonOutputFlag      string
	junitOutputFlag     string
	yamlOutputFlag      string
}

type ContainerConfig struct {
	name string
	tag  string
}

type GemConfig struct {
	name       []string
	executable string
	buildTools bool
	rakefile   string
}

type PuppetConfig struct {
	enabled bool
}
