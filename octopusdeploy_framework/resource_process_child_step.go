package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
)

var _ resource.Resource = &processChildStepResource{}

type processChildStepResource struct {
	*Config
}

func NewProcessChildStepResource() resource.Resource {
	return &processChildStepResource{}
}

func (r *processChildStepResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProcessChildStepResourceName)
}

func (r *processChildStepResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProcessChildStepSchema{}.GetResourceSchema()
}

func (r *processChildStepResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *processChildStepResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.ProcessChildStepResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ParentID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("creating process child step: %s", data.Name.ValueString()))

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("Error creating process child step, unable to find a process", err.Error())
		return
	}

	parent, ok := findStepFromProcessByID(process, parentId)
	if !ok {
		resp.Diagnostics.AddError("Error creating process child step, unable to find a parent step", err.Error())
		return
	}

	action := deployments.NewDeploymentAction(data.Name.ValueString(), data.ActionType.ValueString())
	mapProcessChildStepActionFromState(data, action)

	parent.Actions = append(parent.Actions, action)

	updatedProcess, err := deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to create process child step", err.Error())
		return
	}

	updatedStep, ok := findStepFromProcessByID(updatedProcess, parentId)
	if !ok {
		resp.Diagnostics.AddError("unable to create process child step, unable to find a parent step '%s'", parent.ID)
		return
	}

	createdAction, ok := findActionFromProcessStepByName(updatedStep, action.Name)
	if !ok {
		resp.Diagnostics.AddError("unable to create process child step", action.Name)
		return
	}

	mapProcessChildStepActionToState(createdAction, updatedStep, updatedProcess, data)

	tflog.Info(ctx, fmt.Sprintf("process child step created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processChildStepResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.ProcessChildStepResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading process child step (%s)", data.ID))

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ParentID.ValueString()
	actionId := data.ID.ValueString()

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to find process", err.Error())
		return
	}

	parent, ok := findStepFromProcessByID(process, parentId)
	if !ok {
		resp.Diagnostics.AddError("unable to find parent step '%s'", parentId)
		return
	}

	action, ok := findActionFromProcessStepByID(parent, actionId)
	if !ok {
		resp.Diagnostics.AddError("unable to find process child step", actionId)
		return
	}

	mapProcessChildStepActionToState(action, parent, process, data)

	tflog.Info(ctx, fmt.Sprintf("process chidl step read (%s)", actionId))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processChildStepResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *schemas.ProcessChildStepResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ParentID.ValueString()
	actionId := data.ID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("updating process child step (%s)", actionId))

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to find process", err.Error())
		return
	}

	parent, ok := findStepFromProcessByID(process, parentId)
	if !ok {
		resp.Diagnostics.AddError("unable to find parent step '%s'", parentId)
		return
	}

	action, ok := findActionFromProcessStepByID(parent, actionId)
	if !ok {
		resp.Diagnostics.AddError("unable to find process child step", actionId)
		return
	}

	mapProcessChildStepActionFromState(data, action)

	updatedProcess, err := deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to update process child step", err.Error())
		return
	}

	updatedStep, ok := findStepFromProcessByID(updatedProcess, parentId)
	if !ok {
		resp.Diagnostics.AddError("unable to update process child step, unable to find a parent step '%s'", parent.ID)
		return
	}

	updatedAction, ok := findActionFromProcessStepByID(updatedStep, actionId)
	if !ok {
		resp.Diagnostics.AddError("unable to update process child step", actionId)
		return
	}

	mapProcessChildStepActionToState(updatedAction, updatedStep, updatedProcess, data)

	tflog.Info(ctx, fmt.Sprintf("process child step updated (%s)", actionId))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processChildStepResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *schemas.ProcessChildStepResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ParentID.ValueString()
	actionId := data.ID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("deleting process child step (%s)", data.ID))

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to find process", err.Error())
		return
	}

	parent, ok := findStepFromProcessByID(process, parentId)
	if !ok {
		resp.Diagnostics.AddError("unable to find parent step '%s'", parentId)
		return
	}

	var filteredActions []*deployments.DeploymentAction
	for _, action := range parent.Actions {
		if actionId != action.GetID() {
			filteredActions = append(filteredActions, action)
		}
	}
	parent.Actions = filteredActions

	_, err = deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to delete process child step", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func mapProcessChildStepActionFromState(state *schemas.ProcessChildStepResourceModel, action *deployments.DeploymentAction) {
	action.Name = state.Name.ValueString()
	action.ActionType = state.ActionType.ValueString()

	runOnServer := "False"
	if state.RunOnServer.ValueBool() {
		runOnServer = "True"
	}
	action.Properties["Octopus.Action.RunOnServer"] = core.NewPropertyValue(runOnServer, false)

	action.Properties["Octopus.Action.Script.Syntax"] = core.NewPropertyValue(state.ScriptSyntax.ValueString(), false)
	action.Properties["Octopus.Action.Script.ScriptBody"] = core.NewPropertyValue(state.ScriptBody.ValueString(), false)
	action.Properties["Octopus.Action.MaintainedBy.TerraformProvider"] = core.NewPropertyValue("True", false)
}

func mapProcessChildStepActionToState(action *deployments.DeploymentAction, step *deployments.DeploymentStep, process *deployments.DeploymentProcess, state *schemas.ProcessChildStepResourceModel) {
	state.ID = types.StringValue(action.GetID())
	state.SpaceID = types.StringValue(process.SpaceID)
	state.ProcessID = types.StringValue(process.GetID())
	state.ParentID = types.StringValue(step.GetID())
	state.Name = types.StringValue(action.Name)
	state.ActionType = types.StringValue(action.ActionType)

	value, _ := strconv.ParseBool(action.Properties["Octopus.Action.RunOnServer"].Value)
	state.RunOnServer = types.BoolValue(value)

	state.ScriptSyntax = types.StringValue(action.Properties["Octopus.Action.Script.Syntax"].Value)
	state.ScriptBody = types.StringValue(action.Properties["Octopus.Action.Script.ScriptBody"].Value)
}

func findActionFromProcessStepByID(step *deployments.DeploymentStep, actionId string) (*deployments.DeploymentAction, bool) {
	for _, action := range step.Actions {
		if action.ID == actionId {
			return action, true
		}
	}
	return nil, false
}

func findActionFromProcessStepByName(step *deployments.DeploymentStep, name string) (*deployments.DeploymentAction, bool) {
	for _, action := range step.Actions {
		if action.Name == name {
			return action, true
		}
	}
	return nil, false
}
