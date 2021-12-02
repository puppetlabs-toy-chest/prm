package status_test

import (
	"runtime"
	"testing"

	"github.com/puppetlabs/pdkgo/acceptance/testutils"
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

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("status", "")

	assert.Contains(t, stdout, "Puppet version: 7.0.0")
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

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("status --format json", "")

	assert.Contains(t, stdout, `"PuppetVersion":"7.0.0"`)
	assert.Contains(t, stdout, `"Backend":"docker"`)
	assert.Contains(t, stdout, `"IsAvailable":true`)
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}
