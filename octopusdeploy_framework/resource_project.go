package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &projectResource{}

type projectResource struct {
	*Config
}

func NewProjectResource() resource.Resource {
	return &projectResource{}
}

func (r *projectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProjectResourceName)
}

func (r *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetProjectResourceSchema()
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
	createdProject, err := projects.Add(r.Client, project)
	if err != nil {
		resp.Diagnostics.AddError("Error creating project", err.Error())
		return
	}

	if project.PersistenceSettings != nil && project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		vcsProject, err := projects.ConvertToVCS(r.Client, createdProject, "Converting project to use VCS", "", project.PersistenceSettings.(projects.GitPersistenceSettings))
		if err != nil {
			resp.Diagnostics.AddError("Error converting project to VCS", err.Error())
			_ = projects.DeleteByID(r.Client, plan.SpaceID.ValueString(), createdProject.GetID())
			return
		}
		createdProject.PersistenceSettings = vcsProject.PersistenceSettings
	}

	createdProject, err = projects.GetByID(r.Client, plan.SpaceID.ValueString(), createdProject.GetID())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving created project", err.Error())
		return
	}

	flattenedProject, diags := flattenProject(ctx, createdProject, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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

	project, err := projects.GetByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	flattenedProject, diags := flattenProject(ctx, project, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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

	updatedProject, err = projects.Update(r.Client, updatedProject)
	if err != nil {
		resp.Diagnostics.AddError("Error updating project", err.Error())
		return
	}

	flattenedProject, diags := flattenProject(ctx, updatedProject, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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
