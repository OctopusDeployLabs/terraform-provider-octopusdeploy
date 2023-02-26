package octopusdeploy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestExpandLifecycleWithNil(t *testing.T) {
	lifecycle := expandLifecycle(nil)
	require.Nil(t, lifecycle)
}

func TestExpandLifecycle(t *testing.T) {
	description := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	spaceID := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	releaseRetention := []interface{}{
		map[string]interface{}{
			"quantity_to_keep":    0,
			"should_keep_forever": true,
			"unit":                "Days",
		}}
	tentacleRetention := []interface{}{
		map[string]interface{}{
			"quantity_to_keep":    2,
			"should_keep_forever": false,
			"unit":                "Items",
		}}

	resourceMap := map[string]interface{}{
		"description":               description,
		"name":                      name,
		"space_id":                  spaceID,
		"release_retention_policy":  releaseRetention,
		"tentacle_retention_policy": tentacleRetention,
	}

	d := schema.TestResourceDataRaw(t, getLifecycleSchema(), resourceMap)
	lifecycle := expandLifecycle(d)

	require.Equal(t, lifecycle.Description, description)
	require.NotNil(t, lifecycle.ID)
	require.NotNil(t, lifecycle.Links)
	require.Empty(t, lifecycle.Links)
	require.NotNil(t, lifecycle.ModifiedBy)
	require.Nil(t, lifecycle.ModifiedOn)
	require.Equal(t, lifecycle.Name, name)
	require.Empty(t, lifecycle.Phases)
	require.EqualValues(t, lifecycle.ReleaseRetentionPolicy.QuantityToKeep, 0)
	require.EqualValues(t, lifecycle.TentacleRetentionPolicy.QuantityToKeep, 2)
	require.EqualValues(t, lifecycle.ReleaseRetentionPolicy.ShouldKeepForever, true)
	require.EqualValues(t, lifecycle.TentacleRetentionPolicy.ShouldKeepForever, false)
	require.EqualValues(t, lifecycle.ReleaseRetentionPolicy.Unit, "Days")
	require.EqualValues(t, lifecycle.TentacleRetentionPolicy.Unit, "Items")
	require.Equal(t, lifecycle.SpaceID, spaceID)
}
