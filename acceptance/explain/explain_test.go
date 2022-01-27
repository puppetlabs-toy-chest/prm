package explain_test

import (
	"testing"

	"github.com/puppetlabs/pdkgo/acceptance/testutils"
	"github.com/stretchr/testify/assert"
)

const APP = "prm"

func Test_Explain_With_No_Args(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("explain", "")

	// Assert
	assert.Contains(t, stdout, "telemetry")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_Explain_With_List_Set(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("explain --list", "")

	// Assert
	assert.Contains(t, stdout, "telemetry")
	assert.Contains(t, stdout, "quick-start")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_Explain_With_Format_Json(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("explain --format json", "")

	// Assert
	assert.Contains(t, stdout, "\"Title\":{\"Short\":\"telemetry\",\"Long\":\"Telemetry\"")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_Explain_With_Tag_Filter(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("explain --tag meta", "")

	// Assert
	assert.Contains(t, stdout, "telemetry")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_Explain_With_Cateogry_Filter(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("explain --category concept", "")

	// Assert
	assert.Contains(t, stdout, "telemetry")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}

func Test_Explain_With_Single_Target(t *testing.T) {
	testutils.SkipAcceptanceTest(t)

	// Setup
	testutils.SetAppName(APP)

	// Exec
	stdout, stderr, exitCode := testutils.RunAppCommand("explain telemetry", "")

	// Assert
	assert.Contains(t, stdout, "Telemetry")
	assert.Empty(t, stderr)
	assert.Equal(t, 0, exitCode)
}
