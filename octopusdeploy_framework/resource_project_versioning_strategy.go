package octopusdeploy_framework

import (
	"context"
	"log"
	"net/http"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &projectVersioningStrategyResource{}

type projectVersioningStrategyResource struct {
	*Config
}

func NewProjectVersioningStrategyResource() resource.Resource {
	return &projectVersioningStrategyResource{}
}

func (r *projectVersioningStrategyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProjectVersioningStrategyResourceName)
}

func (r *projectVersioningStrategyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProjectVersioningStrategySchema{}.GetResourceSchema()
}

func (r *projectVersioningStrategyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *projectVersioningStrategyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.ProjectVersioningStrategyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := projects.GetByID(r.Client, plan.SpaceID.ValueString(), plan.ProjectID.ValueString())
	if err != nil {
		if apiError, ok := err.(*core.APIError); ok {
			if apiError.StatusCode == http.StatusNotFound {
				log.Printf("[INFO] associated project (%s) not found; deleting version strategy from state", plan.ProjectID.ValueString())
				resp.State.RemoveResource(ctx)
			}
		} else {
			resp.Diagnostics.AddError("Failed to read associated project", err.Error())
		}
		return
	}
	versioningStrategy := mapStateToProjectVersioningStrategy(&plan)
	project.VersioningStrategy = versioningStrategy

	_, err = projects.Update(r.Client, project)
	if err != nil {
		resp.Diagnostics.AddError("Error updating associated project", err.Error())
		return
	}

	updatedProject, err := projects.GetByID(r.Client, plan.SpaceID.ValueString(), plan.ProjectID.ValueString())
	if err != nil {
		if apiError, ok := err.(*core.APIError); ok {
			if apiError.StatusCode == http.StatusNotFound {
				log.Printf("[INFO] associated project (%s) not found; deleting version strategy from state", plan.ProjectID.ValueString())
				resp.State.RemoveResource(ctx)
			}
		} else {
			resp.Diagnostics.AddError("Failed to read associated project", err.Error())
		}
		return
	}

	mapProjectVersioningStrategyToState(updatedProject.VersioningStrategy, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *projectVersioningStrategyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.ProjectVersioningStrategyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := projects.GetByID(r.Client, state.SpaceID.ValueString(), state.ProjectID.ValueString())
	if err != nil {
		if apiError, ok := err.(*core.APIError); ok {
			if apiError.StatusCode == http.StatusNotFound {
				log.Printf("[INFO] associated project (%s) not found; deleting version strategy from state", state.ProjectID.ValueString())
				resp.State.RemoveResource(ctx)
			}
		} else {
			resp.Diagnostics.AddError("Failed to read associated project", err.Error())
		}
		return
	}
	mapProjectVersioningStrategyToState(project.VersioningStrategy, &state)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *projectVersioningStrategyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.ProjectVersioningStrategyModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	existingProject, err := projects.GetByID(r.Client, plan.SpaceID.ValueString(), plan.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving associated project", err.Error())
		return
	}

	versioningStrategy := mapStateToProjectVersioningStrategy(&plan)
	existingProject.VersioningStrategy = versioningStrategy

	_, err = projects.Update(r.Client, existingProject)
	if err != nil {
		resp.Diagnostics.AddError("Error updating associated project", err.Error())
		return
	}

	updatedProject, err := projects.GetByID(r.Client, plan.SpaceID.ValueString(), plan.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving associated project", err.Error())
		return
	}

	mapProjectVersioningStrategyToState(updatedProject.VersioningStrategy, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *projectVersioningStrategyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.ProjectVersioningStrategyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := projects.GetByID(r.Client, state.SpaceID.ValueString(), state.ProjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving project", err.Error())
		return
	}

	project.VersioningStrategy = &projects.VersioningStrategy{}
	_, err = projects.Update(r.Client, project)
	if err != nil {
		resp.Diagnostics.AddError("Error updating project to remove versioning strategy", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func mapStateToProjectVersioningStrategy(state *schemas.ProjectVersioningStrategyModel) *projects.VersioningStrategy {
	var donorPackageStepID *string
	donorPackageStepIDString := state.DonorPackageStepID.ValueString()
	if donorPackageStepIDString != "" {
		donorPackageStepID = &donorPackageStepIDString
	}

	return &projects.VersioningStrategy{
		Template:           state.Template.ValueString(),
		DonorPackageStepID: donorPackageStepID,
		DonorPackage: &packages.DeploymentActionPackage{
			DeploymentAction: state.DonorPackage.DeploymentAction.ValueString(),
			PackageReference: state.DonorPackage.PackageReference.ValueString(),
		},
	}
}

func mapProjectVersioningStrategyToState(versioningStrategy *projects.VersioningStrategy, state *schemas.ProjectVersioningStrategyModel) {
	if versioningStrategy.DonorPackageStepID != nil {
		state.DonorPackageStepID = types.StringValue(*versioningStrategy.DonorPackageStepID)
	}
	// Template and Donor Package are mutually exclusive options. We won't always have DonorPackage information.
	state.Template = types.StringValue(versioningStrategy.Template)

	if !(versioningStrategy.DonorPackage == nil) {
		state.DonorPackage.PackageReference = types.StringValue(versioningStrategy.DonorPackage.PackageReference)
		state.DonorPackage.DeploymentAction = types.StringValue(versioningStrategy.DonorPackage.DeploymentAction)
	}
}
