package octopusdeploy_framework

import (
	"context"
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

type deploymentFreezeProjectResource struct {
	*Config
}

const description = "deployment freeze project scope"

var _ resource.Resource = &deploymentFreezeProjectResource{}
var _ resource.ResourceWithConfigure = &deploymentFreezeProjectResource{}

func NewDeploymentFreezeProjectResource() resource.Resource {
	return &deploymentFreezeProjectResource{}
}

func (d *deploymentFreezeProjectResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("deployment_freeze_project")
}

func (d *deploymentFreezeProjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.DeploymentFreezeProjectSchema{}.GetResourceSchema()
}

func (d *deploymentFreezeProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	d.Config = ResourceConfiguration(req, resp)
}

func (d *deploymentFreezeProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	util.Create(ctx, description)

	var plan schemas.DeploymentFreezeProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("adding project (%s) to deployment freeze (%s)", plan.ProjectID.ValueString(), plan.DeploymentFreezeID.ValueString()))
	freeze, err := deploymentfreezes.GetById(d.Client, plan.DeploymentFreezeID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("cannot load deployment freeze", err.Error())
	}
	freeze.ProjectEnvironmentScope[plan.ProjectID.ValueString()] = util.ExpandStringList(plan.EnvironmentIDs)

	freeze, err = deploymentfreezes.Update(d.Client, freeze)
	if err != nil {
		resp.Diagnostics.AddError("error while updating deployment freeze", err.Error())
		return
	}

	plan.ID = types.StringValue(util.BuildCompositeId(plan.DeploymentFreezeID.ValueString(), plan.ProjectID.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("scope for project (%s) added to deployment freeze", plan.ProjectID, plan.DeploymentFreezeID))
	util.Created(ctx, description)
}

func (d *deploymentFreezeProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	util.Reading(ctx, description)
	var data schemas.DeploymentFreezeProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bits := util.SplitCompositeId(data.ID.ValueString())
	freezeId := bits[0]
	projectId := bits[1]

	freeze, err := deploymentfreezes.GetById(d.Client, freezeId)
	if err != nil {
		apiError := err.(*core.APIError)
		if apiError.StatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError("unable to load deployment freeze", err.Error())
			return
		}
	}

	data.EnvironmentIDs = util.FlattenStringList(freeze.ProjectEnvironmentScope[projectId])
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	util.Read(ctx, description)
}

func (d *deploymentFreezeProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	util.Update(ctx, description)

	var plan, state schemas.DeploymentFreezeProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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

	tflog.Debug(ctx, fmt.Sprintf("updating project (%s) to deployment freeze (%s)", plan.ProjectID.ValueString(), plan.DeploymentFreezeID.ValueString()))
	freeze.ProjectEnvironmentScope[plan.ProjectID.ValueString()] = util.ExpandStringList(plan.EnvironmentIDs)
	_, err = deploymentfreezes.Update(d.Client, freeze)
	if err != nil {
		resp.Diagnostics.AddError("error while updating deployment freeze", err.Error())
		return
	}

	plan.ID = types.StringValue(util.BuildCompositeId(plan.DeploymentFreezeID.ValueString(), plan.ProjectID.ValueString()))
	plan.EnvironmentIDs = util.FlattenStringList(freeze.ProjectEnvironmentScope[plan.ProjectID.ValueString()])

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("updated project (%s) to deployment freeze (%s)", plan.ProjectID.ValueString(), plan.DeploymentFreezeID.ValueString()))
	util.Updated(ctx, description)
}

func (d *deploymentFreezeProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	util.Delete(ctx, description)

	var data schemas.DeploymentFreezeProjectResourceModel
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
	tflog.Debug(ctx, fmt.Sprintf("before delete: %#v", freeze))

	delete(freeze.ProjectEnvironmentScope, data.ProjectID.ValueString())
	freeze, err = deploymentfreezes.Update(d.Client, freeze)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("cannot remove project scope (%s) from deployment freeze (%s)", data.ProjectID.ValueString(), data.DeploymentFreezeID.ValueString()), err.Error())
	}

	tflog.Debug(ctx, fmt.Sprintf("scope for project (%s) removed from deployment freeze (%s)", data.ProjectID.ValueString(), data.DeploymentFreezeID.ValueString()))
	util.Deleted(ctx, description)
}
