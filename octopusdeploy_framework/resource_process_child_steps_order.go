package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

var (
	_ resource.ResourceWithModifyPlan  = &processChildStepsOrderResource{}
	_ resource.ResourceWithImportState = &processChildStepsOrderResource{}
)

type processChildStepsOrderResource struct {
	*Config
}

func NewProcessChildStepsOrderResource() resource.Resource {
	return &processChildStepsOrderResource{}
}

func (r *processChildStepsOrderResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProcessChildStepsOrderResourceName)
}

func (r *processChildStepsOrderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProcessChildStepsOrderSchema{}.GetResourceSchema()
}

func (r *processChildStepsOrderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *processChildStepsOrderResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	identifiers := strings.Split(request.ID, ":")

	if len(identifiers) != 2 {
		response.Diagnostics.AddError(
			"Incorrect Import Identifier",
			fmt.Sprintf("Expected import identifier with format: ProcessId:ParentStepId (e.g. deploymentprocess-Projects-123:00000000-0000-0000-0000-000000000010). Got: %q", request.ID),
		)
		return
	}

	processId := identifiers[0]
	parentStepId := identifiers[1]

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("process_id"), processId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("parent_id"), parentStepId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), parentStepId)...)

	process, diags := loadProcessWrapperForSteps(r.Config.Client, r.Config.SpaceID, processId)
	if len(diags) > 0 {
		response.Diagnostics.Append(diags...)
		return
	}

	parent, ok := process.FindStepByID(parentStepId)
	if !ok {
		response.Diagnostics.AddError("Error importing process child steps order", fmt.Sprintf("unable to find a parent step (id: %s)", parentStepId))
		return
	}

	// Import all actions, because Read method relies on configured actions to avoid state drifting (see 'mapProcessChildStepsOrderToState')
	var actions []attr.Value
	// Exclude first action, which is embedded into the parent step
	for _, action := range parent.Actions[1:] {
		if action != nil {
			actions = append(actions, types.StringValue(action.GetID()))
		}
	}
	children, _ := types.ListValue(types.StringType, actions)

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("process_id"), processId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("parent_id"), parentStepId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), parentStepId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("children"), children)...)
}

func (r *processChildStepsOrderResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddWarning("Deleting child steps order", "Applying this resource destruction will not update child steps and their order")
		return
	}

	var state *schemas.ProcessChildStepsOrderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.SpaceID.IsUnknown() {
		return
	}

	if state.ProcessID.IsUnknown() {
		return
	}

	if state.ParentID.IsUnknown() {
		return
	}

	spaceId := state.SpaceID.ValueString()
	processId := state.ProcessID.ValueString()
	parentId := state.ParentID.ValueString()

	// Do the validation based on steps stored in Octopus Deploy
	process, diags := loadProcessWrapperForSteps(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	parent, ok := process.FindStepByID(parentId)
	if !ok {
		resp.Diagnostics.AddError("Error modifying plan for child steps order", fmt.Sprintf("unable to find a parent step with id '%s'", parentId))
		return
	}

	configuredActions := util.GetIds(state.Children)

	// Validate that all actions are included in the order resource
	configuredActionsLookup := make(map[string]bool, len(configuredActions))
	for _, id := range configuredActions {
		configuredActionsLookup[id] = true
	}

	var missingActions []string
	for _, action := range parent.Actions[1:] {
		if !configuredActionsLookup[action.GetID()] {
			missingActions = append(missingActions, fmt.Sprintf("'%s' (%s)", action.Name, action.GetID()))
		}
	}

	if len(missingActions) > 0 {
		message := fmt.Sprintf("The following child steps were not included in the steps order and will be added at the end.\nNote that their order at the end is not guaranteed:\n%v", missingActions)
		resp.Diagnostics.AddWarning(
			"Some process child steps were not included in the order",
			message,
		)
	}

	// Validate that included actions belong to the parent step
	existingActions := make(map[string]*deployments.DeploymentAction)
	for _, action := range parent.Actions {
		existingActions[action.GetID()] = action
	}

	var unknownActions []string
	for _, id := range configuredActions {
		_, found := existingActions[id]
		if !found {
			unknownActions = append(unknownActions, id)
		}
	}

	if len(unknownActions) > 0 {
		message := fmt.Sprintf("Following steps are not part of the process: %v", unknownActions)
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Some ordered child steps do not belong to the parent step '%s'", parent.ID),
			message,
		)
	}
}

func (r *processChildStepsOrderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.ProcessChildStepsOrderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ParentID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("creating process child steps order for parent %s", parentId))

	process, diags := loadProcessWrapperForSteps(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	parent, parentFound := process.FindStepByID(parentId)
	if !parentFound {
		resp.Diagnostics.AddError("Error creating process child steps order", fmt.Sprintf("Unable to find a parent step with id '%s'", parentId))
		return
	}

	mapDiagnostics := mapProcessChildStepsOrderFromState(data, parent)
	resp.Diagnostics.Append(mapDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedProcess, err := process.Update(r.Config.Client)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create process step", err.Error())
		return
	}

	updatedParent, updatedParentFound := updatedProcess.FindStepByID(parentId)
	if !updatedParentFound {
		resp.Diagnostics.AddError("Error creating process child steps order", fmt.Sprintf("Unable to find a parent step with id '%s'", parentId))
		return
	}

	mapProcessChildStepsOrderToState(updatedProcess, updatedParent, data)

	tflog.Info(ctx, fmt.Sprintf("process child steps order created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processChildStepsOrderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.ProcessChildStepsOrderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("reading process child steps order (%s)", parentId))

	process, diags := loadProcessWrapperForSteps(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	parent, ok := process.FindStepByID(parentId)
	if !ok {
		resp.Diagnostics.AddError("Error reading process child steps order", fmt.Sprintf("Unable to find a parent step (id: %s)", parentId))
		return
	}

	mapProcessChildStepsOrderToState(process, parent, data)

	tflog.Info(ctx, fmt.Sprintf("process child steps order read (%s)", parentId))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processChildStepsOrderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *schemas.ProcessChildStepsOrderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("updating process child steps order (%s)", parentId))

	process, diags := loadProcessWrapperForSteps(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	parent, parentFound := process.FindStepByID(parentId)
	if !parentFound {
		resp.Diagnostics.AddError("Error updating process child steps order", fmt.Sprintf("unable to find a parent step (id: %s)", parentId))
		return
	}

	mapDiagnostics := mapProcessChildStepsOrderFromState(data, parent)
	resp.Diagnostics.Append(mapDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedProcess, err := process.Update(r.Config.Client)
	if err != nil {
		resp.Diagnostics.AddError("unable to update process child steps order", err.Error())
		return
	}

	updatedParent, updatedParentFound := updatedProcess.FindStepByID(parentId)
	if !updatedParentFound {
		resp.Diagnostics.AddError("Error updating process child steps order", fmt.Sprintf("unable to find a parent step (id: %s)", parentId))
		return
	}

	mapProcessChildStepsOrderToState(updatedProcess, updatedParent, data)

	tflog.Info(ctx, fmt.Sprintf("process steps order updated (%s)", processId))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processChildStepsOrderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *schemas.ProcessChildStepsOrderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	parentId := data.ID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("deleting process steps order (%s)", processId))

	process, diags := loadProcessWrapperForSteps(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	_, ok := process.FindStepByID(parentId)
	if !ok {
		resp.Diagnostics.AddError("Cannot delete process child steps order", fmt.Sprintf("Unable to find a parent step (%s)", parentId))
		return
	}

	// Do nothing

	resp.State.RemoveResource(ctx)
}

func mapProcessChildStepsOrderFromState(state *schemas.ProcessChildStepsOrderResourceModel, step *deployments.DeploymentStep) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if len(step.Actions) == 0 {
		diags.AddError("Cannot map child steps ordering", "Parent step is missing embedded execution action")
		return diags
	}

	embeddedAction := step.Actions[0]

	lookup := make(map[string]*deployments.DeploymentAction)
	for _, action := range step.Actions[1:] {
		lookup[action.GetID()] = action
	}

	orderedIds := util.GetIds(state.Children)

	reorderedActions := []*deployments.DeploymentAction{embeddedAction}
	for _, id := range orderedIds {
		action, found := lookup[id]
		if !found {
			diags.AddError("Error mapping child steps order", fmt.Sprintf("Child step (id: %s) does not belong to the parent step", id))
			continue
		}

		delete(lookup, id)
		reorderedActions = append(reorderedActions, action)
	}

	// Append unordered actions
	var missingActions []string
	for _, action := range lookup {
		missingActions = append(missingActions, fmt.Sprintf("'%s' (%s)", action.Name, action.GetID()))
		reorderedActions = append(reorderedActions, action)
	}

	if len(missingActions) > 0 {
		message := fmt.Sprintf("The following child steps were not included in the steps order and will be added at the end.\nNote that their order at the end is not guaranteed:\n%v", missingActions)
		diags.AddWarning(
			"Some child steps were not included in the order",
			message,
		)
	}

	if diags.HasError() {
		return diags
	}

	step.Actions = reorderedActions

	return diags
}

func mapProcessChildStepsOrderToState(process processWrapper, step *deployments.DeploymentStep, state *schemas.ProcessChildStepsOrderResourceModel) {
	state.ID = types.StringValue(step.GetID())
	state.SpaceID = types.StringValue(process.GetSpaceID())
	state.ProcessID = types.StringValue(process.GetID())
	state.ParentID = types.StringValue(step.GetID())

	// parent step's first action is "embedded" into a step resource - we want to exclude it from ordering and keep it always first
	childActions := step.Actions[1:]

	configuredActions := min(len(state.Children.Elements()), len(childActions))
	var actions []attr.Value
	// Take only "configured" amount of steps to avoid state drifting when practitioner didn't include all steps into the order resource
	for _, action := range childActions[:configuredActions] {
		if action != nil {
			actions = append(actions, types.StringValue(action.GetID()))
		}
	}
	state.Children, _ = types.ListValue(types.StringType, actions)
}
