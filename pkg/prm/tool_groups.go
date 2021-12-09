package prm

/*
  This package contains ToolGroups, which defines
  what inidividual tools make up a specific ToolGroup.

  This allows users to specify 'group/mygroup' and have PRM
  execute against a larger list of tools, without needing to define
  or understand what is being called.
*/

var (
	ToolGroups = map[string][]string{
		// TODO: we may need to define group as a reserved word
		"group/modules": {
			"puppetlabs/rubocop",
			"puppetlabs/rspec-puppet",
			"puppetlabs/puppet-lint",
			"puppetlabs/puppet-syntax",
			"puppetlabs/puppet-strings",
		},
	}
)
