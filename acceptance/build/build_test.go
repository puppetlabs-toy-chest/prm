package build_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pdkgo/acceptance/testutils"
	"github.com/stretchr/testify/assert"
)

const APP = "prm"

func Test_PrmBuild_Outputs_TarGz(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	toolName := "good-project"

	sourceDir, _ := filepath.Abs("../../acceptance/build/testdata")
	toolDir := filepath.Join(sourceDir, toolName)
	wd := testutils.GetTmpDir(t)

	cmd := fmt.Sprintf("build --sourcedir %v --targetdir %v", toolDir, wd)
	stdout, stderr, exitCode := testutils.RunAppCommand(cmd, "")

	expectedOutputFilePath := filepath.Join(wd, fmt.Sprintf("%v.tar.gz", toolName))

	assert.Contains(t, stdout, fmt.Sprintf("Packaged tool output to %v", expectedOutputFilePath))
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
	assert.FileExists(t, expectedOutputFilePath)
}

func Test_PrmBuild_With_NoTargetDir_Outputs_TarGz(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	toolName := "good-project"

	sourceDir, _ := filepath.Abs("../../acceptance/build/testdata")
	toolDir := filepath.Join(sourceDir, toolName)
	wd := testutils.GetTmpDir(t)

	cmd := fmt.Sprintf("build --sourcedir %v", toolDir)
	stdout, stderr, exitCode := testutils.RunAppCommand(cmd, wd)

	expectedOutputFilePath := filepath.Join(wd, "pkg", fmt.Sprintf("%v.tar.gz", toolName))

	assert.Contains(t, stdout, fmt.Sprintf("Packaged tool output to %v", expectedOutputFilePath))
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
	assert.FileExists(t, expectedOutputFilePath)
}

func Test_PrmBuild_With_EmptySourceDir_Errors(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	toolName := "no-project-here"

	sourceDir, _ := filepath.Abs("../../acceptance/build/testdata")
	toolDir := filepath.Join(sourceDir, toolName)

	cmd := fmt.Sprintf("build --sourcedir %v", toolDir)
	stdout, stderr, exitCode := testutils.RunAppCommand(cmd, "")

	assert.Contains(t, stdout, fmt.Sprintf("No project directory at %v", toolDir))
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmBuild_With_NoPrmConfig_Errors(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	toolName := "no-prm-config-project"

	sourceDir, _ := filepath.Abs("../../acceptance/build/testdata")
	toolDir := filepath.Join(sourceDir, toolName)

	cmd := fmt.Sprintf("build --sourcedir %v", toolDir)
	stdout, stderr, exitCode := testutils.RunAppCommand(cmd, "")

	assert.Contains(t, stdout, fmt.Sprintf("No 'prm-config.yml' found in %v", toolDir))
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmBuild_With_NoContentDir_Errors(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	testutils.SetAppName(APP)

	toolName := "no-content-dir-project"

	sourceDir, _ := filepath.Abs("../../acceptance/build/testdata")
	toolDir := filepath.Join(sourceDir, toolName)

	cmd := fmt.Sprintf("build --sourcedir %v", toolDir)
	stdout, stderr, exitCode := testutils.RunAppCommand(cmd, "")

	assert.Contains(t, stdout, fmt.Sprintf("No 'content' dir found in %v", toolDir))
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}
