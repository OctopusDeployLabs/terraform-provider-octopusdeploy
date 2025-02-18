package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &processStepResource{}

type processStepResource struct {
	*Config
}

func NewProcessStepResource() resource.Resource {
	return &processStepResource{}
}

func (r *processStepResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProcessStepResourceName)
}

func (r *processStepResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProcessStepSchema{}.GetResourceSchema()
}

func (r *processStepResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *processStepResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.ProcessStepResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("creating process step: %s", data.Name.ValueString()))

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("Error creating process step, unable to find a process", err.Error())
		return
	}

	step := deployments.NewDeploymentStep(data.Name.ValueString())

	fromStateDiagnostics := mapProcessStepFromState(ctx, data, step)
	resp.Diagnostics.Append(fromStateDiagnostics...)
	if fromStateDiagnostics.HasError() {
		return
	}

	process.Steps = append(process.Steps, step)

	updatedProcess, err := deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to create process step", err.Error())
		return
	}

	createdStep, exists := findStepFromProcessByName(updatedProcess, step.Name)
	if !exists {
		resp.Diagnostics.AddError("unable to create process step '%s'", step.Name)
		return
	}

	toStateDiagnostics := mapProcessStepToState(updatedProcess, createdStep, data)
	resp.Diagnostics.Append(toStateDiagnostics...)
	if toStateDiagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("process step created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processStepResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.ProcessStepResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	stepId := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("reading process step (%s)", data.ID))

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to find process", err.Error())
		return
	}

	step, exists := findStepFromProcessByID(process, stepId)
	if !exists {
		tflog.Info(ctx, fmt.Sprintf("process step read (id: %s), but not found, removing ...", stepId))
		resp.State.RemoveResource(ctx)
		return
	}

	mapProcessStepToState(process, step, data)

	tflog.Info(ctx, fmt.Sprintf("process step read (%s)", step.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processStepResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *schemas.ProcessStepResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	stepId := data.ID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("updating process step (%s)", stepId))

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to load process", err.Error())
		return
	}

	step, exists := findStepFromProcessByID(process, stepId)
	if !exists {
		resp.Diagnostics.AddError("unable to find process step '%s'", stepId)
		return
	}

	diagnostics := mapProcessStepFromState(ctx, data, step)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	updatedProcess, err := deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to update process step", err.Error())
		return
	}

	updatedStep, exists := findStepFromProcessByID(updatedProcess, stepId)
	if !exists {
		resp.Diagnostics.AddError("unable to find updated process step '%s'", stepId)
		return
	}

	mapProcessStepToState(updatedProcess, updatedStep, data)

	tflog.Info(ctx, fmt.Sprintf("process step updated (%s)", updatedStep.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processStepResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *schemas.ProcessStepResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	stepId := data.ID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("deleting process step (%s)", stepId))

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to load process", err.Error())
		return
	}

	var filteredSteps []*deployments.DeploymentStep
	for _, step := range process.Steps {
		if stepId != step.GetID() {
			filteredSteps = append(filteredSteps, step)
		}
	}
	process.Steps = filteredSteps

	_, err = deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to delete process step", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func mapProcessStepFromState(ctx context.Context, state *schemas.ProcessStepResourceModel, step *deployments.DeploymentStep) diag.Diagnostics {
	step.StartTrigger = deployments.DeploymentStepStartTrigger(state.StartTrigger.ValueString())
	step.PackageRequirement = deployments.DeploymentStepPackageRequirement(state.PackageRequirement.ValueString())
	step.Condition = deployments.DeploymentStepConditionType(state.Condition.ValueString())

	if state.StepProperties.IsNull() {
		step.Properties = make(map[string]core.PropertyValue)
	} else {
		stateProperties := make(map[string]types.String, len(state.StepProperties.Elements()))
		diags := state.StepProperties.ElementsAs(ctx, &stateProperties, false)
		if diags.HasError() {
			return diags
		}

		properties := make(map[string]core.PropertyValue, len(stateProperties))
		for key, value := range stateProperties {
			if value.IsNull() {
				properties[key] = core.NewPropertyValue("", false)
			} else {
				properties[key] = core.NewPropertyValue(value.ValueString(), false)
			}
		}

		step.Properties = properties
	}

	return mapProcessStepEmbeddedActionFromState(ctx, state, step)
}

func mapProcessStepEmbeddedActionFromState(ctx context.Context, state *schemas.ProcessStepResourceModel, step *deployments.DeploymentStep) diag.Diagnostics {
	actionType := state.ActionType.ValueString()
	name := state.Name.ValueString()

	if step.Actions == nil || len(step.Actions) == 0 {
		newAction := deployments.NewDeploymentAction(name, actionType)
		step.Actions = []*deployments.DeploymentAction{newAction}
	}

	if step.Actions[0] == nil {
		step.Actions[0] = deployments.NewDeploymentAction(name, actionType)
	}

	return mapProcessStepActionFromState(ctx, state, step.Actions[0])
}

func mapProcessStepActionFromState(ctx context.Context, state *schemas.ProcessStepResourceModel, action *deployments.DeploymentAction) diag.Diagnostics {
	action.Name = state.Name.ValueString()
	action.Slug = state.Slug.ValueString() // update only embedded action slug(step slug remains original), same as UI behaviour
	action.ActionType = state.ActionType.ValueString()
	// action.Condition is not updated, replicates UI behaviour where condition of the first action of step always remains as default value (Success)

	action.IsRequired = state.IsRequired.ValueBool()
	action.IsDisabled = state.IsDisabled.ValueBool()
	action.Notes = state.Notes.ValueString()
	action.WorkerPool = state.WorkerPoolId.ValueString()
	action.WorkerPoolVariable = state.WorkerPoolVariable.ValueString()
	if state.Container == nil {
		action.Container = nil
	} else {
		action.Container = deployments.NewDeploymentActionContainer(state.Container.FeedId.ValueStringPointer(), state.Container.Image.ValueStringPointer())
	}

	diags := diag.Diagnostics{}
	if state.TenantTags.IsNull() {
		action.TenantTags = nil
	} else {
		action.TenantTags, diags = util.SetToStringArray(ctx, state.TenantTags)
		if diags.HasError() {
			return diags
		}
	}

	if state.Environments.IsNull() {
		action.Environments = nil
	} else {
		action.Environments, diags = util.SetToStringArray(ctx, state.Environments)
		if diags.HasError() {
			return diags
		}
	}

	if state.ExcludedEnvironments.IsNull() {
		action.ExcludedEnvironments = nil
	} else {
		action.ExcludedEnvironments, diags = util.SetToStringArray(ctx, state.ExcludedEnvironments)
		if diags.HasError() {
			return diags
		}
	}

	if state.Channels.IsNull() {
		action.Channels = nil
	} else {
		action.Channels, diags = util.SetToStringArray(ctx, state.Channels)
		if diags.HasError() {
			return diags
		}
	}

	if state.ActionProperties.IsNull() {
		action.Properties = nil
	} else {
		stateProperties := make(map[string]types.String, len(state.ActionProperties.Elements()))
		propertiesDiags := state.ActionProperties.ElementsAs(ctx, &stateProperties, false)
		if propertiesDiags.HasError() {
			return propertiesDiags
		}

		properties := make(map[string]core.PropertyValue, len(stateProperties))
		for key, value := range stateProperties {
			if value.IsNull() {
				properties[key] = core.NewPropertyValue("", false)
			} else {
				properties[key] = core.NewPropertyValue(value.ValueString(), false)
			}
		}

		action.Properties = properties
	}

	return diag.Diagnostics{}
}

func mapProcessStepToState(process *deployments.DeploymentProcess, step *deployments.DeploymentStep, state *schemas.ProcessStepResourceModel) diag.Diagnostics {
	state.ID = types.StringValue(step.GetID())
	state.SpaceID = types.StringValue(process.SpaceID)
	state.ProcessID = types.StringValue(process.GetID())
	state.Name = types.StringValue(step.Name)
	state.StartTrigger = types.StringValue(string(step.StartTrigger))
	state.PackageRequirement = types.StringValue(string(step.PackageRequirement))
	state.Condition = types.StringValue(string(step.Condition))

	stepProperties := make(map[string]attr.Value, len(step.Properties))
	for key, value := range step.Properties {
		stepProperties[key] = types.StringValue(value.Value)
	}

	stateProperties, diags := types.MapValue(types.StringType, stepProperties)
	if diags.HasError() {
		return diags
	}

	state.StepProperties = stateProperties

	if len(step.Actions) > 0 && step.Actions[0] != nil {
		return mapProcessStepActionToState(step.Actions[0], state)
	}

	return diag.Diagnostics{}
}

func mapProcessStepActionToState(action *deployments.DeploymentAction, state *schemas.ProcessStepResourceModel) diag.Diagnostics {
	state.ActionType = types.StringValue(action.ActionType)
	state.Slug = types.StringValue(action.Slug)
	state.IsRequired = types.BoolValue(action.IsRequired)
	state.IsDisabled = types.BoolValue(action.IsDisabled)
	state.Notes = types.StringValue(action.Notes)
	state.WorkerPoolId = types.StringValue(action.WorkerPool)
	state.WorkerPoolVariable = types.StringValue(action.WorkerPoolVariable)

	if action.Container == nil {
		state.Container = nil
	} else {
		state.Container = &schemas.ProcessStepActionContainerModel{
			FeedId: types.StringValue(action.Container.FeedID),
			Image:  types.StringValue(action.Container.Image),
		}
	}

	if action.TenantTags == nil {
		state.TenantTags = types.SetValueMust(types.StringType, []attr.Value{})
	} else {
		state.TenantTags = types.SetValueMust(types.StringType, util.ToValueSlice(action.TenantTags))
	}

	if action.Environments == nil {
		state.Environments = types.SetValueMust(types.StringType, []attr.Value{})
	} else {
		state.Environments = types.SetValueMust(types.StringType, util.ToValueSlice(action.Environments))
	}

	if action.ExcludedEnvironments == nil {
		state.ExcludedEnvironments = types.SetValueMust(types.StringType, []attr.Value{})
	} else {
		state.ExcludedEnvironments = types.SetValueMust(types.StringType, util.ToValueSlice(action.ExcludedEnvironments))
	}

	if action.Channels == nil {
		state.Channels = types.SetValueMust(types.StringType, []attr.Value{})
	} else {
		state.Channels = types.SetValueMust(types.StringType, util.ToValueSlice(action.Channels))
	}

	actionProperties := make(map[string]attr.Value, len(action.Properties))
	for key, value := range action.Properties {
		actionProperties[key] = types.StringValue(value.Value)
	}

	stateProperties, diags := types.MapValue(types.StringType, actionProperties)
	if diags.HasError() {
		return diags
	}

	state.ActionProperties = stateProperties

	return diag.Diagnostics{}
}

func findStepFromProcessByID(process *deployments.DeploymentProcess, stepID string) (*deployments.DeploymentStep, bool) {
	for _, step := range process.Steps {
		if step.ID == stepID {
			return step, true
		}
	}
	return nil, false
}

func findStepFromProcessByName(process *deployments.DeploymentProcess, name string) (*deployments.DeploymentStep, bool) {
	for _, step := range process.Steps {
		if step.Name == name {
			return step, true
		}
	}
	return nil, false
}
