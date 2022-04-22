package status_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/puppetlabs/pct/acceptance/testutils"
	"github.com/stretchr/testify/assert"
)

const APP = "prm"

func Test_PrmStatus_NoArgs(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	if runtime.GOOS == "darwin" {
		t.Skip("Docker based acceptance tests currently fail on MacOS")
	}

	// Setup
	testutils.SetAppName(APP)
	tempDir := testutils.GetTmpDir(t)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("status --config %s", filepath.Join(tempDir, ".prm.yaml")), "")

	assert.Contains(t, stdout, "Puppet version: 7.15.0")
	assert.Contains(t, stdout, "Backend: docker (running)")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_PrmStatus_Json(t *testing.T) {
	testutils.SkipAcceptanceTest(t)
	if runtime.GOOS == "darwin" {
		t.Skip("Docker based acceptance tests currently fail on MacOS")
	}

	// Setup
	testutils.SetAppName(APP)
	tempDir := testutils.GetTmpDir(t)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand(fmt.Sprintf("status --format json --config %s", filepath.Join(tempDir, ".prm.yaml")), "")

	assert.Contains(t, stdout, `"PuppetVersion":"7.15.0"`)
	assert.Contains(t, stdout, `"Backend":"docker"`)
	assert.Contains(t, stdout, `"IsAvailable":true`)
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}
