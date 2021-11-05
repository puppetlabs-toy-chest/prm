//nolint:structcheck,unused
package prm

import (
	"github.com/Masterminds/semver"
)

const (
	PuppetCmdFlag   string = "puppet"
	PuppetVerCfgKey string = "puppet.version"
)

type Config struct {
	puppetVersion semver.Version
	backend       BackendI
}
