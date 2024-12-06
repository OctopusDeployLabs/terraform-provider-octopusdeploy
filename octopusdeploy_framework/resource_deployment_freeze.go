package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deploymentfreezes"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

const deploymentFreezeResourceName = "deployment_freeze"

type deploymentFreezeModel struct {
	Name  types.String      `tfsdk:"name"`
	Start timetypes.RFC3339 `tfsdk:"start"`
	End   timetypes.RFC3339 `tfsdk:"end"`

	schemas.ResourceModel
}

type deploymentFreezeResource struct {
	*Config
}

var _ resource.Resource = &deploymentFreezeResource{}

func NewDeploymentFreezeResource() resource.Resource {
	return &deploymentFreezeResource{}
}

func (f *deploymentFreezeResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(deploymentFreezeResourceName)
}

func (f *deploymentFreezeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.DeploymentFreezeSchema{}.GetResourceSchema()
}

func (f *deploymentFreezeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	f.Config = ResourceConfiguration(req, resp)
}

func (f *deploymentFreezeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var state *deploymentFreezeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deploymentFreeze, err := deploymentfreezes.GetById(f.Config.Client, state.GetID())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "deployment freeze"); err != nil {
			resp.Diagnostics.AddError("unable to load deployment freeze", err.Error())
		}
		return
	}

	if deploymentFreeze.Name != state.Name.ValueString() {
		state.Name = types.StringValue(deploymentFreeze.Name)
	}

	mapToState(ctx, state, deploymentFreeze)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (f *deploymentFreezeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var plan *deploymentFreezeModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deploymentFreeze, diags := mapFromState(plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	createdFreeze, err := deploymentfreezes.Add(f.Config.Client, deploymentFreeze)
	if err != nil {
		resp.Diagnostics.AddError("error while creating deployment freeze", err.Error())
		return
	}

	diags.Append(mapToState(ctx, plan, createdFreeze)...)
	if diags.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (f *deploymentFreezeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var plan *deploymentFreezeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	existingFreeze, err := deploymentfreezes.GetById(f.Config.Client, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load deployment freeze", err.Error())
		return
	}

	updatedFreeze, diags := mapFromState(plan)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// this resource doesn't include scopes, need to copy it from the fetched resource
	updatedFreeze.ProjectEnvironmentScope = existingFreeze.ProjectEnvironmentScope

	updatedFreeze.SetID(existingFreeze.GetID())
	updatedFreeze.Links = existingFreeze.Links

	updatedFreeze, err = deploymentfreezes.Update(f.Config.Client, updatedFreeze)
	if err != nil {
		resp.Diagnostics.AddError("error while updating deployment freeze", err.Error())
		return
	}

	diags.Append(mapToState(ctx, plan, updatedFreeze)...)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (f *deploymentFreezeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var state *deploymentFreezeModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	freeze, err := deploymentfreezes.GetById(f.Config.Client, state.GetID())
	if err != nil {
		resp.Diagnostics.AddError("unable to load deployment freeze", err.Error())
		return
	}

	err = deploymentfreezes.Delete(f.Config.Client, freeze)
	if err != nil {
		resp.Diagnostics.AddError("unable to delete deployment freeze", err.Error())
	}

	resp.State.RemoveResource(ctx)
}

func mapToState(ctx context.Context, state *deploymentFreezeModel, deploymentFreeze *deploymentfreezes.DeploymentFreeze) diag.Diagnostics {
	state.ID = types.StringValue(deploymentFreeze.ID)
	state.Name = types.StringValue(deploymentFreeze.Name)

	updatedStart, diags := calculateStateTime(ctx, state.Start, *deploymentFreeze.Start)
	if diags.HasError() {
		return diags
	}
	state.Start = updatedStart

	updatedEnd, diags := calculateStateTime(ctx, state.End, *deploymentFreeze.End)
	if diags.HasError() {
		return diags
	}
	state.End = updatedEnd

	return nil
}

func calculateStateTime(ctx context.Context, stateValue timetypes.RFC3339, updatedValue time.Time) (timetypes.RFC3339, diag.Diagnostics) {
	stateTime, diags := stateValue.ValueRFC3339Time()
	if diags.HasError() {
		return timetypes.RFC3339{}, diags
	}
	stateTimeUTC := timetypes.NewRFC3339TimeValue(stateTime.UTC())
	updatedValueUTC := updatedValue.UTC()
	valuesAreEqual, diags := stateTimeUTC.StringSemanticEquals(ctx, timetypes.NewRFC3339TimeValue(updatedValueUTC))
	if diags.HasError() {
		return timetypes.NewRFC3339Null(), diags
	}

	if valuesAreEqual {
		return stateValue, diags
	}

	location := stateTime.Location()
	newValue := timetypes.NewRFC3339TimeValue(updatedValueUTC.In(location))
	return newValue, diags
}

func mapFromState(state *deploymentFreezeModel) (*deploymentfreezes.DeploymentFreeze, diag.Diagnostics) {
	start, diags := state.Start.ValueRFC3339Time()
	if diags.HasError() {
		return nil, diags
	}
	start = start.UTC()

	end, diags := state.End.ValueRFC3339Time()
	if diags.HasError() {
		return nil, diags
	}
	end = end.UTC()

	freeze := deploymentfreezes.DeploymentFreeze{
		Name:  state.Name.ValueString(),
		Start: &start,
		End:   &end,
	}

	freeze.ID = state.ID.String()
	return &freeze, nil
}
