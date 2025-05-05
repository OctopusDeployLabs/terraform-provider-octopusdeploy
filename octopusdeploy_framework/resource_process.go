package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbookprocess"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

func (r *processResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProcessResourceName)
}

func (r *processResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProcessSchema{}.GetResourceSchema()
}

func (r *processResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *processResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	process, diags := loadProcessWrapperByProcessId(r.Config.Client, r.Config.SpaceID, request.ID)
	if len(diags) > 0 {
		response.Diagnostics.Append(diags...)
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("project_id"), process.GetProjectID())...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), process.GetID())...)
}

func (r *processResource) ModifyPlan(_ context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
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
	projectId := data.ProjectID.ValueString()
	runbookId := data.RunbookID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("creating process for owner: %s", projectId))

	project, projectError := projects.GetByID(r.Config.Client, spaceId, projectId)
	if projectError != nil {
		resp.Diagnostics.AddError("Error creating process, unable to find associated project", projectError.Error())
		return
	}

	if project.PersistenceSettings != nil && project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		resp.Diagnostics.AddError("Cannot create process for version controlled project", "Version controlled resources will not be modified via terraform")
		return
	}

	// Empty process is created as part of the project or runbook creation
	var process processWrapper
	if runbookId != "" {
		runbook, runbookError := runbooks.GetByID(r.Config.Client, spaceId, data.RunbookID.ValueString())
		if runbookError != nil {
			resp.Diagnostics.AddError("Error creating process, unable to find associated runbook", runbookError.Error())
			return
		}

		if runbook.ProjectID != project.ID {
			resp.Diagnostics.AddError("Error creating process", "Provided runbook does not belong to the given project")
			return
		}

		runbookProcess, processError := runbookprocess.GetByID(r.Config.Client, spaceId, runbook.RunbookProcessID)
		if processError != nil {
			resp.Diagnostics.AddError("Error creating runbook process", processError.Error())
			return
		}

		process = runbookProcessWrapper{runbookProcess}
	} else {
		deploymentProcess, processError := deployments.GetDeploymentProcessByID(r.Config.Client, spaceId, project.DeploymentProcessID)
		if processError != nil {
			resp.Diagnostics.AddError("Error creating deployment process", processError.Error())
			return
		}

		process = deploymentProcessWrapper{deploymentProcess}
	}

	process.PopulateState(data)

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
	projectId := data.ProjectID.ValueString()
	processId := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("reading process (%s)", processId))

	process, diags := loadProcessWrapper(r.Config.Client, spaceId, projectId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	process.PopulateState(data)

	tflog.Info(ctx, fmt.Sprintf("process read (%s)", process.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *schemas.ProcessResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	projectId := data.ProjectID.ValueString()
	processId := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("updating process (%s)", data.ID))

	process, diags := loadProcessWrapper(r.Config.Client, spaceId, projectId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Nothing to update, when projectId or runbookId are changed we want to replace this resource with process from another owner
	process.PopulateState(data)

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
	projectId := data.ProjectID.ValueString()
	processId := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("deleting process (%s)", data.ID))

	_, diags := loadProcessWrapper(r.Config.Client, spaceId, projectId, processId)
	if len(diags) > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Do nothing, because process can not be deleted from the project
	resp.Diagnostics.AddWarning("Deleting process", "Destruction of this resource will not remove the process from the system")

	resp.State.RemoveResource(ctx)
}
