package octopusdeploy_framework

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deploymentfreezes"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
)

type deploymentFreezeTenantResource struct {
	*Config
}

const tenantDescription = "deployment freeze tenant scope"

var _ resource.Resource = &deploymentFreezeTenantResource{}
var _ resource.ResourceWithConfigure = &deploymentFreezeTenantResource{}

func NewDeploymentFreezeTenantResource() resource.Resource {
	return &deploymentFreezeTenantResource{}
}

func (d *deploymentFreezeTenantResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("deployment_freeze_tenant")
}

func (d *deploymentFreezeTenantResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.DeploymentFreezeTenantSchema{}.GetResourceSchema()
}

func (d *deploymentFreezeTenantResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	d.Config = ResourceConfiguration(req, resp)
}

func (d *deploymentFreezeTenantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	util.Create(ctx, tenantDescription)

	var plan schemas.DeploymentFreezeTenantResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("adding tenant (%s) to deployment freeze (%s)", plan.TenantID.ValueString(), plan.DeploymentFreezeID.ValueString()))
	freeze, err := deploymentfreezes.GetById(d.Client, plan.DeploymentFreezeID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("cannot load deployment freeze", err.Error())
		return
	}

	// Create new tenant scope
	tenantScope := deploymentfreezes.TenantProjectEnvironment{
		TenantId:      plan.TenantID.ValueString(),
		ProjectId:     plan.ProjectID.ValueString(),
		EnvironmentId: plan.EnvironmentID.ValueString(),
	}

	// Add to existing scopes
	freeze.TenantProjectEnvironmentScope = append(freeze.TenantProjectEnvironmentScope, tenantScope)

	tflog.Info(ctx, fmt.Sprintf("[API Request] Updating deployment freeze with new scope. Total scopes: %d", len(freeze.TenantProjectEnvironmentScope)))

	freeze, err = deploymentfreezes.Update(d.Client, freeze)
	if err != nil {
		resp.Diagnostics.AddError("error while updating deployment freeze", err.Error())
		return
	}

	if freezeJSON, err := json.MarshalIndent(freeze, "", "  "); err == nil {
		tflog.Info(ctx, fmt.Sprintf("[API Response] Updated deployment freeze:\n%s", string(freezeJSON)))
	}

	plan.ID = types.StringValue(util.BuildCompositeId(plan.DeploymentFreezeID.ValueString(), plan.TenantID.ValueString(), plan.ProjectID.ValueString(), plan.EnvironmentID.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("tenant scope (%s) added to deployment freeze (%s)", plan.TenantID.ValueString(), plan.DeploymentFreezeID.ValueString()))
	util.Created(ctx, tenantDescription)
}
func (d *deploymentFreezeTenantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	util.Reading(ctx, tenantDescription)

	var data schemas.DeploymentFreezeTenantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bits := util.SplitCompositeId(data.ID.ValueString())
	freezeId := bits[0]
	tenantId := bits[1]
	projectId := bits[2]
	environmentId := bits[3]

	freeze, err := deploymentfreezes.GetById(d.Client, freezeId)
	if err != nil {
		apiError, ok := err.(*core.APIError)
		if !ok {
			resp.Diagnostics.AddError("unable to load deployment freeze", err.Error())
			return
		}

		if apiError.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("unable to load deployment freeze", apiError.Error())
		return
	}

	exists := false
	for _, scope := range freeze.TenantProjectEnvironmentScope {
		if scope.TenantId == tenantId && scope.ProjectId == projectId && scope.EnvironmentId == environmentId {
			exists = true
			break
		}
	}

	if !exists {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	util.Read(ctx, tenantDescription)
}

func (d *deploymentFreezeTenantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	util.Update(ctx, tenantDescription)

	var plan, state schemas.DeploymentFreezeTenantResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	freeze, err := deploymentfreezes.GetById(d.Client, state.DeploymentFreezeID.ValueString())
	if err != nil {
		apiError := err.(*core.APIError)
		if apiError.StatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError("unable to load deployment freeze", err.Error())
			return
		}
	}

	// Remove old scope
	bits := util.SplitCompositeId(state.ID.ValueString())
	oldTenantId := bits[1]
	oldProjectId := bits[2]
	oldEnvironmentId := bits[3]

	newScopes := make([]deploymentfreezes.TenantProjectEnvironment, 0)
	for _, scope := range freeze.TenantProjectEnvironmentScope {
		if scope.TenantId != oldTenantId || scope.ProjectId != oldProjectId || scope.EnvironmentId != oldEnvironmentId {
			newScopes = append(newScopes, scope)
		}
	}

	// Add new scope
	newScopes = append(newScopes, deploymentfreezes.TenantProjectEnvironment{
		TenantId:      plan.TenantID.ValueString(),
		ProjectId:     plan.ProjectID.ValueString(),
		EnvironmentId: plan.EnvironmentID.ValueString(),
	})

	freeze.TenantProjectEnvironmentScope = newScopes

	freeze, err = deploymentfreezes.Update(d.Client, freeze)
	if err != nil {
		resp.Diagnostics.AddError("error while updating deployment freeze", err.Error())
		return
	}

	plan.ID = types.StringValue(util.BuildCompositeId(plan.DeploymentFreezeID.ValueString(), plan.TenantID.ValueString(), plan.ProjectID.ValueString(), plan.EnvironmentID.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("updated tenant scope (%s) in deployment freeze (%s)", plan.TenantID.ValueString(), plan.DeploymentFreezeID.ValueString()))
	util.Updated(ctx, tenantDescription)
}

func (d *deploymentFreezeTenantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	util.Delete(ctx, tenantDescription)

	var data schemas.DeploymentFreezeTenantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	freeze, err := deploymentfreezes.GetById(d.Client, data.DeploymentFreezeID.ValueString())
	if err != nil {
		apiError := err.(*core.APIError)
		if apiError.StatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError("unable to load deployment freeze", err.Error())
			return
		}
	}

	// Remove the tenant scope
	bits := util.SplitCompositeId(data.ID.ValueString())
	tenantId := bits[1]
	projectId := bits[2]
	environmentId := bits[3]

	newScopes := make([]deploymentfreezes.TenantProjectEnvironment, 0)
	for _, scope := range freeze.TenantProjectEnvironmentScope {
		if scope.TenantId != tenantId || scope.ProjectId != projectId || scope.EnvironmentId != environmentId {
			newScopes = append(newScopes, scope)
		}
	}

	freeze.TenantProjectEnvironmentScope = newScopes

	freeze, err = deploymentfreezes.Update(d.Client, freeze)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("cannot remove tenant scope (%s) from deployment freeze (%s)", data.TenantID.ValueString(), data.DeploymentFreezeID.ValueString()), err.Error())
	}

	tflog.Debug(ctx, fmt.Sprintf("tenant scope (%s) removed from deployment freeze (%s)", data.TenantID.ValueString(), data.DeploymentFreezeID.ValueString()))
	util.Deleted(ctx, tenantDescription)
}
