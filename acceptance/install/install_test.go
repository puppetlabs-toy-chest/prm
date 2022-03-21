package install_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pct/acceptance/testutils"
	"github.com/stretchr/testify/assert"
)

var defaultToolPath string

const APP = "prm"

// This test may not work locally if the default tool path is set to a different location in the `~/.`config/.prm.yaml` file.
func Test_PrmInstall_InstallsTo_DefaultToolPath(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)

	// Setup
	toolPkgPath, _ := filepath.Abs(fmt.Sprintf("../../acceptance/install/testdata/%v.tar.gz", "good-project"))
	installCmd := fmt.Sprintf("install %v", toolPkgPath)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

	// Assert
	assert.Contains(t, stdout, fmt.Sprintf("Tool installed to %v", filepath.Join(getDefaultToolPath(), "gooder", "good-project", "0.1.0")))
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
	assert.FileExists(t, filepath.Join(getDefaultToolPath(), "gooder", "good-project", "0.1.0", "prm-config.yml"))
	assert.FileExists(t, filepath.Join(getDefaultToolPath(), "gooder", "good-project", "0.1.0", "content", "empty.txt"))
	assert.FileExists(t, filepath.Join(getDefaultToolPath(), "gooder", "good-project", "0.1.0", "content", "goodfile.txt.tmpl"))

	stdout, stderr, exitCode = testutils.RunAppCommand("exec --list", "")
	assert.Regexp(t, "Good\\sProject\\s+\\|\\sgooder\\s+\\|\\sgood-project\\s+\\|\\shttps://github.com/puppetlabs/pct-good-project\\s+\\|\\s0.1.0", stdout)
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)

	// Tear Down
	removeInstalledTool(filepath.Join(getDefaultToolPath(), "gooder", "good-project", "0.1.0"))
}

type toolData struct {
	name          string
	author        string
	listExpRegex  string
	expectedFiles []string
	gitUri        string
}

func Test_PrmInstall_InstallsTo_DefinedToolPath(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)

	// Setup
	toolPath := testutils.GetTmpDir(t)

	toolPkgs := []toolData{
		{
			name:         "additional-project",
			author:       "adder",
			listExpRegex: "Additional\\sProject\\s+\\|\\sadder\\s+\\|\\sadditional-project\\s+\\|\\shttps://github.com/puppetlabs/pct-additional-project\\s+\\|\\s0.1.0",
			expectedFiles: []string{
				"prm-config.yml",
				"content/empty.txt",
				"content/goodfile.txt.tmpl",
			},
		},
		{
			name:         "good-project",
			author:       "gooder",
			listExpRegex: "Good\\sProject\\s+\\|\\sgooder\\s+\\|\\sgood-project\\s+\\|\\shttps://github.com/puppetlabs/pct-good-project\\s+\\|\\s0.1.0",
			expectedFiles: []string{
				"prm-config.yml",
				"content/empty.txt",
				"content/goodfile.txt.tmpl",
			},
		},
	}

	for _, tool := range toolPkgs {
		// Setup
		toolPkgPath, _ := filepath.Abs(fmt.Sprintf("../../acceptance/install/testdata/%v.tar.gz", tool.name))
		installCmd := fmt.Sprintf("install %v --toolpath %v", toolPkgPath, toolPath)

		// Exec
		stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

		// Assert
		assert.Contains(t, stdout, fmt.Sprintf("Tool installed to %v", filepath.Join(toolPath, tool.author, tool.name, "0.1.0")))
		assert.Equal(t, "", stderr)
		assert.Equal(t, 0, exitCode)
	}

	for _, tool := range toolPkgs {
		// Assert
		for _, file := range tool.expectedFiles {
			assert.FileExists(t, filepath.Join(toolPath, tool.author, tool.name, "0.1.0", file))
		}

		listCmd := fmt.Sprintf("exec --list --toolpath %v", toolPath)
		stdout, stderr, exitCode := testutils.RunAppCommand(listCmd, "")

		assert.Regexp(t, tool.listExpRegex, stdout)
		assert.Equal(t, "", stderr)
		assert.Equal(t, 0, exitCode)
	}

	// Tear Down
	for _, tool := range toolPkgs {
		removeInstalledTool(filepath.Join(toolPath, tool.author, tool.name, "0.1.0"))
	}
}

func Test_PrmInstall_InstallsFrom_RemoteToolPath(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)

	// Setup
	toolPath := testutils.GetTmpDir(t)

	toolPkgs := []toolData{
		{
			name:         "additional-project",
			author:       "adder",
			listExpRegex: "Additional\\sProject\\s+\\|\\sadder\\s+\\|\\sadditional-project\\s+\\|\\shttps://github.com/puppetlabs/pct-additional-project\\s+\\|\\s0.1.0",
			expectedFiles: []string{
				"prm-config.yml",
				"content/empty.txt",
				"content/goodfile.txt.tmpl",
			},
		},
		{
			name:         "good-project",
			author:       "gooder",
			listExpRegex: "Good\\sProject\\s+\\|\\sgooder\\s+\\|\\sgood-project\\s+\\|\\shttps://github.com/puppetlabs/pct-good-project\\s+\\|\\s0.1.0",
			expectedFiles: []string{
				"prm-config.yml",
				"content/empty.txt",
				"content/goodfile.txt.tmpl",
			},
		},
	}

	for _, tool := range toolPkgs {
		// Setup
		toolPkgPath := fmt.Sprintf("https://github.com/puppetlabs/prm/raw/main/acceptance/install/testdata/%s.tar.gz", tool.name)
		installCmd := fmt.Sprintf("install %v --toolpath %v", toolPkgPath, toolPath)

		// Exec
		stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

		// Assert
		assert.Contains(t, stdout, fmt.Sprintf("Tool installed to %v", filepath.Join(toolPath, tool.author, tool.name, "0.1.0")))
		assert.Equal(t, "", stderr)
		assert.Equal(t, 0, exitCode)
	}

	for _, tool := range toolPkgs {
		// Assert
		for _, file := range tool.expectedFiles {
			assert.FileExists(t, filepath.Join(toolPath, tool.author, tool.name, "0.1.0", file))
		}

		listCmd := fmt.Sprintf("exec --list --toolpath %v", toolPath)
		stdout, stderr, exitCode := testutils.RunAppCommand(listCmd, "")

		assert.Regexp(t, tool.listExpRegex, stdout)
		assert.Equal(t, "", stderr)
		assert.Equal(t, 0, exitCode)
	}

	// Tear Down
	for _, tool := range toolPkgs {
		removeInstalledTool(filepath.Join(toolPath, tool.author, tool.name, "0.1.0"))
	}
}

func Test_PrmInstall_Errors_When_NoToolPkgDefined(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("install", "")

	// Assert
	assert.Contains(t, stdout, "Path to tool package (tar.gz) should be first argument")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmInstall_Errors_When_ToolPkgNotExist(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)

	// Setup
	toolPkgPath, _ := filepath.Abs("/path/to/nowhere/good-project.tar.gz")
	installCmd := fmt.Sprintf("install %v", toolPkgPath)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

	// Assert
	assert.Contains(t, stdout, fmt.Sprintf("No package at %v", toolPkgPath))
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmInstall_Errors_When_InvalidGzProvided(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)

	// Setup
	toolPkgPath, _ := filepath.Abs("../../acceptance/install/testdata/invalid-gz-project.tar.gz")
	installCmd := fmt.Sprintf("install %v", toolPkgPath)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

	// Assert
	assert.Contains(t, stdout, fmt.Sprintf("Could not extract TAR from GZIP (%v)", toolPkgPath))
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmInstall_Errors_When_InvalidTarProvided(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)

	// Setup
	toolPkgPath, _ := filepath.Abs("../../acceptance/install/testdata/invalid-tar-project.tar.gz")
	installCmd := fmt.Sprintf("install %v", toolPkgPath)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

	assert.Contains(t, stdout, fmt.Sprintf("Could not UNTAR package (%v)", toolPkgPath))
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmInstall_FailsWhenToolAlreadyExists(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)

	// Setup
	toolPath := testutils.GetTmpDir(t)

	// Install tool
	toolPkgPath, _ := filepath.Abs("../../acceptance/install/testdata/additional-project.tar.gz")
	installCmd := fmt.Sprintf("install %v --toolpath %v", toolPkgPath, toolPath)
	stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

	// verify the tool installed
	assert.Contains(t, stdout, fmt.Sprintf("Tool installed to %v", filepath.Join(toolPath, "adder", "additional-project", "0.1.0")))
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)

	// Attempt to reinstall the tool
	toolPkgPath, _ = filepath.Abs("../../acceptance/install/testdata/additional-project.tar.gz")
	installCmd = fmt.Sprintf("install %v --toolpath %v", toolPkgPath, toolPath)
	stdout, stderr, exitCode = testutils.RunAppCommand(installCmd, "")

	// verify that the tool failed to install
	assert.Contains(t, stdout, "Unable to install in namespace: Package already installed")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)

	// Tear Down
	removeInstalledTool(filepath.Join(toolPath, "adder", "additional-project", "0.1.0"))
}

func Test_PrmInstall_ForceSuccessWhenToolAlreadyExists(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)

	// Setup
	toolPath := testutils.GetTmpDir(t)

	// Install tool
	toolPkgPath, _ := filepath.Abs("../../acceptance/install/testdata/additional-project.tar.gz")
	installCmd := fmt.Sprintf("install %v --toolpath %v", toolPkgPath, toolPath)
	stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

	// verify the tool installed
	assert.Contains(t, stdout, fmt.Sprintf("Tool installed to %v", filepath.Join(toolPath, "adder", "additional-project", "0.1.0")))
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)

	// Attempt to reinstall the tool
	toolPkgPath, _ = filepath.Abs("../../acceptance/install/testdata/additional-project.tar.gz")
	installCmd = fmt.Sprintf("install %v --force --toolpath %v", toolPkgPath, toolPath)
	stdout, stderr, exitCode = testutils.RunAppCommand(installCmd, "")

	// verify that the tool reinstall exited successfully
	assert.Contains(t, stdout, fmt.Sprintf("Tool installed to %v", filepath.Join(toolPath, "adder", "additional-project", "0.1.0")))
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)

	// Tear Down
	removeInstalledTool(filepath.Join(toolPath, "adder", "additional-project", "0.1.0"))
}

func Test_PrmInstall_WithGitUri_InstallTool(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)

	// Setup
	toolPath := testutils.GetTmpDir(t)

	toolPkgs := []toolData{
		{
			name:         "test-tool-1",
			author:       "test-user",
			listExpRegex: "Test\\sTool\\s1\\s+\\|\\stest-user\\s+\\|\\stest-tool-1",
			expectedFiles: []string{
				"prm-config.yml",
			},
			gitUri: "https://github.com/puppetlabs/prm-test-tool-01.git",
		},
		{
			name:         "test-tool-2",
			author:       "test-user",
			listExpRegex: "Test\\sTool\\s2\\s+\\|\\stest-user\\s+\\|\\stest-tool-2",
			expectedFiles: []string{
				"prm-config.yml",
			},
			gitUri: "https://github.com/puppetlabs/prm-test-tool-02.git",
		},
	}

	for _, tool := range toolPkgs {
		installCmd := fmt.Sprintf("install --git-uri %v --toolpath %v", tool.gitUri, toolPath)

		// Exec
		stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

		// Assert
		assert.Contains(t, stdout, fmt.Sprintf("Tool installed to %v", filepath.Join(toolPath, tool.author, tool.name, "0.1.0")))
		assert.Equal(t, "", stderr)
		assert.Equal(t, 0, exitCode)
	}

	for _, tool := range toolPkgs {
		// Assert
		for _, file := range tool.expectedFiles {
			assert.FileExists(t, filepath.Join(toolPath, tool.author, tool.name, "0.1.0", file))
		}

		listCmd := fmt.Sprintf("exec --list --toolpath %v", toolPath)
		stdout, stderr, exitCode := testutils.RunAppCommand(listCmd, "")

		assert.Regexp(t, tool.listExpRegex, stdout)
		assert.Equal(t, "", stderr)
		assert.Equal(t, 0, exitCode)
	}

	// Tear Down
	for _, tool := range toolPkgs {
		removeInstalledTool(filepath.Join(toolPath, tool.author, tool.name, "0.1.0"))
	}
}

func Test_PrmInstall_WithGitUri_FailsWithNonExistentUri(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("install --git-uri https://example.com/fake-git-uri", "")

	// Assert
	assert.Contains(t, stdout, "Could not clone git repository:")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmInstall_WithGitUri_FailsWithInvalidUri(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("install --git-uri example.com/invalid-git-uri", "")

	// Assert
	assert.Contains(t, stdout, "Could not parse package uri")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmInstall_WithGitUri_RemovesHiddenGitDir(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	testutils.SetAppName(APP)

	// Setup
	toolPath := testutils.GetTmpDir(t)

	// Install tool
	installCmd := fmt.Sprint("install --git-uri https://github.com/puppetlabs/prm-test-tool-01.git --toolpath ", toolPath)
	stdout, stderr, exitCode := testutils.RunAppCommand(installCmd, "")

	// Verify the tool installed
	assert.Contains(t, stdout, fmt.Sprintf("Tool installed to %v", filepath.Join(toolPath, "test-user", "test-tool-1", "0.1.0")))
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)

	// Check .git directory has been deleted
	assert.NoDirExists(t, filepath.Join(toolPath, "test-user", "test-tool-1", "0.1.0", ".git"))

	// Tear Down
	removeInstalledTool(filepath.Join(toolPath, "test-user", "test-tool-1", "0.1.0"))
}

// Util Functions
func removeInstalledTool(toolPath string) {
	_, err := os.Stat(toolPath)
	if err != nil {
		note := "NOTE: This test may not work locally if the default tool path is set to a different location in the `~/.`config/.prm.yaml` file."
		panic(fmt.Sprintf("removeInstalledTool(): Could not determine if tool path (%v) exists: %v\n%v", toolPath, err, note))
	}

	os.RemoveAll(toolPath)
	if err != nil {
		panic(fmt.Sprintf("remoteTool(): Could not remove %v: %v", toolPath, err))
	}
}

func getDefaultToolPath() string {
	if defaultToolPath != "" {
		return defaultToolPath
	}

	entries, err := filepath.Glob("../../dist/prm_*/tools")
	if err != nil {
		panic("getDefaultToolPath(): Could not determine default tool path")
	}
	if len(entries) != 1 {
		panic(fmt.Sprintf("getDefaultToolPath(): Could not determine default tool path; matched entries: %v", len(entries)))
	}

	defaultToolPath, err := filepath.Abs(entries[0])
	if err != nil {
		panic(fmt.Sprintf("getDefaultToolPath(): Could not create absolute path to toolpath: %v", err))
	}

	return defaultToolPath
}
