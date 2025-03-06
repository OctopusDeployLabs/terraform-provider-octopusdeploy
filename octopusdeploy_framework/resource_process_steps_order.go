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

var _ resource.ResourceWithModifyPlan = &processStepsOrderResource{}

type processStepsOrderResource struct {
	*Config
}

func NewProcessStepsOrderResource() resource.Resource {
	return &processStepsOrderResource{}
}

func (r *processStepsOrderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProcessStepsOrderResourceName)
}

func (r *processStepsOrderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProcessStepsOrderSchema{}.GetResourceSchema()
}

func (r *processStepsOrderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *processStepsOrderResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddWarning("Deleting steps order", "Applying this resource destruction will not update process steps and their order")
		return
	}

	var state *schemas.ProcessStepsOrderResourceModel
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

	spaceId := state.SpaceID.ValueString()
	processId := state.ProcessID.ValueString()

	// Do the validation based on steps stored in Octopus Deploy
	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to find process", err.Error())
		return
	}

	orderedIds := util.GetIds(state.Steps)

	// Validate that all steps are included in the order resource
	includedStepsLookup := make(map[string]bool, len(orderedIds))
	for _, id := range orderedIds {
		includedStepsLookup[id] = true
	}

	var missingSteps []string
	for _, step := range process.Steps {
		if !includedStepsLookup[step.GetID()] {
			missingSteps = append(missingSteps, fmt.Sprintf("'%s' (%s)", step.Name, step.GetID()))
		}
	}

	if len(missingSteps) > 0 {
		message := fmt.Sprintf("The following steps were not included in the steps order and will be added at the end.\nNote that their order at the end is not guaranteed:\n%v", missingSteps)
		resp.Diagnostics.AddWarning(
			"Some process steps were not included in the order",
			message,
		)
	}

	// Validate that included steps are part of the process
	lookup := make(map[string]*deployments.DeploymentStep)
	for _, step := range process.Steps {
		lookup[step.GetID()] = step
	}

	var unknownSteps []string
	for _, id := range orderedIds {
		_, found := lookup[id]
		if !found {
			unknownSteps = append(unknownSteps, id)
		}
	}

	if len(unknownSteps) > 0 {
		message := fmt.Sprintf("Following steps are not part of the process: %v", unknownSteps)
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Some ordered steps are not part of the process '%s'", process.ID),
			message,
		)
	}
}

func (r *processStepsOrderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.ProcessStepsOrderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("creating process steps order: %s", processId))

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("Error creating process steps order, unable to find a process", err.Error())
		return
	}

	diags := mapProcessStepsOrderFromState(data, process)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedProcess, err := deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to create process step", err.Error())
		return
	}

	mapProcessStepsOrderToState(updatedProcess, data)

	tflog.Info(ctx, fmt.Sprintf("process steps order created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processStepsOrderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.ProcessStepsOrderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading process steps order (%s)", data.ID))

	client := r.Config.Client
	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to find process", err.Error())
		return
	}

	mapProcessStepsOrderToState(process, data)

	tflog.Info(ctx, fmt.Sprintf("process steps order read (%s)", processId))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processStepsOrderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *schemas.ProcessStepsOrderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("updating process steps order (%s)", data.ProcessID))

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to load process", err.Error())
		return
	}

	diags := mapProcessStepsOrderFromState(data, process)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedProcess, err := deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to update process steps order", err.Error())
		return
	}

	mapProcessStepsOrderToState(updatedProcess, data)

	tflog.Info(ctx, fmt.Sprintf("process steps order updated (%s)", processId))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processStepsOrderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *schemas.ProcessStepsOrderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("deleting process steps order (%s)", processId))

	client := r.Config.Client
	_, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to load process", err.Error())
		return
	}

	// Do nothing or delete all steps

	resp.State.RemoveResource(ctx)
}

func mapProcessStepsOrderFromState(state *schemas.ProcessStepsOrderResourceModel, process *deployments.DeploymentProcess) diag.Diagnostics {
	diags := diag.Diagnostics{}

	lookup := make(map[string]*deployments.DeploymentStep)
	for _, step := range process.Steps {
		lookup[step.GetID()] = step
	}

	orderedIds := util.GetIds(state.Steps)

	var reorderedSteps []*deployments.DeploymentStep
	for _, id := range orderedIds {
		step, found := lookup[id]
		if !found {
			diags.AddError("Ordered step with id '%s' is not part of the process", id)
			continue
		}

		delete(lookup, id)
		reorderedSteps = append(reorderedSteps, step)
	}

	// Append unordered steps to the end
	var missingSteps []string
	for _, step := range lookup {
		missingSteps = append(missingSteps, fmt.Sprintf("'%s' (%s)", step.Name, step.GetID()))
		reorderedSteps = append(reorderedSteps, step)
	}

	if len(missingSteps) > 0 {
		message := fmt.Sprintf("The following steps were not included in the steps order and will be added at the end.\nNote that their order at the end is not guaranteed:\n%v", missingSteps)
		diags.AddWarning(
			"Some process steps were not included in the order",
			message,
		)
	}

	if diags.HasError() {
		return diags
	}

	process.Steps = reorderedSteps

	return diags
}

func mapProcessStepsOrderToState(process *deployments.DeploymentProcess, state *schemas.ProcessStepsOrderResourceModel) {
	state.ID = types.StringValue(process.GetID())
	state.SpaceID = types.StringValue(process.SpaceID)
	state.ProcessID = types.StringValue(process.GetID())

	configuredSteps := min(len(state.Steps.Elements()), len(process.Steps))
	var steps []attr.Value
	// Take only "configured" amount of steps to avoid state drifting when practitioner didn't include all steps into the order resource
	for _, step := range process.Steps[:configuredSteps] {
		if step != nil {
			steps = append(steps, types.StringValue(step.GetID()))
		}
	}
	state.Steps, _ = types.ListValue(types.StringType, steps)
}
