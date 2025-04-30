package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

var _ resource.Resource = &projectResource{}
var _ resource.ResourceWithImportState = &projectResource{}

type projectResource struct {
	*Config
}

func NewProjectResource() resource.Resource {
	return &projectResource{}
}

func (r *projectResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProjectResourceName)
}

func (r *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProjectSchema{}.GetResourceSchema()
}

func (r *projectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan projectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project := expandProject(ctx, plan)
	// PersistenceSettings.Password doesn't return from API so this is work around
	persistenceSettings := project.PersistenceSettings
	createdProject, err := projects.Add(r.Client, project)
	if err != nil {
		resp.Diagnostics.AddError("Error creating project", err.Error())
		return
	}

	if persistenceSettings != nil && persistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		_, err := projects.ConvertToVCS(r.Client, createdProject, "Converting project to use VCS", "", persistenceSettings.(projects.GitPersistenceSettings))
		if err != nil {
			resp.Diagnostics.AddError("Error converting project to VCS", err.Error())
			_ = projects.DeleteByID(r.Client, plan.SpaceID.ValueString(), createdProject.GetID())
			return
		}
	}

	createdProject, err = projects.GetByID(r.Client, plan.SpaceID.ValueString(), createdProject.GetID())
	if persistenceSettings != nil {
		createdProject.PersistenceSettings = persistenceSettings
	}

	flattenedProject, diags := flattenProject(ctx, createdProject, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deploymentDiags := r.updateStateWithDeploymentSettings(createdProject, flattenedProject, &plan)
	if deploymentDiags.HasError() {
		return
	}

	diags = resp.State.Set(ctx, flattenedProject)
	resp.Diagnostics.Append(diags...)
}

func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state projectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	stateProject := expandProject(ctx, state)
	// PersistenceSettings.Password doesn't return from API so this is work around
	persistenceSettings := stateProject.PersistenceSettings

	project, err := projects.GetByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "project"); err != nil {
			resp.Diagnostics.AddError("Error reading project", err.Error())
		}
		return
	}
	if persistenceSettings != nil {
		project.PersistenceSettings = persistenceSettings
	}

	flattenedProject, diags := flattenProject(ctx, project, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diagFromUpdate := r.updateStateWithDeploymentSettings(project, flattenedProject, &state)
	if diagFromUpdate.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, flattenedProject)...)
}

func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan projectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	existingProject, err := projects.GetByID(r.Client, plan.SpaceID.ValueString(), plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving project", err.Error())
		return
	}

	updatedProject := expandProject(ctx, plan)
	updatedProject.ID = existingProject.ID
	updatedProject.Links = existingProject.Links
	// PersistenceSettings.Password doesn't return from API so this is work around
	persistenceSettings := updatedProject.PersistenceSettings

	if updatedProject.PersistenceSettings != nil && updatedProject.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		if existingProject.PersistenceSettings == nil || existingProject.PersistenceSettings.Type() != projects.PersistenceSettingsTypeVersionControlled {
			vcsProject, err := projects.ConvertToVCS(r.Client, existingProject, "Converting project to use VCS", "", updatedProject.PersistenceSettings.(projects.GitPersistenceSettings))
			if err != nil {
				resp.Diagnostics.AddError("Error converting project to VCS", err.Error())
				return
			}
			updatedProject.PersistenceSettings = vcsProject.PersistenceSettings
		}
	}

	if updatedProject.AutoCreateRelease == true && updatedProject.ReleaseCreationStrategy == nil {
		// This condition is possible when 'built_in_trigger' resource is used to maintain release creation strategy
		// For this scenario we want to send persisted strategy to the API to avoid an error(missing package for ARC) which practitioner will not be able to escape
		updatedProject.ReleaseCreationStrategy = existingProject.ReleaseCreationStrategy
	}

	updatedProject, err = projects.Update(r.Client, updatedProject)
	if err != nil {
		resp.Diagnostics.AddError("Error updating project", err.Error())
		return
	}

	if persistenceSettings != nil {
		updatedProject.PersistenceSettings = persistenceSettings
	}

	flattenedProject, diags := flattenProject(ctx, updatedProject, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diagFromUpdate := r.updateStateWithDeploymentSettings(updatedProject, flattenedProject, &plan)
	if diagFromUpdate.HasError() {
		return
	}

	diags = resp.State.Set(ctx, flattenedProject)
	resp.Diagnostics.Append(diags...)
}

func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state projectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := projects.DeleteByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting project", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (*projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *projectResource) updateStateWithDeploymentSettings(project *projects.Project, newState *projectResourceModel, originalState *projectResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	var gitRef string
	if project.IsVersionControlled {
		if gitSettings, ok := project.PersistenceSettings.(projects.GitPersistenceSettings); ok {
			gitRef = gitSettings.DefaultBranch()
		}
		if gitRef == "" {
			gitRef = "main"
		}
	}

	deploymentSettings, err := r.Client.Deployments.GetDeploymentSettings(project, gitRef)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading deployment settings: %w", err))
	}

	// Update the state with the deployment settings
	if deploymentSettings.ConnectivityPolicy != nil && !originalState.ConnectivityPolicy.IsNull() {
		newState.ConnectivityPolicy = flattenConnectivityPolicy(deploymentSettings.ConnectivityPolicy)
	}
	newState.DefaultGuidedFailureMode = types.StringValue(string(deploymentSettings.DefaultGuidedFailureMode))
	newState.DefaultToSkipIfAlreadyInstalled = types.BoolValue(deploymentSettings.DefaultToSkipIfAlreadyInstalled)
	newState.DeploymentChangesTemplate = types.StringValue(deploymentSettings.DeploymentChangesTemplate)
	newState.ReleaseNotesTemplate = types.StringValue(deploymentSettings.ReleaseNotesTemplate)
	if deploymentSettings.VersioningStrategy != nil && !originalState.VersioningStrategy.IsNull() {
		newState.VersioningStrategy = flattenVersioningStrategy(deploymentSettings.VersioningStrategy)
	}

	return diags
}
