package set_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/puppetlabs/pct/acceptance/testutils"
	"github.com/stretchr/testify/assert"
)

const APP = "prm"

func Test_PrmSet_NoArgs(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("set", "")

	// Assert
	assert.Contains(t, stdout, "Sets the specified configuration to the specified value\n\nUsage:\n  prm set")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmSet_Puppet_NoArgs(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("set puppet", "")

	// Assert
	assert.Equal(t, "Error: please specify a Puppet version after 'set puppet'\n", stdout)
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)
}

func Test_PrmSet_Puppet_ValidVer(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)
	tempDir := testutils.GetTmpDir(t)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("set puppet 1.2.3 --config %s", filepath.Join(tempDir, ".prm.yaml")), "")

	// Assert
	assert.Empty(t, stdout)
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)

	stdout, stderr, exitCode = testutils.RunAppCommand(fmt.Sprintf("get puppet --config %s", filepath.Join(tempDir, ".prm.yaml")), "")
	assert.Contains(t, stdout, "Puppet version is configured to: 1.2.3")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmSet_Puppet_InvalidVer(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)
	tempDir := testutils.GetTmpDir(t)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("set puppet foo.bar --config %s", filepath.Join(tempDir, ".prm.yaml")), "")

	// Assert
	assert.Contains(t, stdout, "'foo.bar' is not a semantic (x.y.z) Puppet version: Invalid Semantic Version")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)

	stdout, stderr, exitCode = testutils.RunAppCommand(fmt.Sprintf("get puppet --config %s", filepath.Join(tempDir, ".prm.yaml")), "")
	assert.Contains(t, stdout, "Puppet version is configured to: 7.15.0")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmSet_Backend_InvalidOpt(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)
	tempDir := testutils.GetTmpDir(t)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("set backend foo.bar --config %s", filepath.Join(tempDir, ".prm.yaml")), "")

	// Assert
	assert.Contains(t, stdout, "'foo.bar' is not a valid backend type, please specify one of the following backend types:")
	assert.Equal(t, "exit status 1", stderr)
	assert.Equal(t, 1, exitCode)

	stdout, stderr, exitCode = testutils.RunAppCommand(fmt.Sprintf("get backend --config %s", filepath.Join(tempDir, ".prm.yaml")), "")
	assert.Contains(t, stdout, "Backend is configured to: docker")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}
