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

	resourceMap := map[string]interface{}{
		"description": description,
		"name":        name,
		"space_id":    spaceID,
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
	require.NotNil(t, lifecycle.ReleaseRetentionPolicy)
	require.Equal(t, lifecycle.SpaceID, spaceID)
	require.NotNil(t, lifecycle.TentacleRetentionPolicy)
}
