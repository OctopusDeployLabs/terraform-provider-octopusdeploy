package octopusdeploy_framework

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbookprocess"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccMapProcessStepsOrderingFromStateAddMissingStepToTheEnd(t *testing.T) {
	state := schemas.ProcessStepsOrderResourceModel{
		SpaceID:   types.StringValue("Spaces-1"),
		ProcessID: types.StringValue("Processes-1"),
		Steps: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("step-1"),
			types.StringValue("step-2"),
			types.StringValue("step-3"),
		}),
	}

	step1 := deployments.NewDeploymentStep("Step One")
	step1.SetID("step-1")

	step2 := deployments.NewDeploymentStep("Step Two")
	step2.SetID("step-2")

	step3 := deployments.NewDeploymentStep("Step Three")
	step3.SetID("step-3")

	step4 := deployments.NewDeploymentStep("Step Four")
	step4.SetID("step-4")

	process := deployments.DeploymentProcess{
		Steps: []*deployments.DeploymentStep{step4, step3, step2, step1},
	}

	diags := mapProcessStepsOrderFromState(&state, deploymentProcessWrapper{&process})

	orderedStepIds := make([]string, len(process.Steps))
	for i, step := range process.Steps {
		orderedStepIds[i] = step.ID
	}
	expectedStepIds := []string{"step-1", "step-2", "step-3", "step-4"}
	assert.Equal(t, expectedStepIds, orderedStepIds, "Should put missing steps to the end of the order")

	diagnostics := make([]diag.Severity, len(diags))
	for i, d := range diags {
		diagnostics[i] = d.Severity()
	}
	expectedDiagnostics := []diag.Severity{diag.SeverityWarning}
	assert.Equal(t, expectedDiagnostics, diagnostics, "Expects to have warning diagnostics")
}

func TestAccMapProcessStepsOrderingFromStateAddMissingStepToTheEndForRunbook(t *testing.T) {
	state := schemas.ProcessStepsOrderResourceModel{
		SpaceID:   types.StringValue("Spaces-1"),
		ProcessID: types.StringValue("Processes-1"),
		Steps: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("step-1"),
			types.StringValue("step-2"),
			types.StringValue("step-3"),
		}),
	}

	step1 := deployments.NewDeploymentStep("Step One")
	step1.SetID("step-1")

	step2 := deployments.NewDeploymentStep("Step Two")
	step2.SetID("step-2")

	step3 := deployments.NewDeploymentStep("Step Three")
	step3.SetID("step-3")

	step4 := deployments.NewDeploymentStep("Step Four")
	step4.SetID("step-4")

	process := runbookprocess.RunbookProcess{
		Steps: []*deployments.DeploymentStep{step4, step3, step2, step1},
	}

	diags := mapProcessStepsOrderFromState(&state, runbookProcessWrapper{&process})

	orderedStepIds := make([]string, len(process.Steps))
	for i, step := range process.Steps {
		orderedStepIds[i] = step.ID
	}
	expectedStepIds := []string{"step-1", "step-2", "step-3", "step-4"}
	assert.Equal(t, expectedStepIds, orderedStepIds, "Should put missing steps to the end of the order")

	diagnostics := make([]diag.Severity, len(diags))
	for i, d := range diags {
		diagnostics[i] = d.Severity()
	}
	expectedDiagnostics := []diag.Severity{diag.SeverityWarning}
	assert.Equal(t, expectedDiagnostics, diagnostics, "Expects to have warning diagnostics")
}

func TestAccMapProcessStepsOrderingFromStateAddsErrorWhenStepIsNotPartOfTheProcess(t *testing.T) {
	state := schemas.ProcessStepsOrderResourceModel{
		SpaceID:   types.StringValue("Spaces-1"),
		ProcessID: types.StringValue("Processes-1"),
		Steps: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("step-1"),
			types.StringValue("step-2"),
			types.StringValue("step-3"),
		}),
	}

	step1 := deployments.NewDeploymentStep("Step One")
	step1.SetID("step-1")

	step2 := deployments.NewDeploymentStep("Step Two")
	step2.SetID("step-2")

	step4 := deployments.NewDeploymentStep("Step Four (Unordered)")
	step4.SetID("step-4")

	process := deployments.DeploymentProcess{
		Steps: []*deployments.DeploymentStep{step4, step2, step1},
	}

	diags := mapProcessStepsOrderFromState(&state, deploymentProcessWrapper{&process})

	orderedStepIds := make([]string, len(process.Steps))
	for i, step := range process.Steps {
		orderedStepIds[i] = step.ID
	}
	expectedStepIds := []string{"step-4", "step-2", "step-1"}
	assert.Equal(t, expectedStepIds, orderedStepIds, "Should not update original order when invalid steps are found")

	diagnostics := make([]diag.Severity, len(diags))
	for i, d := range diags {
		diagnostics[i] = d.Severity()
	}
	expectedDiagnostics := []diag.Severity{diag.SeverityError, diag.SeverityWarning}
	assert.Equal(t, expectedDiagnostics, diagnostics, "Expects to have an error about invalid step and warning about not included step")
}

func TestAccMapProcessStepsOrderingToStateTakesOnlyConfiguredAmountOfSteps(t *testing.T) {
	step1 := deployments.NewDeploymentStep("Step One")
	step1.SetID("00000000-0000-0000-0000-000000000001")

	step2 := deployments.NewDeploymentStep("Step Two")
	step2.SetID("00000000-0000-0000-0000-000000000002")

	step3 := deployments.NewDeploymentStep("Step Three")
	step3.SetID("00000000-0000-0000-0000-000000000003")

	process := &deployments.DeploymentProcess{
		SpaceID:   "Spaces-1",
		ProjectID: "Projects-1",
		Steps:     []*deployments.DeploymentStep{step1, step2, step3},
	}
	process.SetID("Processes-1")

	state := schemas.ProcessStepsOrderResourceModel{
		Steps: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(step3.ID),
			types.StringValue(step2.ID),
		}),
	}

	mapProcessStepsOrderToState(deploymentProcessWrapper{process}, &state)

	expectedState := schemas.ProcessStepsOrderResourceModel{
		SpaceID:   types.StringValue(process.SpaceID),
		ProcessID: types.StringValue(process.ID),
		Steps: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(step1.ID),
			types.StringValue(step2.ID),
		}),
	}
	expectedState.ID = types.StringValue(process.ID)

	assert.Equal(t, expectedState, state)
}

func TestAccMapProcessStepsOrderingToStateTakesOnlyConfiguredAmountOfStepsForRunbooks(t *testing.T) {
	step1 := deployments.NewDeploymentStep("Step One")
	step1.SetID("00000000-0000-0000-0000-000000000001")

	step2 := deployments.NewDeploymentStep("Step Two")
	step2.SetID("00000000-0000-0000-0000-000000000002")

	step3 := deployments.NewDeploymentStep("Step Three")
	step3.SetID("00000000-0000-0000-0000-000000000003")

	process := &runbookprocess.RunbookProcess{
		SpaceID:   "Spaces-1",
		ProjectID: "Projects-1",
		RunbookID: "Runbook-1",
		Steps:     []*deployments.DeploymentStep{step1, step2, step3},
	}
	process.SetID("Processes-1")

	state := schemas.ProcessStepsOrderResourceModel{
		Steps: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(step3.ID),
			types.StringValue(step2.ID),
		}),
	}

	mapProcessStepsOrderToState(runbookProcessWrapper{process}, &state)

	expectedState := schemas.ProcessStepsOrderResourceModel{
		SpaceID:   types.StringValue(process.SpaceID),
		ProcessID: types.StringValue(process.ID),
		Steps: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(step1.ID),
			types.StringValue(step2.ID),
		}),
	}
	expectedState.ID = types.StringValue(process.ID)

	assert.Equal(t, expectedState, state)
}

func TestAccMapProcessStepsOrderingToStateWhenConfiguredStepsMoreThanProvided(t *testing.T) {
	step1 := deployments.NewDeploymentStep("Step One")
	step1.SetID("00000000-0000-0000-0000-000000000001")

	step2 := deployments.NewDeploymentStep("Step Two")
	step2.SetID("00000000-0000-0000-0000-000000000002")

	step3 := deployments.NewDeploymentStep("Step Three")
	step3.SetID("00000000-0000-0000-0000-000000000003")

	process := &deployments.DeploymentProcess{
		SpaceID:   "Spaces-1",
		ProjectID: "Projects-1",
		Steps:     []*deployments.DeploymentStep{step1, step2},
	}
	process.SetID("Processes-1")

	state := schemas.ProcessStepsOrderResourceModel{
		Steps: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(step1.ID),
			types.StringValue(step2.ID),
			types.StringValue(step3.ID),
		}),
	}

	mapProcessStepsOrderToState(deploymentProcessWrapper{process}, &state)

	expectedState := schemas.ProcessStepsOrderResourceModel{
		SpaceID:   types.StringValue(process.SpaceID),
		ProcessID: types.StringValue(process.ID),
		Steps: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(step1.ID),
			types.StringValue(step2.ID),
		}),
	}
	expectedState.ID = types.StringValue(process.ID)

	assert.Equal(t, expectedState, state)
}
