package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type builtInTriggerResource struct {
	*Config
}

func NewBuiltInTriggerResource() resource.Resource {
	return &builtInTriggerResource{}
}

var _ resource.ResourceWithImportState = &builtInTriggerResource{}

func (r *builtInTriggerResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("built_in_trigger")
}

func (r *builtInTriggerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.BuiltInTriggerSchema{}.GetResourceSchema()
}

func (r *builtInTriggerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *builtInTriggerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *builtInTriggerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.BuiltInTriggerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectId := data.ProjectID.ValueString()
	project, err := projects.GetByID(r.Client, data.SpaceID.ValueString(), projectId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read associated project for built-in trigger", err.Error())
		return
	}

	mapBuiltInTriggerFromState(&data, project)

	_, err = projects.Update(r.Client, project)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update associated project for built-in trigger", err.Error())
		return
	}

	// Reload project in case different values were computed for release strategy
	updatedProject, err := projects.GetByID(r.Client, data.SpaceID.ValueString(), projectId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read associated project for built-in trigger", err.Error())
		return
	}

	exists := mapBuiltInTriggerToState(updatedProject, &data)
	if !exists {
		resp.Diagnostics.AddError("Failed to map built-in trigger from updated project", "Release strategy or package are missing")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *builtInTriggerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.BuiltInTriggerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectId := state.ProjectID.ValueString()
	project, err := projects.GetByID(r.Client, state.SpaceID.ValueString(), projectId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read associated project for built-in trigger", err.Error())
		return
	}

	exists := mapBuiltInTriggerToState(project, &state)
	if !exists {
		// Remove from state when release creation strategy or associated package are missing from the project
		tflog.Info(ctx, fmt.Sprintf("unable to find built-in trigger from project (id: %s), removing from state ...", projectId))
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *builtInTriggerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data schemas.BuiltInTriggerResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectId := data.ProjectID.ValueString()
	existingProject, err := projects.GetByID(r.Client, data.SpaceID.ValueString(), projectId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read associated project for built-in trigger", err.Error())
		return
	}

	mapBuiltInTriggerFromState(&data, existingProject)

	_, err = projects.Update(r.Client, existingProject)
	if err != nil {
		resp.Diagnostics.AddError("Error updating associated project for built-in trigger", err.Error())
		return
	}

	updatedProject, err := projects.GetByID(r.Client, data.SpaceID.ValueString(), projectId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read associated project for built-in trigger", err.Error())
		return
	}

	exists := mapBuiltInTriggerToState(updatedProject, &data)
	if !exists {
		resp.Diagnostics.AddError("Failed to map built-in trigger from updated project", "Release strategy or package are missing")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *builtInTriggerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.BuiltInTriggerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectId := state.ProjectID.ValueString()
	project, err := projects.GetByID(r.Client, state.SpaceID.ValueString(), projectId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read associated project for built-in trigger", err.Error())
		return
	}

	project.ReleaseCreationStrategy = &projects.ReleaseCreationStrategy{}
	project.AutoCreateRelease = false

	_, err = projects.Update(r.Client, project)
	if err != nil {
		resp.Diagnostics.AddError("Error updating project to remove release creation strategy(built-in trigger)", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func mapBuiltInTriggerFromState(state *schemas.BuiltInTriggerResourceModel, project *projects.Project) {
	var packageStepId string
	configuredPackageStepId := state.ReleaseCreationPackageStepID.ValueString()
	if configuredPackageStepId != "" {
		packageStepId = configuredPackageStepId
	}

	project.AutoCreateRelease = true
	project.ReleaseCreationStrategy = &projects.ReleaseCreationStrategy{
		ReleaseCreationPackageStepID: packageStepId,
		ChannelID:                    state.ChannelID.ValueString(),
		ReleaseCreationPackage: &packages.DeploymentActionPackage{
			DeploymentAction: state.ReleaseCreationPackage.DeploymentAction.ValueString(),
			PackageReference: state.ReleaseCreationPackage.PackageReference.ValueString(),
		},
	}
}

func mapBuiltInTriggerToState(project *projects.Project, state *schemas.BuiltInTriggerResourceModel) bool {
	if project.ReleaseCreationStrategy == nil {
		return false
	}

	if project.ReleaseCreationStrategy.ReleaseCreationPackage == nil {
		return false
	}

	releaseStrategy := project.ReleaseCreationStrategy

	if releaseStrategy.ReleaseCreationPackageStepID != "" {
		state.ReleaseCreationPackageStepID = types.StringValue(releaseStrategy.ReleaseCreationPackageStepID)
	}

	state.ChannelID = types.StringValue(releaseStrategy.ChannelID)
	state.ReleaseCreationPackage.PackageReference = types.StringValue(releaseStrategy.ReleaseCreationPackage.PackageReference)
	state.ReleaseCreationPackage.DeploymentAction = types.StringValue(releaseStrategy.ReleaseCreationPackage.DeploymentAction)
	state.SpaceID = types.StringValue(project.SpaceID)

	return true
}
