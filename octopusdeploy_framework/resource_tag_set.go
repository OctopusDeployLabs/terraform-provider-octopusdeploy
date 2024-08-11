package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &tagSetResource{}
var _ resource.ResourceWithImportState = &tagSetResource{}

type tagSetResource struct {
	*Config
}

func NewTagSetResource() resource.Resource {
	return &tagSetResource{}
}

func (r *tagSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.TagSetResourceName)
}

func (r *tagSetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetTagSetResourceSchema()
}

func (r *tagSetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *tagSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.TagSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagSet := expandTagSet(plan)
	createdTagSet, err := tagsets.Add(r.Client, tagSet)
	if err != nil {
		resp.Diagnostics.AddError("Error creating tag set", err.Error())
		return
	}

	state := flattenTagSet(createdTagSet)
	state.ID = types.StringValue(createdTagSet.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *tagSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.TagSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagSet, err := tagsets.GetByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "tagSetResource"); err != nil {
			resp.Diagnostics.AddError("unable to load tag set", err.Error())
		}
		return
	}

	newState := flattenTagSet(tagSet)
	state.ID = types.StringValue(tagSet.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *tagSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.TagSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagSet := expandTagSet(plan)
	updatedTagSet, err := tagsets.Update(r.Client, tagSet)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tag set", err.Error())
		return
	}

	state := flattenTagSet(updatedTagSet)
	state.ID = types.StringValue(updatedTagSet.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *tagSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.TagSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := tagsets.DeleteByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting tag set", err.Error())
		return
	}
}

func (r *tagSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tagSetID := req.ID

	tagSet, err := tagsets.GetByID(r.Client, r.Client.GetSpaceID(), tagSetID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tag set",
			fmt.Sprintf("Unable to read tag set with ID %s: %s", tagSetID, err.Error()),
		)
		return
	}

	state := flattenTagSet(tagSet)
	state.ID = types.StringValue(tagSet.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func expandTagSet(model schemas.TagSetResourceModel) *tagsets.TagSet {
	tagSet := tagsets.NewTagSet(model.Name.ValueString())
	tagSet.ID = model.ID.ValueString()
	tagSet.Description = model.Description.ValueString()
	tagSet.SpaceID = model.SpaceID.ValueString()

	if !model.SortOrder.IsNull() {
		tagSet.SortOrder = int32(model.SortOrder.ValueInt64())
	}

	return tagSet
}

func flattenTagSet(tagSet *tagsets.TagSet) schemas.TagSetResourceModel {
	return schemas.TagSetResourceModel{
		Name:        types.StringValue(tagSet.Name),
		Description: types.StringValue(tagSet.Description),
		SortOrder:   types.Int64Value(int64(tagSet.SortOrder)),
		SpaceID:     types.StringValue(tagSet.SpaceID),
	}
}
