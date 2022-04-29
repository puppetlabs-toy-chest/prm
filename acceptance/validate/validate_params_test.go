package validate_test

import (
	"fmt"
	"github.com/puppetlabs/pct/acceptance/testutils"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_PrmValidate_List_Flag_Valid_ToolDir(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	skipValidationTests(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/validate/testdata/tooldir")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("validate --toolpath %s -l", toolDir), "")

	// Assert
	mustContain := []string{
		"  Rubocop     | puppetlabs | rubocop     | https://github.com/rubocop/rubocop         | 0.1.0    \n",
		"  puppet-lint | puppetlabs | puppet-lint | https://github.com/puppetlabs/puppet-lint/ | 0.1.0    \n",
	}

	for _, s := range mustContain {
		assert.Contains(t, stdout, s)
	}
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmValidate_List_Flag_Invalid_ToolDir(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	skipValidationTests(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/invalid/tooldir")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("validate -l --toolpath %s", toolDir), "")

	assert.Contains(t, stdout, fmt.Sprintf("no validators found in %s", toolDir))
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmValidate_Tool_Timeout_Flag_Exceeded(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	skipValidationTests(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/validate/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("validate puppetlabs/timeout --toolpath %s --cachedir %s --codedir %s --toolTimeout 5", toolDir, cacheDir, codeDir), "")

	assert.Contains(t, stdout, "context deadline exceeded")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmValidate_Invalid_Tool_Timeout_Flag(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	skipValidationTests(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/validate/testdata/tooldir")
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("validate puppetlabs/timeout --toolpath %s --codedir %s --toolTimeout -5", toolDir, codeDir), "")

	assert.Contains(t, stdout, "the --toolTimeout flag must be set to a value greater than 1")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}
