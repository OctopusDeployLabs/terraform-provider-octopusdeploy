package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccMapProcessStepFromStateWithAllAttributes(t *testing.T) {
	ctx := context.Background()
	state := schemas.ProcessStepResourceModel{
		SpaceID:            types.StringValue("Spaces-1"),
		ProcessID:          types.StringValue("Processes-1"),
		Name:               types.StringValue("Step One"),
		StartTrigger:       types.StringValue("StartAfterPrevious"),
		PackageRequirement: types.StringValue("LetOctopusDecide"),
		Condition:          types.StringValue("Always"),
		StepProperties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.MaxParallelism": types.StringValue("2"),
			"Octopus.Action.TargetRoles":    types.StringValue("agent-1,agent-2"),
		}),
		ActionType:         types.StringValue("Octopus.Script"),
		Slug:               types.StringValue("step-one"),
		IsRequired:         types.BoolValue(true),
		IsDisabled:         types.BoolValue(false),
		Notes:              types.StringValue(`Some notes`),
		WorkerPoolId:       types.StringValue("WorkerPools-1"),
		WorkerPoolVariable: types.StringValue("#{Environment.WorkerPools.Default}"),
		TenantTags: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("tag-1"),
			types.StringValue("tag-2"),
		}),
		Environments: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Environments-1"),
			types.StringValue("Environments-2"),
		}),
		ExcludedEnvironments: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Environments-13"),
		}),
		Channels: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Channels-1"),
		}),
		Container: &schemas.ProcessStepActionContainerModel{
			FeedId: types.StringValue("Feeds-1"),
			Image:  types.StringValue("docker.io/library/dummy:latest"),
		},
		ActionProperties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.RunOnServer":       types.StringValue("True"),
			"Octopus.Action.Script.ScriptBody": types.StringValue("Write-Host \"Step 1, Action 1\""),
		}),
	}

	step := deployments.NewDeploymentStep("Step One")

	diags := mapProcessStepFromState(ctx, &state, step)

	assert.False(t, diags.HasError(), "Expected no errors in diagnostics")

	expectedStep := &deployments.DeploymentStep{
		Name:               "Step One",
		StartTrigger:       "StartAfterPrevious",
		PackageRequirement: "LetOctopusDecide",
		Condition:          "Always",
		Properties: map[string]core.PropertyValue{
			"Octopus.Action.MaxParallelism": core.NewPropertyValue("2", false),
			"Octopus.Action.TargetRoles":    core.NewPropertyValue("agent-1,agent-2", false),
		},
		Actions: []*deployments.DeploymentAction{
			{
				Name:                 "Step One",
				Slug:                 "step-one",
				ActionType:           "Octopus.Script",
				IsRequired:           true,
				IsDisabled:           false,
				Notes:                "Some notes",
				WorkerPool:           "WorkerPools-1",
				WorkerPoolVariable:   "#{Environment.WorkerPools.Default}",
				TenantTags:           []string{"tag-1", "tag-2"},
				Environments:         []string{"Environments-1", "Environments-2"},
				ExcludedEnvironments: []string{"Environments-13"},
				Channels:             []string{"Channels-1"},
				Container: &deployments.DeploymentActionContainer{
					FeedID: "Feeds-1",
					Image:  "docker.io/library/dummy:latest",
				},
				Properties: map[string]core.PropertyValue{
					"Octopus.Action.RunOnServer":       core.NewPropertyValue("True", false),
					"Octopus.Action.Script.ScriptBody": core.NewPropertyValue("Write-Host \"Step 1, Action 1\"", false),
				},
				Resource: *resources.NewResource(),
			},
		},
		TargetRoles: []string{},
		Resource:    *resources.NewResource(),
	}

	assert.Equal(t, expectedStep, step)
}

func TestAccMapProcessStepFromStateForScriptStep(t *testing.T) {
	ctx := context.Background()
	state := schemas.ProcessStepResourceModel{
		SpaceID:            types.StringValue("Spaces-1"),
		ProcessID:          types.StringValue("Processes-1"),
		Name:               types.StringValue("Run Script"),
		StartTrigger:       types.StringValue("StartAfterPrevious"),
		PackageRequirement: types.StringValue("LetOctopusDecide"),
		Condition:          types.StringValue("Success"),
		ActionType:         types.StringValue("Octopus.Script"),
		ActionProperties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.Script.ScriptBody": types.StringValue("Write-Host \"Minimum attributes\""),
		}),
	}

	step := deployments.NewDeploymentStep("Run Script")

	diags := mapProcessStepFromState(ctx, &state, step)

	assert.False(t, diags.HasError(), "Expected no errors in diagnostics")

	expectedStep := &deployments.DeploymentStep{
		Name:               "Run Script",
		StartTrigger:       "StartAfterPrevious",
		PackageRequirement: "LetOctopusDecide",
		Condition:          "Success",
		Properties:         map[string]core.PropertyValue{},
		Actions: []*deployments.DeploymentAction{
			{
				Name:       "Run Script",
				ActionType: "Octopus.Script",
				Properties: map[string]core.PropertyValue{
					"Octopus.Action.Script.ScriptBody": core.NewPropertyValue("Write-Host \"Minimum attributes\"", false),
				},
				Resource: *resources.NewResource(),
			},
		},
		TargetRoles: []string{},
		Resource:    *resources.NewResource(),
	}

	assert.Equal(t, expectedStep, step)
}

func TestAccMapProcessStepToStateWithAllAttributes(t *testing.T) {
	//ctx := context.Background()

	action := deployments.NewDeploymentAction("Step One", "Octopus.Script")
	action.SetID("12345678-1234-1234-1234-123456789000")
	action.Name = "Step One"
	action.Slug = "step-one"
	action.ActionType = "Octopus.Script"
	action.IsRequired = true
	action.IsDisabled = false
	action.Notes = "Some notes"
	action.WorkerPool = "WorkerPools-1"
	action.WorkerPoolVariable = "#{Environment.WorkerPools.Default}"
	action.TenantTags = []string{"tag-1", "tag-2"}
	action.Environments = []string{"Environments-1", "Environments-2"}
	action.ExcludedEnvironments = []string{"Environments-13"}
	action.Channels = []string{"Channels-1"}
	action.Container = &deployments.DeploymentActionContainer{
		FeedID: "Feeds-1",
		Image:  "docker.io/library/dummy:latest",
	}
	action.Properties = map[string]core.PropertyValue{
		"Octopus.Action.RunOnServer":       core.NewPropertyValue("True", false),
		"Octopus.Action.Script.ScriptBody": core.NewPropertyValue("Write-Host \"Step 1, Action 1\"", false),
	}

	step := deployments.NewDeploymentStep("Step One")
	step.SetID("12345678-1234-1234-1234-123456789001")
	step.StartTrigger = "StartAfterPrevious"
	step.PackageRequirement = "LetOctopusDecide"
	step.Condition = "Success"
	step.Properties = map[string]core.PropertyValue{
		"Octopus.Action.MaxParallelism": core.NewPropertyValue("2", false),
		"Octopus.Action.TargetRoles":    core.NewPropertyValue("agent-1,agent-2", false),
	}
	step.Actions = []*deployments.DeploymentAction{action}

	process := &deployments.DeploymentProcess{
		SpaceID:   "Spaces-1",
		ProjectID: "Projects-1",
		Steps:     []*deployments.DeploymentStep{step},
	}
	process.SetID("Processes-1")

	state := schemas.ProcessStepResourceModel{
		SpaceID:   types.StringValue(process.SpaceID),
		ProcessID: types.StringValue(process.ID),
	}
	diags := mapProcessStepToState(process, step, &state)

	assert.False(t, diags.HasError(), "Expected no errors in diagnostics")

	expectedState := schemas.ProcessStepResourceModel{
		SpaceID:            types.StringValue("Spaces-1"),
		ProcessID:          types.StringValue("Processes-1"),
		Name:               types.StringValue("Step One"),
		StartTrigger:       types.StringValue("StartAfterPrevious"),
		PackageRequirement: types.StringValue("LetOctopusDecide"),
		Condition:          types.StringValue("Success"),
		StepProperties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.MaxParallelism": types.StringValue("2"),
			"Octopus.Action.TargetRoles":    types.StringValue("agent-1,agent-2"),
		}),
		ActionType:         types.StringValue("Octopus.Script"),
		Slug:               types.StringValue("step-one"),
		IsRequired:         types.BoolValue(true),
		IsDisabled:         types.BoolValue(false),
		Notes:              types.StringValue(`Some notes`),
		WorkerPoolId:       types.StringValue("WorkerPools-1"),
		WorkerPoolVariable: types.StringValue("#{Environment.WorkerPools.Default}"),
		TenantTags: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("tag-1"),
			types.StringValue("tag-2"),
		}),
		Environments: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Environments-1"),
			types.StringValue("Environments-2"),
		}),
		ExcludedEnvironments: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Environments-13"),
		}),
		Channels: types.SetValueMust(types.StringType, []attr.Value{
			types.StringValue("Channels-1"),
		}),
		Container: &schemas.ProcessStepActionContainerModel{
			FeedId: types.StringValue("Feeds-1"),
			Image:  types.StringValue("docker.io/library/dummy:latest"),
		},
		ActionProperties: types.MapValueMust(types.StringType, map[string]attr.Value{
			"Octopus.Action.RunOnServer":       types.StringValue("True"),
			"Octopus.Action.Script.ScriptBody": types.StringValue("Write-Host \"Step 1, Action 1\""),
		}),
	}
	expectedState.ID = types.StringValue(step.ID)

	assert.Equal(t, expectedState, state)
}
