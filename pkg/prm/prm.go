//nolint:structcheck,unused
package prm

import "github.com/Masterminds/semver"

type Prm struct {
	codeDir string
	cache   []toolCache
}

type toolCache struct {
	toolName string
	tool     *Tool
}
type PuppetVersion struct {
	version semver.Version
}

// Given a list of tool names, check if these are groups, and return
// an expanded list containing all the toolNames
func (*Prm) checkGroups(tools []string) []string {
	// TODO
	return []string{}
}

// Look within codeDir for a "validate.yml" containing
// a list of tools and/or tool groups that should be run against
// code within codeDir.
func (*Prm) checkLocalConfig() []string {
	// TODO
	return []string{}
}

// Check to see if the requested tool can be found installed.
// If installed read the tool configuration and return
func (*Prm) isToolAvailable(tool string) (Tool, bool) {
	return Tool{}, false
}

// Check to see if the tool is ready to execute
func (*Prm) isToolReady(tool *Tool) bool {
	return false
}

// save traversing to the filesystem
func (*Prm) cacheTool(tool *Tool) error {
	// TODO
	return nil
}

// What version of Puppet is requested by the user
func (*Prm) getPuppetVersion() PuppetVersion {
	return PuppetVersion{}
}
