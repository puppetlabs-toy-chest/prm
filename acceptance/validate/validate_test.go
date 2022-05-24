package validate_test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"

	dircopy "github.com/otiai10/copy"
	"github.com/puppetlabs/pct/acceptance/testutils"
	"github.com/stretchr/testify/assert"
)

const APP = "prm"

func Test_PrmValidate_Single_Tool_Pass_Output_To_Terminal(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	skipValidationTests(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/validate/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Exec
	_, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("validate puppetlabs/puppet-lint --toolpath %s --cachedir %s --codedir %s --toolArgs=valid/manifests/valid.pp", toolDir, cacheDir, codeDir), "")

	// Assert
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmValidate_Single_Tool_Fail_Output_To_Terminal(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	skipValidationTests(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/validate/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("validate puppetlabs/puppet-lint --toolpath %s --cachedir %s --codedir %s --toolArgs=invalid.pp", toolDir, cacheDir, codeDir), "")

	// Assert
	assert.Contains(t, stdout, "Validation returned 1 error")
	assert.Contains(t, stdout, "ERROR: invalid not in autoload module layout on line 1 (check: autoloader_layout)")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmValidate_Single_Tool_Pass_Output_To_File(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	skipValidationTests(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/validate/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("validate puppetlabs/puppet-lint --toolpath %s --cachedir %s --codedir %s --resultsView file --toolArgs=valid/manifests/valid.pp", toolDir, cacheDir, codeDir), "")

	// Assert
	mustContain := []string{""}
	checkLogFileAndContents(t, path.Join(codeDir, ".prm-validate"), mustContain, stdout, "puppet-lint")
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmValidate_Single_Tool_Fail_Output_To_File(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	skipValidationTests(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/validate/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("validate puppetlabs/puppet-lint --toolpath %s --cachedir %s --codedir %s --resultsView file --toolArgs=invalid.pp", toolDir, cacheDir, codeDir), "")

	// Assert
	mustContain := []string{
		"ERROR: invalid not in autoload module layout on line 1 (check: autoloader_layout)",
		"WARNING: class not documented on line 1 (check: documentation)",
	}
	checkLogFileAndContents(t, path.Join(codeDir, ".prm-validate"), mustContain, stdout, "puppet-lint")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmValidate_Validate_File_Multitool_Fail_Output_To_Terminal(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	skipValidationTests(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/validate/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")
	validateFileContents := `groups:
  - id: test_group_1
    tools:
      - name: puppetlabs/puppet-lint
        args: [invalid.pp]
      - name: puppetlabs/rubocop
`
	createValidateFile(codeDir, validateFileContents)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("validate --toolpath %s --cachedir %s --codedir %s --group test_group_1 --resultsView terminal", toolDir, cacheDir, codeDir), "")

	// Assert
	assert.Contains(t, stdout, "WARNING: class not documented on line 1 (check: documentation)")
	assert.Contains(t, stdout, "ERROR: invalid not in autoload module layout on line 1 (check: autoloader_layout)")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmValidate_Validate_File_Multitool_Fail_Output_To_File(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	skipValidationTests(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/validate/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")
	validateFileContents := `groups:
  - id: test_group_1
    tools:
      - name: puppetlabs/puppet-lint
        args: [invalid.pp]
      - name: puppetlabs/rubocop
`
	createValidateFile(codeDir, validateFileContents)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("validate --toolpath %s --cachedir %s --codedir %s --group test_group_1 --resultsView file", toolDir, cacheDir, codeDir), "")

	// Assert
	mustContain := map[string][]string{
		"puppet-lint": {
			"WARNING: class not documented on line 1 (check: documentation)",
		},
		"rubocop": {
			"Inspecting 0 files", // TODO: Setup a proper tool for this test
			"0 files inspected, no offenses detected",
		},
	}

	// Assert
	for key, value := range mustContain {
		checkLogFileAndContents(t, path.Join(codeDir, ".prm-validate", "test_group_1"), value, stdout, key)
	}
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmValidate_Validate_File_Multitool_Pass_Output_To_Terminal(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	skipValidationTests(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/validate/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")
	validateFileContents := `groups:
  - id: test_group_1
    tools:
      - name: puppetlabs/puppet-lint
        args: [valid/manifests/valid.pp]
      - name: puppetlabs/rubocop
`
	createValidateFile(codeDir, validateFileContents)

	// Exec
	_, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("validate --toolpath %s --cachedir %s --codedir %s --group test_group_1 --resultsView terminal", toolDir, cacheDir, codeDir), "")

	// Assert
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmValidate_Validate_File_Multitool_Pass_Output_To_File(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	skipValidationTests(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/validate/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")
	validateFileContents := `groups:
  - id: test_group_1
    tools:
      - name: puppetlabs/puppet-lint
        args: [valid/manifests/valid.pp]
      - name: puppetlabs/rubocop
`
	createValidateFile(codeDir, validateFileContents)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("validate --toolpath %s --cachedir %s --codedir %s --group test_group_1 --resultsView file", toolDir, cacheDir, codeDir), "")

	mustContain := map[string][]string{
		"puppet-lint": {
			"",
		},
		"rubocop": {
			"Inspecting 0 files", // TODO: Setup a proper tool for this test
			"0 files inspected, no offenses detected",
		},
	}

	// Assert
	for key, value := range mustContain {
		checkLogFileAndContents(t, path.Join(codeDir, ".prm-validate", "test_group_1"), value, stdout, key)
	}
	assert.Equal(t, "", stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmValidate_Validate_File_Multitool_Duplicate_Tool_Error(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	skipValidationTests(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/validate/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")
	validateFileContents := `groups:
  - id: test_group_1
    tools:
      - name: puppetlabs/puppet-lint
        args: [invalid.pp]
      - name: puppetlabs/puppet-lint
`
	createValidateFile(codeDir, validateFileContents)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("validate --toolpath %s --cachedir %s --codedir %s --group test_group_1 --resultsView terminal", toolDir, cacheDir, codeDir), "")

	// Assert
	assert.Contains(t, stdout, "duplicate tool 'puppetlabs/puppet-lint' found. Validation groups cannot contain duplicate tools")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmValidate_Validate_File_Multitool_No_Tools_Error(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	skipValidationTests(t)

	// Setup
	testutils.SetAppName(APP)
	toolDir, _ := filepath.Abs("../../acceptance/validate/testdata/tooldir")
	cacheDir := testutils.GetTmpDir(t)
	codeDir := createCodeDir(t, "puppet-lint-playground")
	testutils.RunAppCommand("set puppet 7.12.1", "")
	validateFileContents := `groups:
  - id: test_group_1
    tools:
`
	createValidateFile(codeDir, validateFileContents)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("validate --toolpath %s --cachedir %s --codedir %s --group test_group_1 --resultsView terminal", toolDir, cacheDir, codeDir), "")

	// Assert
	assert.Contains(t, stdout, "no tools provided for validation")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func createCodeDir(t *testing.T, codeDirSrc string) (codeDir string) {
	if codeDirSrc != "" {
		codeDirSrc, _ = filepath.Abs(filepath.Join("../../acceptance/exec/testdata/codedirs", codeDirSrc))
	}
	codeDir = testutils.GetTmpDir(t)
	dircopy.Copy(codeDirSrc, codeDir) //nolint:gosec,errcheck // we should know what we've given this func
	return codeDir
}

func createValidateFile(codeDir string, validateContent string) {
	os.WriteFile(filepath.Join(codeDir, "validate.yml"), []byte(validateContent), 0644) //nolint:gosec,errcheck
}

func checkLogFileAndContents(t *testing.T, logDirPath string, expectedLogContents []string, stdout string, toolName string) {
	logFilePath := filepath.Join(logDirPath, getFileLocationFromStdout(stdout, toolName))
	assert.FileExists(t, logFilePath)
	logContents, err := os.ReadFile(logFilePath)
	contentsStr := string(logContents)
	assert.NoError(t, err)
	for _, expected := range expectedLogContents {
		assert.Contains(t, contentsStr, expected)
	}
}

func getFileLocationFromStdout(stdout, toolName string) string {
	re := regexp.MustCompile(fmt.Sprint(toolName, `_(.*?).log`))
	matches := re.FindStringSubmatch(stdout)
	if len(matches) > 0 {
		return matches[0]
	}
	return ""
}

func skipValidationTests(t *testing.T) {
	if _, isCI := os.LookupEnv("CI"); (runtime.GOOS == "darwin" || runtime.GOOS == "windows") && isCI {
		t.Skip("Skipping exec acceptance tests for darwin and windows in CI")
	}
}
