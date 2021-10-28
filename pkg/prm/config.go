//nolint:structcheck,unused
package prm

import "github.com/Masterminds/semver"

type Config struct {
	puppetVersion semver.Version
	backend       BackendI
}
