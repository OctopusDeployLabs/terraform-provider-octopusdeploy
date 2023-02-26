package octopusdeploy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/stretchr/testify/require"
)

func TestExpandPhaseWithNil(t *testing.T) {
	phase := expandPhase(nil)
	require.Nil(t, phase)
}

func TestExpandPhase(t *testing.T) {
	automaticDeploymentTargets := []string{
		acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
	}
	isOptionalPhase := true
	minimumEnvironmentsBeforePromotion := acctest.RandIntRange(1, 1000)
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
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
		"automatic_deployment_targets":          automaticDeploymentTargets,
		"is_optional_phase":                     isOptionalPhase,
		"minimum_environments_before_promotion": minimumEnvironmentsBeforePromotion,
		"name":                                  name,
		"release_retention_policy":              releaseRetention,
		"tentacle_retention_policy":             tentacleRetention,
	}

	phase := expandPhase(resourceMap)

	require.NotNil(t, phase.ID)
	require.Equal(t, automaticDeploymentTargets, phase.AutomaticDeploymentTargets)
	require.Equal(t, isOptionalPhase, phase.IsOptionalPhase)
	require.EqualValues(t, minimumEnvironmentsBeforePromotion, phase.MinimumEnvironmentsBeforePromotion)
	require.Equal(t, name, phase.Name)
	require.EqualValues(t, phase.ReleaseRetentionPolicy.QuantityToKeep, 0)
	require.EqualValues(t, phase.TentacleRetentionPolicy.QuantityToKeep, 2)
	require.EqualValues(t, phase.ReleaseRetentionPolicy.ShouldKeepForever, true)
	require.EqualValues(t, phase.TentacleRetentionPolicy.ShouldKeepForever, false)
	require.EqualValues(t, phase.ReleaseRetentionPolicy.Unit, "Days")
	require.EqualValues(t, phase.TentacleRetentionPolicy.Unit, "Items")
}
