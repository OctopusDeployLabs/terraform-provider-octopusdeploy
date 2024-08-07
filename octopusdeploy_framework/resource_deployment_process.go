package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/mappers"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type deploymentProcessResource struct {
	*Config
}

// implementation checks
var _ resource.Resource = &deploymentProcessResource{}
var _ resource.ResourceWithImportState = &deploymentProcessResource{}

func NewDeploymentProcessResource() resource.Resource {
	return &deploymentProcessResource{}
}

func (d *deploymentProcessResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("deployment_process")
}

func (d *deploymentProcessResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetDeploymentProcessResourceSchema()
}

func (d *deploymentProcessResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	d.Config = ResourceConfiguration(req, resp)
}

func (d *deploymentProcessResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.DeploymentProcessResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("creating deployment process: %#v", plan))
	spaceID := plan.SpaceID.ValueString()
	project, err := projects.GetByID(d.Client, plan.SpaceID.ValueString(), plan.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error getting project %s", plan.ProjectID.ValueString()), err.Error())
		return
	}

	var current *deployments.DeploymentProcess
	if project.PersistenceSettings != nil && project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		current, err = deployments.GetDeploymentProcessByGitRef(d.Client, spaceID, project, plan.Branch.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by git ref", err.Error())
			return
		}
	} else {
		current, err = deployments.GetDeploymentProcessByID(d.Client, spaceID, project.DeploymentProcessID)
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by ID", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(mappers.MapStateToDeploymentProcess(ctx, &plan, current)...)

	if resp.Diagnostics.HasError() {
		return
	}

	current, err = deployments.UpdateDeploymentProcess(d.Client, current)
	if err != nil {
		resp.Diagnostics.AddError("unable to update deployment process", err.Error())
		return
	}

	resp.Diagnostics.Append(mappers.MapDeploymentProcessToState(ctx, current, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (d *deploymentProcessResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.DeploymentProcessResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading deployment process: %#v", state))
	spaceID := state.SpaceID.ValueString()
	project, err := projects.GetByID(d.Client, state.SpaceID.ValueString(), state.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error getting project %s", state.ProjectID.ValueString()), err.Error())
		return
	}

	var current *deployments.DeploymentProcess
	if project.PersistenceSettings != nil && project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		current, err = deployments.GetDeploymentProcessByGitRef(d.Client, spaceID, project, state.Branch.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by git ref", err.Error())
			return
		}
	} else {
		current, err = deployments.GetDeploymentProcessByID(d.Client, spaceID, project.DeploymentProcessID)
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by ID", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(mappers.MapDeploymentProcessToState(ctx, current, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *deploymentProcessResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.DeploymentProcessResourceModel
	var state schemas.DeploymentProcessResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading deployment process: %#v", plan))
	spaceID := plan.SpaceID.ValueString()
	project, err := projects.GetByID(d.Client, plan.SpaceID.ValueString(), plan.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error getting project %s", plan.ProjectID.ValueString()), err.Error())
		return
	}

	var current *deployments.DeploymentProcess
	if project.PersistenceSettings != nil && project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		current, err = deployments.GetDeploymentProcessByGitRef(d.Client, spaceID, project, plan.Branch.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by git ref", err.Error())
			return
		}
	} else {
		current, err = deployments.GetDeploymentProcessByID(d.Client, spaceID, project.DeploymentProcessID)
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by ID", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(mappers.MapStateToDeploymentProcess(ctx, &plan, current)...)

	if resp.Diagnostics.HasError() {
		return
	}

	current, err = deployments.UpdateDeploymentProcess(d.Client, current)
	if err != nil {
		resp.Diagnostics.AddError("unable to update deployment process", err.Error())
		return
	}

	resp.Diagnostics.Append(mappers.MapDeploymentProcessToState(ctx, current, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (d *deploymentProcessResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.DeploymentProcessResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading deployment process: %#v", state))
	spaceID := state.SpaceID.ValueString()
	project, err := projects.GetByID(d.Client, state.SpaceID.ValueString(), state.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error getting project %s", state.ProjectID.ValueString()), err.Error())
		return
	}

	var current *deployments.DeploymentProcess
	if project.PersistenceSettings != nil && project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		current, err = deployments.GetDeploymentProcessByGitRef(d.Client, spaceID, project, state.Branch.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by git ref", err.Error())
			return
		}
	} else {
		current, err = deployments.GetDeploymentProcessByID(d.Client, spaceID, project.DeploymentProcessID)
		if err != nil {
			resp.Diagnostics.AddError("unable to retrieve deployment process by ID", err.Error())
			return
		}
	}

	current.Steps = []*deployments.DeploymentStep{}
	deployments.UpdateDeploymentProcess(d.Client, current)
	resp.Diagnostics.Append(mappers.MapDeploymentProcessToState(ctx, current, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Info(ctx, "deployment process deleted")
}

func (d *deploymentProcessResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//TODO implement me
	//panic("implement me")
}
