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
)

var (
	_ resource.ResourceWithModifyPlan  = &processStepsOrderResource{}
	_ resource.ResourceWithImportState = &processStepsOrderResource{}
)

type processStepsOrderResource struct {
	*Config
}

func NewProcessStepsOrderResource() resource.Resource {
	return &processStepsOrderResource{}
}

func (r *processStepsOrderResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProcessStepsOrderResourceName)
}

func (r *processStepsOrderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProcessStepsOrderSchema{}.GetResourceSchema()
}

func (r *processStepsOrderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *processStepsOrderResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	processId := request.ID

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, r.Config.SpaceID, processId)
	if len(diags) > 0 {
		response.Diagnostics.Append(diags...)
		return
	}

	// Import all steps, because Read method relies on configured steps to avoid state drifting (see 'mapProcessStepsOrderToState')
	var identifiers []attr.Value
	for _, step := range process.GetSteps() {
		if step != nil {
			identifiers = append(identifiers, types.StringValue(step.GetID()))
		}
	}

	importedSteps, stepDiagnostics := types.ListValue(types.StringType, identifiers)
	response.Diagnostics.Append(stepDiagnostics...)

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), processId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("process_id"), processId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("steps"), importedSteps)...)
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
	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	orderedIds := util.GetIds(state.Steps)

	// Validate that all steps are included in the order resource
	includedStepsLookup := make(map[string]bool, len(orderedIds))
	for _, id := range orderedIds {
		includedStepsLookup[id] = true
	}

	var missingSteps []string
	for _, step := range process.GetSteps() {
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
	for _, step := range process.GetSteps() {
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
			fmt.Sprintf("Some ordered steps are not part of the process '%s'", process.GetID()),
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

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	mappingDiags := mapProcessStepsOrderFromState(data, process)
	resp.Diagnostics.Append(mappingDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedProcess, err := process.Update(r.Config.Client)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create process step", err.Error())
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

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
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

	process, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	mappingDiags := mapProcessStepsOrderFromState(data, process)
	resp.Diagnostics.Append(mappingDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatedProcess, err := process.Update(r.Config.Client)
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

	_, diags := loadProcessWrapperByProcessId(r.Config.Client, spaceId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Do nothing

	resp.State.RemoveResource(ctx)
}

func mapProcessStepsOrderFromState(state *schemas.ProcessStepsOrderResourceModel, process processWrapper) diag.Diagnostics {
	diags := diag.Diagnostics{}

	lookup := make(map[string]*deployments.DeploymentStep)
	for _, step := range process.GetSteps() {
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

	process.ReplaceSteps(reorderedSteps)

	return diags
}

func mapProcessStepsOrderToState(process processWrapper, state *schemas.ProcessStepsOrderResourceModel) {
	state.ID = types.StringValue(process.GetID())
	state.SpaceID = types.StringValue(process.GetSpaceID())
	state.ProcessID = types.StringValue(process.GetID())

	processSteps := process.GetSteps()
	configuredSteps := min(len(state.Steps.Elements()), len(processSteps))

	var steps []attr.Value
	// Take only "configured" amount of steps to avoid state drifting when practitioner didn't include all steps into the order resource
	for _, step := range processSteps[:configuredSteps] {
		if step != nil {
			steps = append(steps, types.StringValue(step.GetID()))
		}
	}
	state.Steps, _ = types.ListValue(types.StringType, steps)
}
