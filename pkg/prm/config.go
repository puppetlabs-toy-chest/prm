//nolint:structcheck,unused
package prm

import (
	"github.com/Masterminds/semver"
)

const (
	PuppetVerCfgKey string = "puppet.version"
)

type Config struct {
	puppetVersion semver.Version
	backend       BackendI
}
