package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &processStepsOrderResource{}

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

	mapProcessStepsOrderFromState(data, process)

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

	mapProcessStepsOrderFromState(data, process)

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

func mapProcessStepsOrderFromState(state *schemas.ProcessStepsOrderResourceModel, process *deployments.DeploymentProcess) {
	// Validate that ordering include all steps defined in the process

	lookup := make(map[string]*deployments.DeploymentStep)
	for _, step := range process.Steps {
		lookup[step.GetID()] = step
	}

	orderedIds := util.GetIds(state.Steps)

	reorderedSteps := make([]*deployments.DeploymentStep, len(process.Steps))
	for i, id := range orderedIds {
		step, _ := lookup[id]
		reorderedSteps[i] = step
	}
	process.Steps = reorderedSteps
}

func mapProcessStepsOrderToState(process *deployments.DeploymentProcess, state *schemas.ProcessStepsOrderResourceModel) {
	state.ID = types.StringValue(process.GetID())
	state.SpaceID = types.StringValue(process.SpaceID)
	state.ProcessID = types.StringValue(process.GetID())

	var steps = make([]attr.Value, len(process.Steps))
	for i, step := range process.Steps {
		steps[i] = types.StringValue(step.GetID())
	}
	state.Steps, _ = types.ListValue(types.StringType, steps)
}
