package tool

/*
  This package contains ToolGroups, which defines
  what inidividual tools make up a specific ToolGroup.

  This allows users to specify 'group/mygroup' and have PRM
  execute against a larger list of tools, without needing to define
  or understand what is being called.
*/
type ToolInst struct {
	Name string   `yaml:"name"`
	Args []string `yaml:"args"`
}

var (
	ToolGroups = map[string][]ToolInst{
		// TODO: we may need to define group as a reserved word
		"group/modules": {
			{Name: "puppetlabs/spec_cache"},
			{
				Name: "puppetlabs/spec_puppet",
				Args: []string{
					"spec_prep",
				},
			},
			{Name: "puppetlabs/spec_puppet"},
			{Name: "puppetlabs/rubocop"},
			{Name: "puppetlabs/puppet-lint"},
			{Name: "puppetlabs/puppet-syntax"},
			{Name: "puppetlabs/puppet-strings"},
		},
	}
)

func compareToolInst(t1 ToolInst, t2 ToolInst) bool {
	return t1.Name == t2.Name && equal(t1.Args, t2.Args)
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
