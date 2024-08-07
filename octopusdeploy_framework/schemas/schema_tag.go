package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetTagSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"id": util.GetIdResourceSchema(),
		"canonical_tag_name": resourceSchema.StringAttribute{
			Computed: true,
		},
		"color": resourceSchema.StringAttribute{
			Required: true,
		},
		"description": util.GetDescriptionResourceSchema("tag"),
		"name":        util.GetNameResourceSchema(true),
		"sort_order": resourceSchema.Int64Attribute{
			Computed: true,
			Optional: true,
		},
		"tag_set_id": resourceSchema.StringAttribute{
			Description: "The ID of the associated tag set.",
			Required:    true,
		},
		"tag_set_space_id": resourceSchema.StringAttribute{
			Description: "The Space ID of the associated tag set. Required if the tag set is not in the same space as what is configured on the provider",
			Computed:    true,
			Optional:    true,
		},
	}
}

type TagResourceModel struct {
	ID               types.String `tfsdk:"id"`
	CanonicalTagName types.String `tfsdk:"canonical_tag_name"`
	Color            types.String `tfsdk:"color"`
	Description      types.String `tfsdk:"description"`
	Name             types.String `tfsdk:"name"`
	SortOrder        types.Int64  `tfsdk:"sort_order"`
	TagSetId         types.String `tfsdk:"tag_set_id"`
	TagSetSpaceId    types.String `tfsdk:"tag_set_space_id"`
}

func MapFromStateToTag(data *TagResourceModel) *tagsets.Tag {
	color := data.Color.ValueString()
	name := data.Name.ValueString()

	tag := tagsets.NewTag(name, color)
	tag.ID = data.ID.ValueString()
	tag.CanonicalTagName = data.CanonicalTagName.ValueString()
	tag.Description = data.Description.ValueString()
	tag.SortOrder = int(data.SortOrder.ValueInt64())

	return tag
}

func MapFromTagToState(data *TagResourceModel, tag *tagsets.Tag, tagSet *tagsets.TagSet) {

	data.CanonicalTagName = types.StringValue(tag.CanonicalTagName)
	data.Color = types.StringValue(tag.Color)
	data.Description = types.StringValue(tag.Description)
	data.Name = types.StringValue(tag.Name)
	data.SortOrder = types.Int64Value(int64(tag.SortOrder))

	data.TagSetId = types.StringValue(tagSet.ID)
	data.TagSetSpaceId = types.StringValue(tagSet.SpaceID)

	data.ID = types.StringValue(tag.ID)
}
