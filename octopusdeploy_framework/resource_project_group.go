package octopusdeploy_framework

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projectgroups"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type projectGroupTypeResource struct {
	*Config
}

func NewProjectGroupResource() resource.Resource {
	return &projectGroupTypeResource{}
}

func (r *projectGroupTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("project_group")
}

func (r *projectGroupTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: schemas.GetProjectGroupResourceSchema(),
	}
}

func (r *projectGroupTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *projectGroupTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.ProjectGroupTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newGroup := projectgroups.ProjectGroup{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		SpaceID:     data.SpaceID.ValueString(),
	}

	group, err := projectgroups.Add(r.Config.Client, &newGroup)
	if err != nil {
		resp.Diagnostics.AddError("unable to create project group", err.Error())
		return
	}

	updateProjectGroup(&data, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectGroupTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schemas.ProjectGroupTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := projectgroups.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "project group"); err != nil {
			resp.Diagnostics.AddError("unable to load project group", err.Error())
		}
		return
	}

	updateProjectGroup(&data, group)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectGroupTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state schemas.ProjectGroupTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating project group '%s'", data.ID.ValueString()))

	group, err := projectgroups.GetByID(r.Config.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load project group", err.Error())
		return
	}

	group.Name = data.Name.ValueString()
	group.Description = data.Description.ValueString()
	group.SpaceID = data.SpaceID.ValueString()

	updatedProjectGroup, err := projectgroups.Update(r.Config.Client, *group)
	if err != nil {
		resp.Diagnostics.AddError("unable to update project group", err.Error())
		return
	}

	updateProjectGroup(&data, updatedProjectGroup)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectGroupTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.ProjectGroupTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := projectgroups.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete project group", err.Error())
		return
	}
}

func updateProjectGroup(data *schemas.ProjectGroupTypeResourceModel, group *projectgroups.ProjectGroup) {
	data.ID = types.StringValue(group.ID)
	data.Name = types.StringValue(group.Name)
	data.SpaceID = types.StringValue(group.SpaceID)
	data.Description = types.StringValue(group.Description)
}
