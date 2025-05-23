package octopusdeploy_framework

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnsureResourceCompatibilityByVersionNewerOrSameShouldPass(t *testing.T) {
	testServerVersionShouldPass(t, "2025.1", "2025.1")
	testServerVersionShouldPass(t, "2025.1", "2025.1.24564")
	testServerVersionShouldPass(t, "2025.1", "2025.2")
	testServerVersionShouldPass(t, "2025.1.0", "2025.1")
}

func TestEnsureResourceCompatibilityByVersionOlderShouldFail(t *testing.T) {
	testServerVersionShouldFail(t, "2025.1", "2024.4")
	testServerVersionShouldFail(t, "2025.1.12345", "2025.1.12344")
	testServerVersionShouldFail(t, "2025.1.7108", "2025.1")
}

func TestEnsureResourceCompatibilityByVersionLocalAlwaysPasses(t *testing.T) {
	testServerVersionShouldPass(t, "2025.1", "0.0.0-local")
	testServerVersionShouldPass(t, "2020.3.023498", "0.0.0-local")
}

func TestEnsureResourceCompatibilityByVersionTreatsInvalidVersionAsEmpty(t *testing.T) {
	testServerVersionShouldPass(t, "2025.not-a-number", "2024.4.20456")
	testServerVersionShouldFail(t, "2025.2", "local")
	testServerVersionShouldFail(t, "2025.0.1501", "2025.local.1500")
}

func TestEnsureResourceCompatibilityByVersionIgnoresBranches(t *testing.T) {
	testServerVersionShouldPass(t, "2025.1", "2025.1-my-branch")
	testServerVersionShouldFail(t, "2025.2.100-testing", "2025.2.99-experiment")
}

func testServerVersionShouldPass(t *testing.T, limit string, current string) {
	configuration := Config{OctopusVersion: current}
	diags := configuration.EnsureResourceCompatibilityByVersion("compatible_resource_name", limit)

	assert.False(t, diags.HasError(), "Expected %s to pass limit %s", current, limit)
}

func testServerVersionShouldFail(t *testing.T, limit string, current string) {
	configuration := Config{OctopusVersion: current}
	diags := configuration.EnsureResourceCompatibilityByVersion("incompatible_resource_name", limit)

	assert.True(t, diags.HasError(), "Expected %s to fail limit %s", current, limit)
}
