package octopusdeploy_framework

import (
	"context"
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.ResourceWithModifyPlan  = &processResource{}
	_ resource.ResourceWithImportState = &processResource{}
)

type processResource struct {
	*Config
}

func NewProcessResource() resource.Resource {
	return &processResource{}
}

func (r *processResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProcessResourceName)
}

func (r *processResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProcessSchema{}.GetResourceSchema()
}

func (r *processResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *processResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *processResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddWarning("Deleting process", "Applying this resource destruction will not delete process and it's steps")
		return
	}
}

func (r *processResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.ProcessResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	projectId := data.OwnerID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("creating process for owner: %s", projectId))

	// Empty process is created as part of project creation
	project, err := projects.GetByID(r.Config.Client, spaceId, projectId)
	if err != nil {
		resp.Diagnostics.AddError("Error creating process, unable to find associated project", err.Error())
		return
	}

	if project.PersistenceSettings != nil && project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		resp.Diagnostics.AddWarning("Cannot create process", "Project persisted under version control system. Process of version controlled project cannot be created")
		return
	}

	process, err := deployments.GetDeploymentProcessByID(r.Config.Client, spaceId, project.DeploymentProcessID)
	if err != nil {
		resp.Diagnostics.AddError("Error creating process", err.Error())
		return
	}

	mapProcessToState(process, data)

	tflog.Info(ctx, fmt.Sprintf("process created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.ProcessResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	projectId := data.OwnerID.ValueString()
	processId := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("reading process (%s)", processId))

	process, diags := loadProcess(r.Config.Client, spaceId, projectId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	mapProcessToState(process, data)

	tflog.Info(ctx, fmt.Sprintf("process read (%s)", process.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *schemas.ProcessResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	projectId := data.OwnerID.ValueString()
	processId := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("updating process (%s)", data.ID))

	process, diags := loadProcess(r.Config.Client, spaceId, projectId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Nothing to update, when owner_id is changed we want to replace this resource with process from another owner

	mapProcessToState(process, data)

	tflog.Info(ctx, fmt.Sprintf("process updated (%s)", process.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *schemas.ProcessResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	projectId := data.OwnerID.ValueString()
	processId := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("deleting process (%s)", data.ID))

	_, diags := loadProcess(r.Config.Client, spaceId, projectId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Do nothing, because process can not be deleted from the project
	resp.Diagnostics.AddWarning("Deleting process", "Destruction of this resource will not delete process from the project")

	resp.State.RemoveResource(ctx)
}

func mapProcessToState(process *deployments.DeploymentProcess, state *schemas.ProcessResourceModel) {
	state.ID = types.StringValue(process.ID)
	state.SpaceID = types.StringValue(process.SpaceID)
	state.OwnerID = types.StringValue(process.ProjectID)
}

func loadProcess(client *client.Client, spaceId string, projectId string, processId string) (*deployments.DeploymentProcess, diag.Diagnostics) {
	process, processError := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if processError == nil {
		return process, diag.Diagnostics{}
	}

	processNotFound := diag.Diagnostics{}
	processNotFound.AddError("unable to load process", processError.Error())

	var apiError *core.APIError
	if errors.As(processError, &apiError) && apiError.StatusCode == http.StatusNotFound {
		// Try to load corresponding project to check if it's version controlled
		project, err := projects.GetByID(client, spaceId, projectId)
		if err != nil {
			return nil, processNotFound // return original error when project cannot be loaded
		}

		if project.PersistenceSettings == nil {
			return nil, processNotFound
		}

		if project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
			versionControlled := diag.Diagnostics{}
			versionControlled.AddWarning("process persisted under version control system", "Version controlled resources will not be modified via terraform")
			return nil, versionControlled
		}
	}

	return nil, processNotFound
}

func loadProcessForSteps(client *client.Client, spaceId string, processId string) (*deployments.DeploymentProcess, diag.Diagnostics) {
	projectId := ""

	// Assumes that project id is part of the process identifier
	// This approach allows us to avoid dependencies between resources and avoid adding "project_id" to all process resources
	const prefix = "deploymentprocess-"
	if strings.HasPrefix(processId, prefix) {
		projectId = processId[len(prefix):]
	}

	return loadProcess(client, spaceId, projectId, processId)
}
