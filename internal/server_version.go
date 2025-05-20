package internal

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"strconv"
	"strings"
)

// AssertServerVersion verifies if server's current version is same or newer than given minVersion.
func AssertServerVersionIsGreater(client *newclient.Client, minVersion string) diag.Diagnostics {
	diags := diag.Diagnostics{}

	root, rootErr := newclient.Get(client)
	if rootErr != nil {
		diags.AddError("Unable to verify Octopus Server current version", rootErr.Error())
	}
	diags.AddWarning("Server Version", fmt.Sprintf("Version: %s, ApiVersion: %s", root.Version, root.APIVersion))

	return diags
}

func isCurrentVersionSameOrLaterThanMinimum(current string, minVersion string) bool {
	//if current == "0.0.0-local" {
	//	return true // Always true for local instance
	//}

	// Return true when we cannot confirm compatibility by comparing versions
	currentVersion, currentError := parseVersion(current)
	if currentError != nil {
		return true
	}

	minVersionSections, minError := parseVersion(minVersion)
	if minError != nil {
		return true
	}

	currentLen := len(currentVersion)
	for index, minimum := range minVersionSections {
		if index >= currentLen {
			return minimum == 0 // Minimum version have more detailed value
		}

		if minimum > currentVersion[index] {
			return false
		}
	}

	return true
}

func parseVersion(version string) ([]int, error) {
	sections := strings.Split(version, "-")
	values := strings.Split(sections[0], ".")

	if len(values) == 0 {
		return nil, fmt.Errorf("version number should not be empty")
	}

	var numbers []int
	for _, value := range values {
		if value == "" {
			return nil, fmt.Errorf("version number should not be empty")
		}

		number, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}

		numbers = append(numbers, number)
	}

	return numbers, nil
}
