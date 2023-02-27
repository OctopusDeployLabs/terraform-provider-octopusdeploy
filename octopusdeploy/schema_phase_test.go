package octopusdeploy

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/stretchr/testify/require"
)

func TestExpandPhaseWithNil(t *testing.T) {
	phase := expandPhase(nil)
	require.Nil(t, phase)
}

func TestExpandPhasesWithNil(t *testing.T) {
	phase := expandPhases(nil)
	require.Nil(t, phase)
}

func TestExpandPhaseWithEmptyInput(t *testing.T) {
	phase := expandPhase(map[string]interface{}{})
	require.Nil(t, phase)
}

func TestExpandPhasesWithEmptyInput(t *testing.T) {
	phase := expandPhases([]interface{}{})
	require.Nil(t, phase)
}

func TestFlattenPhaseWithNil(t *testing.T) {
	phase := flattenPhase(nil)
	require.Nil(t, phase)
}

func TestExpandPhaseWithSensibleDefaults(t *testing.T) {
	automaticDeploymentTargets := []string{
		acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
	}
	isOptionalPhase := true
	minimumEnvironmentsBeforePromotion := int32(acctest.RandIntRange(1, 1000))
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	releaseRetentionPolicy := core.NewRetentionPeriod(15, "Items", false)
	tentacleRetentionPolicy := core.NewRetentionPeriod(5, "Days", true)

	actualPhase := lifecycles.NewPhase(name)
	actualPhase.AutomaticDeploymentTargets = automaticDeploymentTargets
	actualPhase.IsOptionalPhase = isOptionalPhase
	actualPhase.MinimumEnvironmentsBeforePromotion = minimumEnvironmentsBeforePromotion
	actualPhase.ReleaseRetentionPolicy = releaseRetentionPolicy
	actualPhase.TentacleRetentionPolicy = tentacleRetentionPolicy

	flattenedPhase := flattenPhase(actualPhase)

	phase := expandPhase(flattenedPhase)

	require.NotNil(t, phase.ID)
	require.Equal(t, automaticDeploymentTargets, phase.AutomaticDeploymentTargets)
	require.Equal(t, isOptionalPhase, phase.IsOptionalPhase)
	require.EqualValues(t, minimumEnvironmentsBeforePromotion, phase.MinimumEnvironmentsBeforePromotion)
	require.Equal(t, name, phase.Name)
	require.Equal(t, releaseRetentionPolicy, phase.ReleaseRetentionPolicy)
	require.Equal(t, tentacleRetentionPolicy, phase.TentacleRetentionPolicy)
}

func TestExpandPhasesWithSensibleDefaults(t *testing.T) {
	flattenedPhases := []interface{}{}

	automaticDeploymentTargets := []string{
		acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
	}
	isOptionalPhase := true
	minimumEnvironmentsBeforePromotion := int32(acctest.RandIntRange(1, 1000))
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	releaseRetentionPolicy := core.NewRetentionPeriod(15, "Items", false)
	tentacleRetentionPolicy := core.NewRetentionPeriod(5, "Days", true)

	actualPhase := lifecycles.NewPhase(name)
	actualPhase.AutomaticDeploymentTargets = automaticDeploymentTargets
	actualPhase.IsOptionalPhase = isOptionalPhase
	actualPhase.MinimumEnvironmentsBeforePromotion = minimumEnvironmentsBeforePromotion
	actualPhase.ReleaseRetentionPolicy = releaseRetentionPolicy
	actualPhase.TentacleRetentionPolicy = tentacleRetentionPolicy

	flattenedPhases = append(flattenedPhases, flattenPhase(actualPhase))

	phases := expandPhases(flattenedPhases)

	require.NotNil(t, phases)
	require.Len(t, phases, 1)

	automaticDeploymentTargets = []string{
		acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
		acctest.RandStringFromCharSet(20, acctest.CharSetAlpha),
	}
	isOptionalPhase = true
	minimumEnvironmentsBeforePromotion = int32(acctest.RandIntRange(1, 1000))
	name = acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	releaseRetentionPolicy = core.NewRetentionPeriod(15, "Items", false)
	tentacleRetentionPolicy = core.NewRetentionPeriod(5, "Days", true)

	actualPhase = lifecycles.NewPhase(name)
	actualPhase.AutomaticDeploymentTargets = automaticDeploymentTargets
	actualPhase.IsOptionalPhase = isOptionalPhase
	actualPhase.MinimumEnvironmentsBeforePromotion = minimumEnvironmentsBeforePromotion
	actualPhase.ReleaseRetentionPolicy = releaseRetentionPolicy
	actualPhase.TentacleRetentionPolicy = tentacleRetentionPolicy

	flattenedPhases = append(flattenedPhases, flattenPhase(actualPhase))

	phases = expandPhases(flattenedPhases)

	require.NotNil(t, phases)
	require.Len(t, phases, 2)
}
