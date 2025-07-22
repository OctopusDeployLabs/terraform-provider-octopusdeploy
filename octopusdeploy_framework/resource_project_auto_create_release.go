package octopusdeploy_framework

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type projectAutoCreateReleaseResource struct {
	*Config
}

func NewProjectAutoCreateReleaseResource() resource.Resource {
	return &projectAutoCreateReleaseResource{}
}

var _ resource.ResourceWithImportState = &projectAutoCreateReleaseResource{}

func (r *projectAutoCreateReleaseResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("project_auto_create_release")
}

func (r *projectAutoCreateReleaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProjectAutoCreateReleaseSchema{}.GetResourceSchema()
}

func (r *projectAutoCreateReleaseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *projectAutoCreateReleaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: project_id
	projectID := req.ID

	// Create empty state with project ID
	empty := &schemas.ProjectAutoCreateReleaseResourceModel{
		ID:                           types.StringValue(projectID + "-auto-create-release"),
		ProjectID:                    types.StringValue(projectID),
		SpaceID:                      types.StringNull(),
		ChannelID:                    types.StringNull(),
		ReleaseCreationPackageStepID: types.StringNull(),
		ReleaseCreationPackage:       []schemas.ProjectAutoCreateReleaseCreationPackage{},
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, empty)...)
}

func (r *projectAutoCreateReleaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.ProjectAutoCreateReleaseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := data.ProjectID.ValueString()
	spaceID := data.SpaceID.ValueString()

	// Get the project
	project, err := projects.GetByID(r.Client, spaceID, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read project", fmt.Sprintf("Unable to read project with ID %s: %s", projectID, err.Error()))
		return
	}

	// Set space ID if not provided
	if data.SpaceID.IsNull() || data.SpaceID.IsUnknown() {
		data.SpaceID = types.StringValue(project.SpaceID)
		spaceID = project.SpaceID
	}

	// Validate the auto create release configuration
	if err := r.validateAutoCreateReleaseConfiguration(ctx, project, &data); err != nil {
		resp.Diagnostics.AddError("Invalid auto create release configuration", err.Error())
		return
	}

	// Configure auto create release
	r.mapDataToProject(ctx, &data, project)

	// Update the project
	_, err = projects.Update(r.Client, project)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update project", fmt.Sprintf("Unable to update project with ID %s: %s", projectID, err.Error()))
		return
	}

	// Set the resource ID
	data.ID = types.StringValue(fmt.Sprintf("%s-auto-create-release", projectID))

	// Re-read the project to get computed values
	updatedProject, err := projects.GetByID(r.Client, spaceID, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read updated project", fmt.Sprintf("Unable to read updated project with ID %s: %s", projectID, err.Error()))
		return
	}

	// Map any computed values back to state
	r.mapProjectToData(ctx, updatedProject, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *projectAutoCreateReleaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.ProjectAutoCreateReleaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := state.ProjectID.ValueString()
	spaceID := state.SpaceID.ValueString()

	// Get the project
	project, err := projects.GetByID(r.Client, spaceID, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read project", fmt.Sprintf("Unable to read project with ID %s: %s", projectID, err.Error()))
		return
	}

	// Check if auto create release is still configured
	if !r.isAutoCreateReleaseConfigured(project) {
		tflog.Info(ctx, fmt.Sprintf("Auto create release not configured for project %s, removing from state", projectID))
		resp.State.RemoveResource(ctx)
		return
	}

	// Map project data to state
	r.mapProjectToData(ctx, project, &state)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *projectAutoCreateReleaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schemas.ProjectAutoCreateReleaseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := data.ProjectID.ValueString()
	spaceID := data.SpaceID.ValueString()

	// Get the project
	project, err := projects.GetByID(r.Client, spaceID, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read project", fmt.Sprintf("Unable to read project with ID %s: %s", projectID, err.Error()))
		return
	}

	// Validate the auto create release configuration
	if err := r.validateAutoCreateReleaseConfiguration(ctx, project, &data); err != nil {
		resp.Diagnostics.AddError("Invalid auto create release configuration", err.Error())
		return
	}

	// Configure auto create release
	r.mapDataToProject(ctx, &data, project)

	// Update the project
	_, err = projects.Update(r.Client, project)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update project", fmt.Sprintf("Unable to update project with ID %s: %s", projectID, err.Error()))
		return
	}

	// Re-read the project to get computed values
	updatedProject, err := projects.GetByID(r.Client, spaceID, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read updated project", fmt.Sprintf("Unable to read updated project with ID %s: %s", projectID, err.Error()))
		return
	}

	// Map any computed values back to state
	r.mapProjectToData(ctx, updatedProject, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *projectAutoCreateReleaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.ProjectAutoCreateReleaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := state.ProjectID.ValueString()
	spaceID := state.SpaceID.ValueString()

	// Get the project
	project, err := projects.GetByID(r.Client, spaceID, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read project", fmt.Sprintf("Unable to read project with ID %s: %s", projectID, err.Error()))
		return
	}

	// Disable auto create release
	project.AutoCreateRelease = false
	project.ReleaseCreationStrategy = &projects.ReleaseCreationStrategy{}

	// Update the project
	_, err = projects.Update(r.Client, project)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update project", fmt.Sprintf("Unable to update project with ID %s: %s", projectID, err.Error()))
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *projectAutoCreateReleaseResource) mapDataToProject(ctx context.Context, data *schemas.ProjectAutoCreateReleaseResourceModel, project *projects.Project) {
	project.AutoCreateRelease = true
	project.ReleaseCreationStrategy = expand(ctx, data)
}

func (r *projectAutoCreateReleaseResource) mapProjectToData(ctx context.Context, project *projects.Project, data *schemas.ProjectAutoCreateReleaseResourceModel) {
	data.SpaceID = types.StringValue(project.SpaceID)
	flatten(ctx, project.ReleaseCreationStrategy, data)
}
