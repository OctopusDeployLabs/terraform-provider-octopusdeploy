package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type builtInTriggerResource struct {
	*Config
}

func NewBuiltInTriggerResource() resource.Resource {
	return &builtInTriggerResource{}
}

var _ resource.ResourceWithImportState = &builtInTriggerResource{}

func (r *builtInTriggerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("built_in_trigger")
}

func (r *builtInTriggerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "Failed to read associated project for built-in trigger", err.Error())
		return
	}
	releaseStrategy := mapStateToReleaseCreationStrategy(&data)
	project.ReleaseCreationStrategy = releaseStrategy
	project.AutoCreateRelease = true

	_, err = projects.Update(r.Client, project)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "Error updating associated project for built-in trigger", err.Error())
		return
	}

	// Reload project in case different values were computed for release strategy
	updatedProject, err := projects.GetByID(r.Client, data.SpaceID.ValueString(), projectId)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "Failed to read associated project for built-in trigger", err.Error())
		return
	}

	mapReleaseCreationStrategyToState(updatedProject, &data)
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
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "Failed to read associated project for built-in trigger", err.Error())
		return
	}

	mapReleaseCreationStrategyToState(project, &state)

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
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "Failed to read associated project for built-in trigger", err.Error())
		return
	}

	releaseStrategy := mapStateToReleaseCreationStrategy(&data)
	existingProject.ReleaseCreationStrategy = releaseStrategy
	existingProject.AutoCreateRelease = true

	_, err = projects.Update(r.Client, existingProject)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "Error updating associated project for built-in trigger", err.Error())
		return
	}

	updatedProject, err := projects.GetByID(r.Client, data.SpaceID.ValueString(), projectId)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "Failed to read associated project for built-in trigger", err.Error())
		return
	}

	mapReleaseCreationStrategyToState(updatedProject, &data)
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
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "Failed to read associated project for built-in trigger", err.Error())
		return
	}

	project.ReleaseCreationStrategy = &projects.ReleaseCreationStrategy{}
	project.AutoCreateRelease = false

	_, err = projects.Update(r.Client, project)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "Error updating project to remove release creation strategy(built-in trigger)", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func mapStateToReleaseCreationStrategy(state *schemas.BuiltInTriggerResourceModel) *projects.ReleaseCreationStrategy {
	var releaseCreationPackageStepId string
	releaseCreationPackageStepIdString := state.ReleaseCreationPackageStepID.ValueString()
	if releaseCreationPackageStepIdString != "" {
		releaseCreationPackageStepId = releaseCreationPackageStepIdString
	}

	return &projects.ReleaseCreationStrategy{
		ReleaseCreationPackageStepID: releaseCreationPackageStepId,
		ChannelID:                    state.ChannelID.ValueString(),
		ReleaseCreationPackage: &packages.DeploymentActionPackage{
			DeploymentAction: state.ReleaseCreationPackage.DeploymentAction.ValueString(),
			PackageReference: state.ReleaseCreationPackage.PackageReference.ValueString(),
		},
	}
}

func mapReleaseCreationStrategyToState(project *projects.Project, state *schemas.BuiltInTriggerResourceModel) {
	releaseStrategy := project.ReleaseCreationStrategy

	if releaseStrategy.ReleaseCreationPackageStepID != "" {
		state.ReleaseCreationPackageStepID = types.StringValue(releaseStrategy.ReleaseCreationPackageStepID)
	}

	state.ChannelID = types.StringValue(releaseStrategy.ChannelID)
	state.ReleaseCreationPackage.PackageReference = types.StringValue(releaseStrategy.ReleaseCreationPackage.PackageReference)
	state.ReleaseCreationPackage.DeploymentAction = types.StringValue(releaseStrategy.ReleaseCreationPackage.DeploymentAction)
	state.SpaceID = types.StringValue(project.SpaceID)
}
