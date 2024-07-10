package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projectgroups"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type projectGroupTypeResource struct {
	*Config
}

type projectGroupTypeResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	SpaceID           types.String `tfsdk:"space_id"`
	Description       types.String `tfsdk:"description"`
	RetentionPolicyID types.String `tfsdk:"retention_policy_id"`
}

func NewProjectGroupResource() resource.Resource {
	return &projectGroupTypeResource{}
}

func (r *projectGroupTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "octopusdeploy_project_group"
}

func (r *projectGroupTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	description := "project group"
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":       util.GetIdResourceSchema(),
			"space_id": util.GetSpaceIdResourceSchema(description),
			"name":     util.GetNameResourceSchema(true),
			"retention_policy_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The ID of the retention policy associated with this project group.",
			},
			"description": util.GetDescriptionResourceSchema(description),
		},
	}
}

func (r *projectGroupTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *projectGroupTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *projectGroupTypeResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newGroup := projectgroups.ProjectGroup{
		Name:              data.Name.ValueString(),
		Description:       data.Description.ValueString(),
		RetentionPolicyID: data.RetentionPolicyID.ValueString(),
		SpaceID:           data.SpaceID.ValueString(),
	}

	group, err := projectgroups.Add(r.Config.Client, &newGroup)
	if err != nil {
		resp.Diagnostics.AddError("unable to create project group", err.Error())
		return
	}

	data.ID = types.StringValue(group.ID)
	data.Name = types.StringValue(group.Name)
	data.SpaceID = types.StringValue(group.SpaceID)
	data.RetentionPolicyID = types.StringValue(group.RetentionPolicyID)
	data.Description = types.StringValue(group.Description)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *projectGroupTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *projectGroupTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := projectgroups.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load project group", err.Error())
		return
	}

	data.ID = types.StringValue(group.ID)
	data.Name = types.StringValue(group.Name)
	data.SpaceID = types.StringValue(group.SpaceID)
	data.RetentionPolicyID = types.StringValue(group.RetentionPolicyID)
	data.Description = types.StringValue(group.Description)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectGroupTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *projectGroupTypeResourceModel

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
	group.RetentionPolicyID = data.RetentionPolicyID.ValueString()
	group.SpaceID = data.SpaceID.ValueString()

	updatedProjectGroup, err := projectgroups.Update(r.Config.Client, *group)
	if err != nil {
		resp.Diagnostics.AddError("unable to update project group", err.Error())
		return
	}

	data.ID = types.StringValue(updatedProjectGroup.ID)
	data.Name = types.StringValue(updatedProjectGroup.Name)
	data.SpaceID = types.StringValue(updatedProjectGroup.SpaceID)
	data.RetentionPolicyID = types.StringValue(updatedProjectGroup.RetentionPolicyID)
	data.Description = types.StringValue(updatedProjectGroup.Description)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectGroupTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *projectGroupTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := projectgroups.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete project group", err.Error())
		return
	}
}

func resourceConfiguration(req resource.ConfigureRequest, resp *resource.ConfigureResponse) *Config {
	if req.ProviderData == nil {
		return nil
	}

	p, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return p
}
