package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssertServerVersionLaterVersionShouldPass(t *testing.T) {
	testServerVersionShouldPass(t, "2025.1", "2025.1")
	testServerVersionShouldPass(t, "2025.1", "2025.1.24564")
	testServerVersionShouldPass(t, "2025.1", "2025.2")
	testServerVersionShouldPass(t, "2025.1.0", "2025.1")
}

func TestAssertServerVersionOlderVersionShouldFail(t *testing.T) {
	testServerVersionShouldFail(t, "2025.1", "2024.4")
	testServerVersionShouldFail(t, "2025.1.12345", "2025.1.12344")
	testServerVersionShouldFail(t, "2025.1.7108", "2025.1")
}

func TestAssertServerVersionLocalAlwaysPasses(t *testing.T) {
	testServerVersionShouldPass(t, "2025.1", "0.0.0-local")
	testServerVersionShouldPass(t, "2020.3.023498", "0.0.0-local")
}

func TestAssertServerVersionInvalidVersionAlwaysPasses(t *testing.T) {
	// When version cannot be parsed - pass and let Server to 'fail' when resource with incompatible version is used
	testServerVersionShouldPass(t, "2025.not-a-number", "0.0.0-local")
	testServerVersionShouldPass(t, "2025.2", "local")
}

func testServerVersionShouldPass(t *testing.T, minVersion string, currentVersion string) {
	pass := isCurrentVersionSameOrLaterThanMinimum(currentVersion, minVersion)
	assert.True(t, pass, "Expected %s to be same or later than %s", currentVersion, minVersion)
}

func testServerVersionShouldFail(t *testing.T, minVersion string, currentVersion string) {
	pass := isCurrentVersionSameOrLaterThanMinimum(currentVersion, minVersion)
	assert.False(t, pass, "Expected %s to be older than %s", currentVersion, minVersion)
}
