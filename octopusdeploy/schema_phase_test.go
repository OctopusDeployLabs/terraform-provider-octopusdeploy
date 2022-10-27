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

	resourceMap := map[string]interface{}{
		"automatic_deployment_targets":          automaticDeploymentTargets,
		"is_optional_phase":                     isOptionalPhase,
		"minimum_environments_before_promotion": minimumEnvironmentsBeforePromotion,
		"name":                                  name,
	}

	phase := expandPhase(resourceMap)

	require.NotNil(t, phase.ID)
	require.Equal(t, automaticDeploymentTargets, phase.AutomaticDeploymentTargets)
	require.Equal(t, isOptionalPhase, phase.IsOptionalPhase)
	require.EqualValues(t, minimumEnvironmentsBeforePromotion, phase.MinimumEnvironmentsBeforePromotion)
	require.Equal(t, name, phase.Name)
	require.Nil(t, phase.ReleaseRetentionPolicy)
	require.Nil(t, phase.TentacleRetentionPolicy)
}
