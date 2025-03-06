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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &processChildStepsOrderResource{}

type processChildStepsOrderResource struct {
	*Config
}

func NewProcessChildStepsOrderResource() resource.Resource {
	return &processChildStepsOrderResource{}
}

func (r *processChildStepsOrderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProcessChildStepsOrderResourceName)
}

func (r *processChildStepsOrderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProcessChildStepsOrderSchema{}.GetResourceSchema()
}

func (r *processChildStepsOrderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
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

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("Error creating process child steps order, unable to find a process", err.Error())
		return
	}

	parent, ok := findStepFromProcessByID(process, parentId)
	if !ok {
		resp.Diagnostics.AddError("Error creating process child steps order", fmt.Sprintf("unable to find a parent step with id '%s'", parentId))
		return
	}

	mapDiagnostics := mapProcessChildStepsOrderFromState(data, parent)
	resp.Diagnostics.Append(mapDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedProcess, err := deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to create process step", err.Error())
		return
	}

	updatedParent, ok := findStepFromProcessByID(updatedProcess, parentId)
	if !ok {
		resp.Diagnostics.AddError("Error creating process child steps order", fmt.Sprintf("unable to find a parent step with id '%s'", parentId))
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

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to find process", err.Error())
		return
	}

	parent, ok := findStepFromProcessByID(process, parentId)
	if !ok {
		resp.Diagnostics.AddError("Error reading process child steps order, unable to find a parent step", err.Error())
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

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to load process", err.Error())
		return
	}

	parent, ok := findStepFromProcessByID(process, parentId)
	if !ok {
		resp.Diagnostics.AddError("Error updating process child steps order, unable to find a parent step", err.Error())
		return
	}

	mapDiagnostics := mapProcessChildStepsOrderFromState(data, parent)
	resp.Diagnostics.Append(mapDiagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedProcess, err := deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to update process child steps order", err.Error())
		return
	}

	updatedParent, ok := findStepFromProcessByID(updatedProcess, parentId)
	if !ok {
		resp.Diagnostics.AddError("Error updating process child steps order, unable to find a parent step", err.Error())
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

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to load process", err.Error())
		return
	}

	_, ok := findStepFromProcessByID(process, parentId)
	if !ok {
		resp.Diagnostics.AddError("Cannot delete process child steps order", fmt.Sprintf("unable to find a parent step (%s)", parentId))
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
			diags.AddError("Ordered child step with id '%s' is not part of the parent step", id)
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
			"Some children steps were not included in the order",
			message,
		)
	}

	if diags.HasError() {
		return diags
	}

	step.Actions = reorderedActions

	return diags
}

func mapProcessChildStepsOrderToState(process *deployments.DeploymentProcess, step *deployments.DeploymentStep, state *schemas.ProcessChildStepsOrderResourceModel) {
	state.ID = types.StringValue(step.GetID())
	state.SpaceID = types.StringValue(process.SpaceID)
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
