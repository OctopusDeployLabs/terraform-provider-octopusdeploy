package octopusdeploy_framework

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithImportState = &runbookTypeResource{}

type runbookTypeResource struct {
	*Config
}

func NewRunbookResource() resource.Resource {
	return &runbookTypeResource{}
}

func (*runbookTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.RunbookResourceDescription)
}

func (*runbookTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetRunbookResourceSchema()
}

func (r *runbookTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *runbookTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.RunbookTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	runbook := mapToRunbook(plan)

	util.Create(ctx, schemas.RunbookResourceDescription, plan)

	createdRunbook, err := runbooks.Add(r.Config.Client, runbook)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to create runbook (%s)", runbook.Name), err.Error())
		return
	}

	mapToState(&plan, createdRunbook)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	util.Created(ctx, schemas.RunbookResourceDescription, createdRunbook)
}

func (r *runbookTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.RunbookTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	util.Reading(ctx, schemas.RunbookResourceDescription, state)

	runbook, err := runbooks.GetByID(r.Config.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to load runbook", err.Error())
		return
	}

	mapToState(&state, runbook)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	util.Read(ctx, schemas.RunbookResourceDescription, runbook)
}

func (r *runbookTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state schemas.RunbookTypeResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	util.Update(ctx, schemas.RunbookResourceDescription, data)

	runbook, err := runbooks.GetByID(r.Config.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load runbook", err.Error())
		return
	}

	updatedRunbook := runbooks.NewRunbook(data.Name.ValueString(), data.ProjectID.ValueString())
	updatedRunbook.ID = runbook.GetID()
	updatedRunbook.SpaceID = runbook.SpaceID
	updatedRunbook.Description = data.Description.ValueString()
	updatedRunbook.RunbookProcessID = data.RunbookProcessID.ValueString()
	updatedRunbook.PublishedRunbookSnapshotID = data.PublishedRunbookSnapshotID.ValueString()
	if !data.MultiTenancyMode.IsNull() {
		updatedRunbook.MultiTenancyMode = core.TenantedDeploymentMode(data.MultiTenancyMode.ValueString())
	}
	updatedRunbook.ConnectivityPolicy = schemas.MapToConnectivityPolicy(data.ConnectivityPolicy)
	updatedRunbook.EnvironmentScope = data.EnvironmentScope.ValueString()
	updatedRunbook.Environments = util.ExpandStringList(data.Environments)
	updatedRunbook.DefaultGuidedFailureMode = data.DefaultGuidedFailureMode.ValueString()
	updatedRunbook.RunRetentionPolicy = schemas.MapToRunbookRetentionPeriod(data.RunRetentionPolicy)
	updatedRunbook.ForcePackageDownload = data.ForcePackageDownload.ValueBool()

	updatedRunbook, err = runbooks.Update(r.Config.Client, updatedRunbook)
	if err != nil {
		resp.Diagnostics.AddError("failed to update runbook", err.Error())
	}

	util.Updated(ctx, schemas.RunbookResourceDescription, updatedRunbook)

	mapToState(&data, updatedRunbook)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (*runbookTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *runbookTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.RunbookTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	util.Delete(ctx, schemas.RunbookResourceDescription, state)

	if err := runbooks.DeleteByID(r.Config.Client, state.SpaceID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("failed to delete runbook", err.Error())
		return
	}

	util.Deleted(ctx, schemas.RunbookResourceDescription, state)
	resp.State.RemoveResource(ctx)
}

func mapToRunbook(data schemas.RunbookTypeResourceModel) *runbooks.Runbook {
	name := data.Name.ValueString()
	projectId := data.ProjectID.ValueString()

	runbook := runbooks.NewRunbook(name, projectId)
	if !data.ID.IsNull() {
		runbook.ID = data.ID.ValueString()
	}

	if !data.Description.IsNull() {
		runbook.Description = data.Description.ValueString()
	}
	if !data.RunbookProcessID.IsNull() {
		runbook.RunbookProcessID = data.RunbookProcessID.ValueString()
	}
	if !data.PublishedRunbookSnapshotID.IsNull() {
		runbook.PublishedRunbookSnapshotID = data.PublishedRunbookSnapshotID.ValueString()
	}
	if !data.SpaceID.IsNull() {
		runbook.SpaceID = data.SpaceID.ValueString()
	}
	if !data.MultiTenancyMode.IsNull() {
		runbook.MultiTenancyMode = core.TenantedDeploymentMode(data.MultiTenancyMode.ValueString())
	}
	if !data.ConnectivityPolicy.IsNull() {
		runbook.ConnectivityPolicy = schemas.MapToConnectivityPolicy(data.ConnectivityPolicy)
	}
	if !data.EnvironmentScope.IsNull() {
		runbook.EnvironmentScope = data.EnvironmentScope.ValueString()
	}
	if !data.Environments.IsNull() {
		runbook.Environments = util.ExpandStringList(data.Environments)
	}
	if !data.DefaultGuidedFailureMode.IsNull() {
		runbook.DefaultGuidedFailureMode = data.DefaultGuidedFailureMode.ValueString()
	}
	if !data.RunRetentionPolicy.IsNull() {
		runbook.RunRetentionPolicy = schemas.MapToRunbookRetentionPeriod(data.RunRetentionPolicy)
	}
	if !data.ForcePackageDownload.IsNull() {
		runbook.ForcePackageDownload = data.ForcePackageDownload.ValueBool()
	}

	return runbook
}

func mapToState(data *schemas.RunbookTypeResourceModel, runbook *runbooks.Runbook) {
	data.ID = types.StringValue(runbook.ID)
	data.Name = types.StringValue(runbook.Name)
	data.ProjectID = types.StringValue(runbook.ProjectID)
	data.Description = types.StringValue(runbook.Description)
	data.RunbookProcessID = types.StringValue(runbook.RunbookProcessID)
	data.PublishedRunbookSnapshotID = types.StringValue(runbook.PublishedRunbookSnapshotID)
	data.SpaceID = types.StringValue(runbook.SpaceID)
	data.MultiTenancyMode = types.StringValue(string(runbook.MultiTenancyMode))
	data.ConnectivityPolicy = types.ListValueMust(
		types.ObjectType{AttrTypes: schemas.GetConnectivityPolicyObjectType()},
		[]attr.Value{
			schemas.MapFromConnectivityPolicy(runbook.ConnectivityPolicy),
		},
	)
	data.EnvironmentScope = types.StringValue(runbook.EnvironmentScope)
	data.Environments = util.FlattenStringList(runbook.Environments)
	data.DefaultGuidedFailureMode = types.StringValue(runbook.DefaultGuidedFailureMode)
	data.ForcePackageDownload = types.BoolValue(runbook.ForcePackageDownload)
	data.RunRetentionPolicy = types.ListValueMust(
		types.ObjectType{AttrTypes: schemas.GetRunbookRetentionPeriodObjectType()},
		[]attr.Value{
			schemas.MapFromRunbookRetentionPeriod(runbook.RunRetentionPolicy),
		},
	)
}
