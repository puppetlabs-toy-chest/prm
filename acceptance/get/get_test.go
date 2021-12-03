package get_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pdkgo/acceptance/testutils"
	"github.com/stretchr/testify/assert"
)

const APP = "prm"

func Test_PrmGet_NoArgs(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("get", "")

	// Assert
	assert.Contains(t, stdout, "Displays the requested configuration value\n\nUsage:\n  prm get")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmGet_Puppet_FirstRun(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)
	tempDir := testutils.GetTmpDir(t)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("get puppet --config %s", filepath.Join(tempDir, ".prm.yaml")), "")

	// Assert
	assert.Contains(t, stdout, "Puppet version is configured to: 7.0.0")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmGet_Backend_FirstRun(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)
	tempDir := testutils.GetTmpDir(t)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("get backend --config %s", filepath.Join(tempDir, ".prm.yaml")), "")

	// Assert
	assert.Contains(t, stdout, "Backend is configured to: docker")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmGet_Puppet_EnsureSet(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)
	tempDir := testutils.GetTmpDir(t)

	// Exec
	_, _, _ = testutils.RunAppCommand(fmt.Sprintf("set puppet 1.2.3 --config %s", filepath.Join(tempDir, ".prm.yaml")), "")
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("get puppet --config %s", filepath.Join(tempDir, ".prm.yaml")), "")

	// Assert
	assert.Contains(t, stdout, "Puppet version is configured to: 1.2.3")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmGet_Backend_EnsureSet(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)
	tempDir := testutils.GetTmpDir(t)

	// Exec
	// TODO: This test should set the value to something other than 'docker' when more options become available, to test
	// this functionality properly
	_, _, _ = testutils.RunAppCommand(fmt.Sprintf("set backend docker --config %s", filepath.Join(tempDir, ".prm.yaml")), "")
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("get backend --config %s", filepath.Join(tempDir, ".prm.yaml")), "")

	// Assert
	assert.Contains(t, stdout, "Backend is configured to: docker")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmGet_InvalidArg(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("get foo", "")

	// Assert
	assert.Contains(t, stdout, "Displays the requested configuration value\n\nUsage:\n  prm get")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}
