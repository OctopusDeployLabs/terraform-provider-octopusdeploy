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

func TestAccMapProcessChildStepsOrderingFromStateAddMissingChildToTheEnd(t *testing.T) {
	child1 := deployments.NewDeploymentAction("Child 1", "Octopus.Script")
	child1.SetID("child-1")

	child2 := deployments.NewDeploymentAction("Child 2", "Octopus.Script")
	child2.SetID("child-2")

	child3 := deployments.NewDeploymentAction("Child 3", "Octopus.Script")
	child3.SetID("child-3")

	child4 := deployments.NewDeploymentAction("Child 4", "Octopus.Script")
	child4.SetID("child-4")

	parent := deployments.NewDeploymentStep("Parent Step")
	parent.SetID("steps-1")
	parent.Actions = []*deployments.DeploymentAction{child1, child4, child2, child3}

	process := deployments.DeploymentProcess{
		Steps: []*deployments.DeploymentStep{parent},
	}

	state := schemas.ProcessChildStepsOrderResourceModel{
		SpaceID:   types.StringValue(process.SpaceID),
		ProcessID: types.StringValue(process.ID),
		ParentID:  types.StringValue(parent.ID),
		Children: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(child2.ID),
			types.StringValue(child3.ID),
		}),
	}

	diags := mapProcessChildStepsOrderFromState(&state, parent)

	orderedActionIds := make([]string, len(parent.Actions))
	for i, action := range parent.Actions {
		orderedActionIds[i] = action.ID
	}
	expectedActionIds := []string{"child-1", "child-2", "child-3", "child-4"}
	assert.Equal(t, expectedActionIds, orderedActionIds, "Should put missing actions to the end of the order")

	diagnostics := make([]diag.Severity, len(diags))
	for i, d := range diags {
		diagnostics[i] = d.Severity()
	}
	expectedDiagnostics := []diag.Severity{diag.SeverityWarning}
	assert.Equal(t, expectedDiagnostics, diagnostics, "Expects to have warning diagnostics")
}

func TestAccMapProcessChildStepsOrderingFromStateAddsErrorWhenChildIsNotPartOfTheParent(t *testing.T) {
	child1 := deployments.NewDeploymentAction("Child 1", "Octopus.Script")
	child1.SetID("child-1")

	child2 := deployments.NewDeploymentAction("Child 2", "Octopus.Script")
	child2.SetID("child-2")

	child3 := deployments.NewDeploymentAction("Child 3", "Octopus.Script")
	child3.SetID("child-3")

	child4 := deployments.NewDeploymentAction("Child 4", "Octopus.Script")
	child4.SetID("child-4")

	parent := deployments.NewDeploymentStep("Parent Step")
	parent.SetID("steps-1")
	parent.Actions = []*deployments.DeploymentAction{child1, child3, child2}

	process := deployments.DeploymentProcess{
		Steps: []*deployments.DeploymentStep{parent},
	}

	state := schemas.ProcessChildStepsOrderResourceModel{
		SpaceID:   types.StringValue(process.SpaceID),
		ProcessID: types.StringValue(process.ID),
		ParentID:  types.StringValue(parent.ID),
		Children: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(child2.ID),
			types.StringValue(child3.ID),
			types.StringValue(child4.ID),
		}),
	}

	diags := mapProcessChildStepsOrderFromState(&state, parent)

	orderedActionIds := make([]string, len(parent.Actions))
	for i, action := range parent.Actions {
		orderedActionIds[i] = action.ID
	}
	expectedActionIds := []string{"child-1", "child-3", "child-2"}
	assert.Equal(t, expectedActionIds, orderedActionIds, "Should not update original order when invalid steps are found")

	diagnostics := make([]diag.Severity, len(diags))
	for i, d := range diags {
		diagnostics[i] = d.Severity()
	}
	expectedDiagnostics := []diag.Severity{diag.SeverityError}
	assert.Equal(t, expectedDiagnostics, diagnostics, "Expects to have an error about invalid child step")
}

func TestAccMapProcessChildStepsOrderingToStateTakesOnlyConfiguredAmountOfSteps(t *testing.T) {
	child1 := deployments.NewDeploymentAction("Child 1", "Octopus.Script")
	child1.SetID("child-1")

	child2 := deployments.NewDeploymentAction("Child 2", "Octopus.Script")
	child2.SetID("child-2")

	child3 := deployments.NewDeploymentAction("Child 3", "Octopus.Script")
	child3.SetID("child-3")

	child4 := deployments.NewDeploymentAction("Child 4", "Octopus.Script")
	child4.SetID("child-4")

	parent := deployments.NewDeploymentStep("Parent Step")
	parent.SetID("steps-1")
	parent.Actions = []*deployments.DeploymentAction{child1, child2, child3, child4}

	process := deployments.DeploymentProcess{
		Steps: []*deployments.DeploymentStep{parent},
	}

	state := schemas.ProcessChildStepsOrderResourceModel{
		Children: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(child3.ID),
			types.StringValue(child4.ID),
		}),
	}

	mapProcessChildStepsOrderToState(deploymentProcessWrapper{&process}, parent, &state)

	expectedState := schemas.ProcessChildStepsOrderResourceModel{
		SpaceID:   types.StringValue(process.SpaceID),
		ProcessID: types.StringValue(process.ID),
		ParentID:  types.StringValue(parent.ID),
		Children: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(child2.ID),
			types.StringValue(child3.ID),
		}),
	}
	expectedState.ID = types.StringValue(parent.ID)

	assert.Equal(t, expectedState, state)
}

func TestAccMapProcessChildStepsOrderingToStateTakesOnlyConfiguredAmountOfStepsForRunbook(t *testing.T) {
	child1 := deployments.NewDeploymentAction("Child 1", "Octopus.Script")
	child1.SetID("child-1")

	child2 := deployments.NewDeploymentAction("Child 2", "Octopus.Script")
	child2.SetID("child-2")

	child3 := deployments.NewDeploymentAction("Child 3", "Octopus.Script")
	child3.SetID("child-3")

	child4 := deployments.NewDeploymentAction("Child 4", "Octopus.Script")
	child4.SetID("child-4")

	parent := deployments.NewDeploymentStep("Parent Step")
	parent.SetID("steps-1")
	parent.Actions = []*deployments.DeploymentAction{child1, child2, child3, child4}

	process := runbookprocess.RunbookProcess{
		Steps: []*deployments.DeploymentStep{parent},
	}

	state := schemas.ProcessChildStepsOrderResourceModel{
		Children: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(child3.ID),
			types.StringValue(child4.ID),
		}),
	}

	mapProcessChildStepsOrderToState(runbookProcessWrapper{&process}, parent, &state)

	expectedState := schemas.ProcessChildStepsOrderResourceModel{
		SpaceID:   types.StringValue(process.SpaceID),
		ProcessID: types.StringValue(process.ID),
		ParentID:  types.StringValue(parent.ID),
		Children: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(child2.ID),
			types.StringValue(child3.ID),
		}),
	}
	expectedState.ID = types.StringValue(parent.ID)

	assert.Equal(t, expectedState, state)
}

func TestAccMapProcessChildStepsOrderingToStateWhenConfiguredChildrenMoreThanProvided(t *testing.T) {
	child1 := deployments.NewDeploymentAction("Child 1", "Octopus.Script")
	child1.SetID("child-1")

	child2 := deployments.NewDeploymentAction("Child 2", "Octopus.Script")
	child2.SetID("child-2")

	child3 := deployments.NewDeploymentAction("Child 3", "Octopus.Script")
	child3.SetID("child-3")

	child4 := deployments.NewDeploymentAction("Child 4", "Octopus.Script")
	child4.SetID("child-4")

	parent := deployments.NewDeploymentStep("Parent Step")
	parent.SetID("steps-1")
	parent.Actions = []*deployments.DeploymentAction{child1, child2, child3}

	process := deployments.DeploymentProcess{
		Steps: []*deployments.DeploymentStep{parent},
	}

	state := schemas.ProcessChildStepsOrderResourceModel{
		Children: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(child2.ID),
			types.StringValue(child3.ID),
			types.StringValue(child4.ID),
		}),
	}

	mapProcessChildStepsOrderToState(deploymentProcessWrapper{&process}, parent, &state)

	expectedState := schemas.ProcessChildStepsOrderResourceModel{
		SpaceID:   types.StringValue(process.SpaceID),
		ProcessID: types.StringValue(process.ID),
		ParentID:  types.StringValue(parent.ID),
		Children: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue(child2.ID),
			types.StringValue(child3.ID),
		}),
	}
	expectedState.ID = types.StringValue(parent.ID)

	assert.Equal(t, expectedState, state)
}
