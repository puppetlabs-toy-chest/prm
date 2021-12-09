package exec_test

import (
	"fmt"
	"path/filepath"
	"testing"

	dircopy "github.com/otiai10/copy"
	"github.com/puppetlabs/pdkgo/acceptance/testutils"
	"github.com/stretchr/testify/assert"
)

const APP = "prm"

func Test_PrmExec_Tool_SingleGem(t *testing.T) {
	// testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/exec/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")


	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("exec puppetlabs/puppet-lint --toolpath %s --cachedir %s --codedir %s --additionalArgs=invalid.pp", toolDir, cacheDir, codeDir), "")

	// Assert
	assert.Contains(t, stdout, "puppetlabs/puppet-lint executed successfully")
	assert.Contains(t, stdout, "ERROR: invalid not in autoload module layout on line 1 (check: autoloader_layout)")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

// Util functions

func createCodeDir(t *testing.T, codeDirSrc string) (codeDir string) {
	codeDirSrc, _ = filepath.Abs(filepath.Join("../../acceptance/exec/testdata/codedirs", codeDirSrc))
	codeDir = testutils.GetTmpDir(t)
	dircopy.Copy(codeDirSrc, codeDir)
	return codeDir
}
