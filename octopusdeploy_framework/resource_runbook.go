package octopusdeploy_framework

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
	resp.Schema = schemas.RunbookSchema{}.GetResourceSchema()
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

	name := plan.Name.ValueString()
	projectId := plan.ProjectID.ValueString()

	runbook := runbooks.NewRunbook(name, projectId)
	if !plan.ID.IsNull() {
		runbook.ID = plan.ID.ValueString()
	}

	runbook.Description = plan.Description.ValueString()
	runbook.RunbookProcessID = plan.RunbookProcessID.ValueString()
	runbook.PublishedRunbookSnapshotID = plan.PublishedRunbookSnapshotID.ValueString()
	runbook.SpaceID = plan.SpaceID.ValueString()
	if !plan.MultiTenancyMode.IsNull() {
		runbook.MultiTenancyMode = core.TenantedDeploymentMode(plan.MultiTenancyMode.ValueString())
	}
	runbook.ConnectivityPolicy = schemas.MapToConnectivityPolicy(plan.ConnectivityPolicy)
	runbook.EnvironmentScope = plan.EnvironmentScope.ValueString()
	runbook.Environments = util.ExpandStringList(plan.Environments)
	runbook.DefaultGuidedFailureMode = plan.DefaultGuidedFailureMode.ValueString()
	runbook.RunRetentionPolicy = schemas.MapToRunbookRetentionPeriod(plan.RunRetentionPolicy)
	runbook.ForcePackageDownload = plan.ForcePackageDownload.ValueBool()

	util.Create(ctx, schemas.RunbookResourceDescription, plan)

	runbooksAreInGit, err := internal.CheckRunbookInGit(r.Config.Client, runbook.SpaceID, runbook.ProjectID)

	if err != nil {
		resp.Diagnostics.AddError("failed to check runbook git ref", err.Error())
	}

	if runbooksAreInGit {
		resp.Diagnostics.AddWarning("Unable to manage CaC Runbooks via Terraform", "Runbook is in git, skipping create")
		resp.State.Set(ctx, &plan)
		return
	}

	createdRunbook, err := runbooks.Add(r.Config.Client, runbook)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to create runbook (%s)", runbook.Name), err.Error())
		return
	}

	resp.Diagnostics.Append(plan.RefreshFromApiResponse(ctx, createdRunbook)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	runbooksAreInGit, err := internal.CheckRunbookInGit(r.Config.Client, state.SpaceID.ValueString(), state.ProjectID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to check runbook git ref", err.Error())
		return
	}

	if runbooksAreInGit {
		resp.Diagnostics.AddWarning("Unable to manage CaC Runbooks via Terraform", "Runbook is in git, skipping read")
		resp.State.Set(ctx, &state)
		return
	}

	runbook, err := runbooks.GetByID(r.Config.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, schemas.RunbookResourceDescription); err != nil {
			resp.Diagnostics.AddError("failed to load runbook", err.Error())
		}
		return
	}

	resp.Diagnostics.Append(state.RefreshFromApiResponse(ctx, runbook)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	util.Read(ctx, schemas.RunbookResourceDescription, runbook)
}

func (r *runbookTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state schemas.RunbookTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	util.Update(ctx, schemas.RunbookResourceDescription, plan)

	runbooksAreInGit, err := internal.CheckRunbookInGit(r.Config.Client, plan.SpaceID.ValueString(), plan.ProjectID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Unable to verify Projects Runbooks persistence settings", err.Error())
		return
	}

	if runbooksAreInGit {
		resp.Diagnostics.AddWarning("Unable to manage CaC Runbooks via Terraform", "Runbook is in git, skipping update")
		resp.State.Set(ctx, &plan)
		return
	}

	runbook, err := runbooks.GetByID(r.Config.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load runbook", err.Error())
		return
	}

	updatedRunbook := runbooks.NewRunbook(plan.Name.ValueString(), plan.ProjectID.ValueString())
	updatedRunbook.ID = runbook.GetID()
	updatedRunbook.SpaceID = runbook.SpaceID
	updatedRunbook.Description = plan.Description.ValueString()
	updatedRunbook.RunbookProcessID = plan.RunbookProcessID.ValueString()
	updatedRunbook.PublishedRunbookSnapshotID = plan.PublishedRunbookSnapshotID.ValueString()
	if !plan.MultiTenancyMode.IsNull() {
		updatedRunbook.MultiTenancyMode = core.TenantedDeploymentMode(plan.MultiTenancyMode.ValueString())
	}
	updatedRunbook.ConnectivityPolicy = schemas.MapToConnectivityPolicy(plan.ConnectivityPolicy)
	updatedRunbook.EnvironmentScope = plan.EnvironmentScope.ValueString()
	updatedRunbook.Environments = util.ExpandStringList(plan.Environments)
	updatedRunbook.DefaultGuidedFailureMode = plan.DefaultGuidedFailureMode.ValueString()
	updatedRunbook.RunRetentionPolicy = schemas.MapToRunbookRetentionPeriod(plan.RunRetentionPolicy)
	updatedRunbook.ForcePackageDownload = plan.ForcePackageDownload.ValueBool()

	updatedRunbook, err = runbooks.Update(r.Config.Client, updatedRunbook)
	if err != nil {
		resp.Diagnostics.AddError("failed to update runbook", err.Error())
	}

	resp.Diagnostics.Append(plan.RefreshFromApiResponse(ctx, updatedRunbook)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	util.Updated(ctx, schemas.RunbookResourceDescription, updatedRunbook)
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

	runbooksAreInGit, err := internal.CheckRunbookInGit(r.Config.Client, state.SpaceID.ValueString(), state.ProjectID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Unable to verify Projects Runbooks persistence settings", err.Error())

		return
	}

	if runbooksAreInGit {
		resp.Diagnostics.AddWarning("Unable to manage CaC Runbooks via Terraform", "Runbook is in git, skipping delete")

		resp.State.RemoveResource(ctx)
		return
	}

	if err := runbooks.DeleteByID(r.Config.Client, state.SpaceID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("failed to delete runbook", err.Error())
		return
	}

	util.Deleted(ctx, schemas.RunbookResourceDescription, state)
	resp.State.RemoveResource(ctx)
}
