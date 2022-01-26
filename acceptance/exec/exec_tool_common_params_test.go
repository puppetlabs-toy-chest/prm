package exec_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	dircopy "github.com/otiai10/copy"
	"github.com/puppetlabs/pdkgo/acceptance/testutils"
	"github.com/stretchr/testify/assert"
)

const APP = "prm"

// Ensure:
// - Git is not installed
// - Default help arg is '--help'
// - Success exit code is 0
func Test_PrmExec_UndefinedCommonParam_Defaults(t *testing.T) {
	skipExecTests(t)
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/exec/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Non zero exit code should be interpreted as a failure
	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("exec puppetlabs/common-undefined --toolpath %s --cachedir %s --codedir %s --toolArgs=foo", toolDir, cacheDir, codeDir), "")

	// Assert
	assert.Contains(t, stdout, "Error executing tool puppetlabs/common-undefined: Tool exited with code: 1") // GH-66
	assert.Contains(t, stdout, "Unknown Puppet subcommand 'foo'")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)

	// Zero exit code should be interpreted as success
	// Exec
	stdout, stderr, exitCode = testutils.RunAppCommand(fmt.Sprintf("exec puppetlabs/common-undefined --toolpath %s --cachedir %s --codedir %s --toolArgs=--version", toolDir, cacheDir, codeDir), "")

	// Assert
	assert.Contains(t, stdout, "Tool puppetlabs/common-undefined executed successfully") // GH-66
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)

	// GH-74: Uncomment when this ticket is addressed
	// --help should invoke --help output from tool
	// // Exec
	// stdout, stderr, exitCode = testutils.RunAppCommand(fmt.Sprintf("exec puppetlabs/common-undefined --toolpath %s --cachedir %s --help", toolDir, cacheDir), "")

	// // Assert
	// assert.Contains(t, stdout, "puppet-lint\n\nBasic Command Line Usage:") // GH-66
	// assert.Empty(t, stderr)
	// assert.Equal(t, 0, exitCode)

}

// GH-79: Uncomment when this ticket is addressed
// func Test_PrmExec_Tool_NeedWriteAccess_Undefined(t *testing.T) {
// 	testutils.SkipAcceptanceTest(t)

// 	// Setup
// 	testutils.SetAppName(APP)
// 	toolDir, _ := filepath.Abs("../../acceptance/exec/testdata/tooldir")
// 	cacheDir := testutils.GetTmpDir(t)
// 	codeDir := createCodeDir(t, "")
// 	testutils.RunAppCommand("set puppet 7.12.1", "")

// 	// Non zero exit code should be interpreted as a failure
// 	// Exec
// 	stdout, stderr, exitCode := RunAppCommandMooCow(fmt.Sprintf("exec puppetlabs/common-undefined --toolpath %s --cachedir %s --codedir %s --toolArgs=apply -e \"file { '/code/foo': ensure => 'present', }\"", toolDir, cacheDir, codeDir), "")

// 	// Assert
// 	assert.Contains(t, stdout, "Error executing tool puppetlabs/common-undefined: Tool exited with code: 1") // GH-66
// 	assert.Contains(t, stdout, "ERROR: Permission denied @ rb_sysopen - /code/foo")
// 	assert.Equal(t, "exit status 1", stderr)
// 	assert.Equal(t, 1, exitCode)
// }

func Test_PrmExec_Tool_NeedWriteAccess(t *testing.T) {
	skipExecTests(t)
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/exec/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("exec puppetlabs/needs-write-access --toolpath %s --cachedir %s --codedir %s --toolArgs=apply -e \"file { '/code/foo': ensure => 'present', }\"", toolDir, cacheDir, codeDir), "")

	// Assert
	assert.Contains(t, stdout, "Tool puppetlabs/needs-write-access executed successfully") // GH-66
	assert.Contains(t, stdout, "File[/code/foo]/ensure: created")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmExec_Tool_UseScript(t *testing.T) {
	skipExecTests(t)
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/exec/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("exec puppetlabs/use-script --alwaysBuild --toolpath %s --cachedir %s --codedir %s", toolDir, cacheDir, codeDir), "")

	// Assert
	assert.Contains(t, stdout, "Tool puppetlabs/use-script executed successfully") // GH-66
	assert.Contains(t, stdout, "This script is so foo bar!")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmExec_Tool_UseScriptBadEntrypoint(t *testing.T) {
	skipExecTests(t)
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/exec/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("exec puppetlabs/use-script-bad --alwaysBuild --toolpath %s --cachedir %s --codedir %s", toolDir, cacheDir, codeDir), "")

	// Assert
	assert.Contains(t, stdout, "exec: \"/tmp/foo-bad.sh\": stat /tmp/foo-bad.sh: no such file or directory: unknown") // GH-66
	assert.Contains(t, stderr, "exit status 1")
	assert.Equal(t, 1, exitCode)
}

// Requires Git works correctly when tested manually, does not work in CI due to a permission denied error
// when removing the .git folder from the cloned repo on test cleanup
// func Test_PrmExec_Tool_RequiresGit(t *testing.T) {
// 	testutils.SkipAcceptanceTest(t)

// 	testutils.SetAppName(APP)
// 	toolDir, _ := filepath.Abs("../../acceptance/exec/testdata/tooldir")
// 	cacheDir := testutils.GetTmpDir(t)
// 	codeDir := createCodeDir(t, "")
// 	testutils.RunAppCommand("set puppet 7.12.1", "")

// 	// Exec
// 	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("exec puppetlabs/requires-git --alwaysBuild --toolpath %s --cachedir %s --codedir %s", toolDir, cacheDir, codeDir), "")

// 	fmt.Printf("stdout: %s", stdout)
// 	// Assert
// 	assert.Contains(t, stdout, "Tool puppetlabs/requires-git executed successfully") // GH-66
// 	// assert.Contains(t, stdout, "This script is so foo bar!")
// 	assert.DirExists(t, filepath.Join(codeDir, "prm-test-template-01", "content"))

// 	assert.Empty(t, stderr)
// 	assert.Equal(t, 1, exitCode)
// }

func Test_PrmExec_Tool_ValidationYaml(t *testing.T) {
	skipExecTests(t)
	testutils.SkipAcceptanceTest(t)
}

func Test_PrmExec_Tool_PuppetEnabled(t *testing.T) {
	skipExecTests(t)
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/exec/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("exec puppetlabs/check-puppet --alwaysBuild --toolpath %s --cachedir %s --codedir %s", toolDir, cacheDir, codeDir), "")

	fmt.Printf("stdout: %s", stdout)
	// Assert
	assert.Contains(t, stdout, "Tool puppetlabs/check-puppet executed successfully") // GH-66

	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmExec_alwaysBuild_Flag(t *testing.T) {
	skipExecTests(t)
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/exec/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "")
	testutils.RunAppCommand("set puppet 7.12.1", "") // Reset script file file

	// Reset script
	scriptPath := filepath.Join(toolDir, "puppetlabs/build-flag/0.1.0/content/script.sh")
	file, err := os.Create(scriptPath)
	assert.NoError(t, err)
	_, err = file.WriteString("#!/bin/bash\necho 'This script is so foo bar!'\nexit 0")
	assert.NoError(t, err)
	err = file.Close()
	assert.NoError(t, err)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("exec puppetlabs/build-flag --alwaysBuild --toolpath %s --cachedir %s --codedir %s", toolDir, cacheDir, codeDir), "")

	fmt.Printf("stdout: %s", stdout)
	// Assert
	assert.Contains(t, stdout, "Tool puppetlabs/build-flag executed successfully") // GH-66
	assert.Contains(t, stdout, "This script is so foo bar!")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)

	// Change the script without the alwayBuild flag, thereby not changing the outputted message
	file, err = os.Create(scriptPath)
	assert.NoError(t, err)
	_, err = file.WriteString("#!/bin/bash\necho 'Different output!'\nexit 0")
	assert.NoError(t, err)
	err = file.Close()
	assert.NoError(t, err)

	stdout, stderr, exitCode = testutils.RunAppCommand(fmt.Sprintf("exec puppetlabs/build-flag --toolpath %s --cachedir %s --codedir %s", toolDir, cacheDir, codeDir), "")

	assert.Contains(t, stdout, "Tool puppetlabs/build-flag executed successfully") // GH-66
	assert.Contains(t, stdout, "This script is so foo bar!")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)

	// Enable alwaysBuild flag to rebuild the image and display new message
	stdout, stderr, exitCode = testutils.RunAppCommand(fmt.Sprintf("exec puppetlabs/build-flag --alwaysBuild --toolpath %s --cachedir %s --codedir %s", toolDir, cacheDir, codeDir), "")

	assert.Contains(t, stdout, "Tool puppetlabs/build-flag executed successfully") // GH-66
	assert.Contains(t, stdout, "Different output!")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func createCodeDir(t *testing.T, codeDirSrc string) (codeDir string) {
	if codeDirSrc != "" {
		codeDirSrc, _ = filepath.Abs(filepath.Join("../../acceptance/exec/testdata/codedirs", codeDirSrc))
	}
	codeDir = testutils.GetTmpDir(t)
	dircopy.Copy(codeDirSrc, codeDir) //nolint:gosec,errcheck // we should know what we've given this func
	return codeDir
}

func skipExecTests(t *testing.T) {
	if _, isCI := os.LookupEnv("CI"); (runtime.GOOS == "darwin" || runtime.GOOS == "windows") && isCI {
		t.Skip("Skipping exec acceptance tests for darwin and windows in CI")
	}
}
