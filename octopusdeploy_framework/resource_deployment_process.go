package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/mappers"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type deploymentProcessResource struct {
	*Config
}

// implementation checks
var _ resource.Resource = &deploymentProcessResource{}
var _ resource.ResourceWithImportState = &deploymentProcessResource{}

func NewDeploymentProcessResource() resource.Resource {
	return &deploymentProcessResource{}
}

func (d *deploymentProcessResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("deployment_process")
}

func (d *deploymentProcessResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetDeploymentProcessResourceSchema()
}

func (d *deploymentProcessResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	d.Config = ResourceConfiguration(req, resp)
}

func (d *deploymentProcessResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.DeploymentProcessResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating deployment process: %#v", plan))
	spaceID := plan.SpaceID.ValueString()
	project, err := projects.GetByID(d.Client, plan.SpaceID.ValueString(), plan.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error getting project %s", plan.ProjectID.ValueString()), err.Error())
		return
	}

	var current *deployments.DeploymentProcess
	if project.PersistenceSettings != nil && project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		current, err = deployments.GetDeploymentProcessByGitRef(d.Client, spaceID, project, plan.Branch.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by git ref", err.Error())
			return
		}
	} else {
		current, err = deployments.GetDeploymentProcessByID(d.Client, spaceID, project.DeploymentProcessID)
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by ID", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(mappers.MapStateToDeploymentProcess(ctx, &plan, current)...)

	if resp.Diagnostics.HasError() {
		return
	}

	current, err = deployments.UpdateDeploymentProcess(d.Client, current)
	if err != nil {
		resp.Diagnostics.AddError("unable to update deployment process", err.Error())
		return
	}

	plan.ID = types.StringValue(current.ID)
	plan.Branch = types.StringValue(current.Branch)
	plan.ProjectID = types.StringValue(current.ProjectID)
	plan.SpaceID = types.StringValue(current.SpaceID)
	plan.Version = types.StringValue(fmt.Sprintf("%d", current.Version))
	plan.LastSnapshotID = types.StringValue(current.LastSnapshotID)

	for _, step := range plan.Steps.Elements() {
		stepAttrs := step.(types.Object).Attributes()
		name := stepAttrs["name"].(types.String).ValueString()
		var currentStep *deployments.DeploymentStep
		for _, processStep := range current.Steps {
			if processStep.Name == name {
				currentStep = processStep
				break
			}
		}

		stepAttrs["id"] = types.StringValue(currentStep.ID)

		for actionKey, _ := range mappers.ActionMappers {
			for _, action := range stepAttrs[actionKey].(types.List).Elements() {
				actionAttrs := action.(types.Object).Attributes()
				actionName := actionAttrs["name"].(types.String).ValueString()
				var currentAction *deployments.DeploymentAction
				for _, stepAction := range currentStep.Actions {
					if stepAction.Name == actionName {
						currentAction = stepAction
						break
					}
				}

				actionAttrs["id"] = types.StringValue(currentAction.ID)
				actionAttrs["channels"] = types.ListValueMust(types.StringType, []attr.Value{})
			}
		}

		//return mapStepsToState(ctx, state, deploymentProcess)
	}

	//resp.Diagnostics.Append(mappers.MapDeploymentProcessToState(ctx, current, &plan)...)
	//if resp.Diagnostics.HasError() {
	//	return
	//}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

//func (r *deploymentProcessResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
//	var plan *schemas.DeploymentProcessResourceModel
//
//	if req.Plan.Raw.IsNull() {
//		return
//	}
//
//	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//
//	if req.State.Raw.IsNull() {
//		return
//	}
//
//	var state *schemas.DeploymentProcessResourceModel
//	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
//	if resp.Diagnostics.HasError() {
//		return
//	}
//
//	for _, stateStep := range state.Steps.Elements() {
//		stateStepAttrs := stateStep.(types.Object).Attributes()
//		for _, planStep := range plan.Steps.Elements() {
//			planStepAttrs := planStep.(types.Object).Attributes()
//			if stateStepAttrs["name"].(types.String).ValueString() == planStepAttrs["name"].(types.String).ValueString() {
//				for actionKey, _ := range schemas.ActionsAttributeToActionTypeMap {
//					for _, stateAction := range stateStepAttrs[actionKey].(types.List).Elements() {
//						stateActionAttrs := stateAction.(types.Object).Attributes()
//						for _, planAction := range planStepAttrs[actionKey].(types.List).Elements() {
//							planActionAttrs := planAction.(types.Object).Attributes()
//							if stateActionAttrs["name"].(types.String).ValueString() == planActionAttrs["name"].(types.String).ValueString() {
//								tflog.Debug(ctx, stateActionAttrs["name"].(types.String).ValueString())
//								//resp.Plan.SetAttribute(ctx, path.Root("step[0].action[0].primary_package[0].properties"), req.State.GetAttribute(ctx, path.Root("step[0].action[0].primary_package[0].properties")).(interface{}))
//							}
//						}
//					}
//				}
//			}
//		}
//
//		//stepAttrs["id"] = types.StringValue(currentStep.ID)
//		//
//		//for actionKey, _ := range schemas.ActionsAttributeToActionTypeMap {
//		//	for _, action := range stepAttrs[actionKey].(types.List).Elements() {
//		//		actionAttrs := action.(types.Object).Attributes()
//		//		actionName := actionAttrs["name"].(types.String).ValueString()
//		//		var currentAction *deployments.DeploymentAction
//		//		for _, stepAction := range currentStep.Actions {
//		//			if stepAction.Name == actionName {
//		//				currentAction = stepAction
//		//				break
//		//			}
//		//		}
//		//
//		//		actionAttrs["id"] = types.StringValue(currentAction.ID)
//		//
//		//	}
//		//}
//		//
//		//steps = append(steps, types.ObjectValueMust(getStepTypeAttrs(), stepAttrs))
//		//tflog.Debug(ctx, name)
//	}
//
//	//templates := plan.Template
//	//expandedActionTemplates := schemas.ExpandActionTemplateParameters(templates)
//	//templateIdsValues := schemas.FlattenTemplateIds(expandedActionTemplates)
//	//resp.Plan.SetAttribute(ctx, path.Root("step[0].action[0].primary_package[0].properties"), templateIdsValues)
//}

func (d *deploymentProcessResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.DeploymentProcessResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading deployment process: %#v", state))
	spaceID := state.SpaceID.ValueString()
	project, err := projects.GetByID(d.Client, state.SpaceID.ValueString(), state.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error getting project %s", state.ProjectID.ValueString()), err.Error())
		return
	}

	var current *deployments.DeploymentProcess
	if project.PersistenceSettings != nil && project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		current, err = deployments.GetDeploymentProcessByGitRef(d.Client, spaceID, project, state.Branch.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by git ref", err.Error())
			return
		}
	} else {
		current, err = deployments.GetDeploymentProcessByID(d.Client, spaceID, project.DeploymentProcessID)
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by ID", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(mappers.MapDeploymentProcessToState(ctx, current, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *deploymentProcessResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.DeploymentProcessResourceModel
	var state schemas.DeploymentProcessResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading deployment process: %#v", plan))
	spaceID := plan.SpaceID.ValueString()
	project, err := projects.GetByID(d.Client, plan.SpaceID.ValueString(), plan.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error getting project %s", plan.ProjectID.ValueString()), err.Error())
		return
	}

	var current *deployments.DeploymentProcess
	if project.PersistenceSettings != nil && project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		current, err = deployments.GetDeploymentProcessByGitRef(d.Client, spaceID, project, plan.Branch.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by git ref", err.Error())
			return
		}
	} else {
		current, err = deployments.GetDeploymentProcessByID(d.Client, spaceID, project.DeploymentProcessID)
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by ID", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(mappers.MapStateToDeploymentProcess(ctx, &plan, current)...)

	if resp.Diagnostics.HasError() {
		return
	}

	current, err = deployments.UpdateDeploymentProcess(d.Client, current)
	if err != nil {
		resp.Diagnostics.AddError("unable to update deployment process", err.Error())
		return
	}

	resp.Diagnostics.Append(mappers.MapDeploymentProcessToState(ctx, current, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (d *deploymentProcessResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.DeploymentProcessResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading deployment process: %#v", state))
	spaceID := state.SpaceID.ValueString()
	project, err := projects.GetByID(d.Client, state.SpaceID.ValueString(), state.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error getting project %s", state.ProjectID.ValueString()), err.Error())
		return
	}

	var current *deployments.DeploymentProcess
	if project.PersistenceSettings != nil && project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		current, err = deployments.GetDeploymentProcessByGitRef(d.Client, spaceID, project, state.Branch.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by git ref", err.Error())
			return
		}
	} else {
		current, err = deployments.GetDeploymentProcessByID(d.Client, spaceID, project.DeploymentProcessID)
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by ID", err.Error())
			return
		}
	}

	current.Steps = []*deployments.DeploymentStep{}
	deployments.UpdateDeploymentProcess(d.Client, current)
	resp.Diagnostics.Append(mappers.MapDeploymentProcessToState(ctx, current, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "deployment process deleted")
}

func (d *deploymentProcessResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//TODO implement me
	//panic("implement me")
}
