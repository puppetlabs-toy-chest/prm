package exec_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pct/acceptance/testutils"
	"github.com/stretchr/testify/assert"
)

// TODO: See GH-66 - exit message validation will need updated when that ticket is addressed

func Test_PrmExec_Tool_SingleGem_ToolExitError(t *testing.T) {
	skipExecTests(t)
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/exec/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("exec puppetlabs/puppet-lint --toolpath %s --cachedir %s --codedir %s --toolArgs=invalid.pp", toolDir, cacheDir, codeDir), "")

	// Assert
	assert.Contains(t, stdout, "Error executing tool puppetlabs/puppet-lint: Tool exited with code: 1") // GH-66
	assert.Contains(t, stdout, "ERROR: invalid not in autoload module layout on line 1 (check: autoloader_layout)")
	assert.Contains(t, stdout, "WARNING: class not documented on line 1 (check: documentation)")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmExec_Tool_SingleGem_ToolExitZero(t *testing.T) {
	skipExecTests(t)
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/exec/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("exec puppetlabs/puppet-lint --toolpath %s --cachedir %s --codedir %s --toolArgs=valid", toolDir, cacheDir, codeDir), "")

	// Assert
	assert.Contains(t, stdout, "Tool puppetlabs/puppet-lint executed successfully") // GH-66
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}
