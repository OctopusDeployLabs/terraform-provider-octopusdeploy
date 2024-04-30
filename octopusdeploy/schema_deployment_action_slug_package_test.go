package octopusdeploy

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/stretchr/testify/require"
)

func TestExpandDeploymentActionSlugPackages(t *testing.T) {
	actual := expandDeploymentActionSlugPackages(nil)
	require.Nil(t, actual)

	actual = expandDeploymentActionSlugPackages([]interface{}{})
	expected := []packages.DeploymentActionSlugPackage{}
	require.Equal(t, expected, actual)

	flattened := []interface{}{
		map[string]interface{}{
			"deployment_action_slug": "",
			"package_reference":      "",
		},
		map[string]interface{}{
			"deployment_action_slug": "test-deployment_action",
			"package_reference":      "test-package_reference",
		},
	}
	expected = []packages.DeploymentActionSlugPackage{
		{DeploymentActionSlug: "", PackageReference: ""},
		{DeploymentActionSlug: "test-deployment_action", PackageReference: "test-package_reference"},
	}
	actual = expandDeploymentActionSlugPackages(flattened)
	require.Equal(t, expected, actual)
}

func TestFlattenDeploymentActionSlugPackages(t *testing.T) {
	actual := flattenDeploymentActionSlugPackages(nil)
	require.Nil(t, actual)

	actual = flattenDeploymentActionSlugPackages([]packages.DeploymentActionSlugPackage{})
	require.Nil(t, actual)

	expanded := []packages.DeploymentActionSlugPackage{{}, {
		DeploymentActionSlug: "test-deployment_action",
		PackageReference:     "test-package_reference",
	}}
	actual = flattenDeploymentActionSlugPackages(expanded)
	expected := []interface{}{
		map[string]interface{}{
			"deployment_action_slug": "",
			"package_reference":      "",
		},
		map[string]interface{}{
			"deployment_action_slug": "test-deployment_action",
			"package_reference":      "test-package_reference",
		},
	}
	require.Equal(t, expected, actual)
}
