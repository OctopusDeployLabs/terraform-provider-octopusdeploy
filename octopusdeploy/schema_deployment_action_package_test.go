package octopusdeploy

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/stretchr/testify/require"
)

func TestExpandDeploymentActionPackage(t *testing.T) {
	actual := expandDeploymentActionPackage(nil)
	require.Nil(t, actual)

	actual = expandDeploymentActionPackage([]interface{}{})
	require.Nil(t, actual)

	flattened := []interface{}{
		map[string]interface{}{
			"deployment_action": "",
			"package_reference": "",
		},
	}

	expected := &octopusdeploy.DeploymentActionPackage{
		DeploymentAction: "",
		PackageReference: "",
	}

	actual = expandDeploymentActionPackage(flattened)
	require.Equal(t, expected, actual)

	flattened = []interface{}{
		map[string]interface{}{
			"deployment_action": "test-deployment_action",
			"package_reference": "test-package_reference",
		},
	}

	actual = expandDeploymentActionPackage(flattened)
	expected = &octopusdeploy.DeploymentActionPackage{
		DeploymentAction: "test-deployment_action",
		PackageReference: "test-package_reference",
	}
	require.Equal(t, expected, actual)
}

func TestFlattenDeploymentActionPackage(t *testing.T) {
	actual := flattenDeploymentActionPackage(nil)
	require.Nil(t, actual)

	expanded := &octopusdeploy.DeploymentActionPackage{}

	actual = flattenDeploymentActionPackage(expanded)
	expected := []interface{}{
		map[string]interface{}{
			"deployment_action": "",
			"package_reference": "",
		},
	}
	require.Equal(t, expected, actual)

	expanded = &octopusdeploy.DeploymentActionPackage{
		DeploymentAction: "test-deployment_action",
		PackageReference: "test-package_reference",
	}

	actual = flattenDeploymentActionPackage(expanded)
	expected = []interface{}{
		map[string]interface{}{
			"deployment_action": "test-deployment_action",
			"package_reference": "test-package_reference",
		},
	}
	require.Equal(t, expected, actual)
}
