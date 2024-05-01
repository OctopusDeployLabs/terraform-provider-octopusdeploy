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

func TestExpandDeploymentActionSlugPrimaryPackages(t *testing.T) {
	actual := expandDeploymentActionSlugPrimaryPackages(nil)
	require.Nil(t, actual)

	actual = expandDeploymentActionSlugPrimaryPackages([]interface{}{})
	expected := []packages.DeploymentActionSlugPackage{}
	require.Equal(t, expected, actual)

	flattened := []interface{}{
		map[string]interface{}{
			"deployment_action_slug": "",
		},
		map[string]interface{}{
			"deployment_action_slug": "test-deployment_action",
			"package_reference":      "should_be_ignored",
		},
	}
	expected = []packages.DeploymentActionSlugPackage{
		{DeploymentActionSlug: "", PackageReference: ""},
		{DeploymentActionSlug: "test-deployment_action", PackageReference: ""},
	}
	actual = expandDeploymentActionSlugPrimaryPackages(flattened)
	require.Equal(t, expected, actual)
}

func TestFlattenDeploymentActionSlugPackages(t *testing.T) {
	actual := flattenDeploymentActionSlugPackages(nil)
	require.Nil(t, actual)

	actual = flattenDeploymentActionSlugPackages([]packages.DeploymentActionSlugPackage{})
	require.Nil(t, actual)

	expanded := []packages.DeploymentActionSlugPackage{
		{
			DeploymentActionSlug: "action-one",
			PackageReference:     "",
		},
		{
			DeploymentActionSlug: "action-two",
			PackageReference:     "",
		},
	}
	actual = flattenDeploymentActionSlugPackages(expanded)
	expected := []interface{}{}
	require.Equal(t, expected, actual)

	expanded = getCommonExpandedPackages()
	actual = flattenDeploymentActionSlugPackages(expanded)
	expected = []interface{}{
		map[string]interface{}{
			"deployment_action_slug": "action-one",
			"package_reference":      "some-package",
		},
	}
	require.Equal(t, expected, actual)

}

func TestFlattenDeploymentActionSlugPrimaryPackages(t *testing.T) {
	actual := flattenDeploymentActionSlugPrimaryPackages(nil)
	require.Nil(t, actual)

	actual = flattenDeploymentActionSlugPrimaryPackages([]packages.DeploymentActionSlugPackage{})
	require.Nil(t, actual)

	expanded := []packages.DeploymentActionSlugPackage{
		{
			DeploymentActionSlug: "action-one",
			PackageReference:     "some-package",
		},
		{
			DeploymentActionSlug: "action-two",
			PackageReference:     "some-other-package",
		},
	}
	actual = flattenDeploymentActionSlugPrimaryPackages(expanded)
	expected := []interface{}{}

	expanded = getCommonExpandedPackages()
	actual = flattenDeploymentActionSlugPrimaryPackages(expanded)
	expected = []interface{}{
		map[string]interface{}{
			"deployment_action_slug": "action-two",
		},
		map[string]interface{}{
			"deployment_action_slug": "action-three",
		},
	}
	require.Equal(t, expected, actual)
}

func getCommonExpandedPackages() []packages.DeploymentActionSlugPackage {
	return []packages.DeploymentActionSlugPackage{
		{
			DeploymentActionSlug: "action-one",
			PackageReference:     "some-package",
		},
		{
			DeploymentActionSlug: "action-two",
			PackageReference:     "",
		},
		{
			DeploymentActionSlug: "action-three",
		},
	}
}
