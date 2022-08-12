package octopusdeploy

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/stretchr/testify/require"
)

func TestExpandDeploymentAction(t *testing.T) {
	actual := expandAction(nil)
	require.Nil(t, actual)

	actual = expandAction(map[string]interface{}{})
	require.Nil(t, actual)

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	actionType := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	flattened := map[string]interface{}{
		"action_type": actionType,
		"name":        name,
	}

	expected := deployments.NewDeploymentAction(name, actionType)

	actual = expandAction(flattened)
	require.Equal(t, expected, actual)
}

func TestFlattenDeploymentAction(t *testing.T) {
	actual := flattenAction(nil)
	require.Nil(t, actual)

	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	actionType := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	expanded := deployments.NewDeploymentAction(name, actionType)

	actual = flattenAction(expanded)
	expected := map[string]interface{}{
		"can_be_used_for_project_versioning": false,
		"is_disabled":                        false,
		"is_required":                        false,
		"name":                               name,
	}
	require.Equal(t, expected, actual)

	expanded.CanBeUsedForProjectVersioning = true
	expanded.IsDisabled = true
	expanded.IsRequired = true
	actual = flattenAction(expanded)
	expected = map[string]interface{}{
		"can_be_used_for_project_versioning": true,
		"is_disabled":                        true,
		"is_required":                        true,
		"name":                               name,
	}
	require.Equal(t, expected, actual)

	expanded.Channels = append(expanded.Channels, "channel")
	expanded.Condition = "condition"
	actual = flattenAction(expanded)
	expected = map[string]interface{}{
		"can_be_used_for_project_versioning": true,
		"channels":                           []string{"channel"},
		"condition":                          "condition",
		"is_disabled":                        true,
		"is_required":                        true,
		"name":                               name,
	}
	require.Equal(t, expected, actual)
}
