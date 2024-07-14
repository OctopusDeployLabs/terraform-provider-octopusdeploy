package octopusdeploy_framework

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExpandLifecycleWithNil(t *testing.T) {
	lifecycle := expandLifecycle(nil)
	require.Nil(t, lifecycle)
}

func TestExpandLifecycle(t *testing.T) {
	description := "test-description"
	name := "test-name"
	spaceID := "test-space-id"
	Id := "test-id"
	releaseRetention := core.NewRetentionPeriod(0, "Days", true)
	tentacleRetention := core.NewRetentionPeriod(2, "Items", false)

	data := &lifecycleTypeResourceModel{
		ID:          types.StringValue(Id),
		Description: types.StringValue(description),
		Name:        types.StringValue(name),
		SpaceID:     types.StringValue(spaceID),
		ReleaseRetentionPolicy: types.ListValueMust(
			types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()},
			[]attr.Value{
				types.ObjectValueMust(
					getRetentionPeriodAttrTypes(),
					map[string]attr.Value{
						"quantity_to_keep":    types.Int64Value(int64(releaseRetention.QuantityToKeep)),
						"should_keep_forever": types.BoolValue(releaseRetention.ShouldKeepForever),
						"unit":                types.StringValue(releaseRetention.Unit),
					},
				),
			},
		),
		TentacleRetentionPolicy: types.ListValueMust(
			types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()},
			[]attr.Value{
				types.ObjectValueMust(
					getRetentionPeriodAttrTypes(),
					map[string]attr.Value{
						"quantity_to_keep":    types.Int64Value(int64(tentacleRetention.QuantityToKeep)),
						"should_keep_forever": types.BoolValue(tentacleRetention.ShouldKeepForever),
						"unit":                types.StringValue(tentacleRetention.Unit),
					},
				),
			},
		),
	}

	lifecycle := expandLifecycle(data)

	require.Equal(t, description, lifecycle.Description)
	require.NotEmpty(t, lifecycle.ID)
	require.NotNil(t, lifecycle.Links)
	require.Empty(t, lifecycle.Links)
	require.Equal(t, name, lifecycle.Name)
	require.Empty(t, lifecycle.Phases)
	require.Equal(t, releaseRetention, lifecycle.ReleaseRetentionPolicy)
	require.Equal(t, tentacleRetention, lifecycle.TentacleRetentionPolicy)
	require.Equal(t, spaceID, lifecycle.SpaceID)
}
