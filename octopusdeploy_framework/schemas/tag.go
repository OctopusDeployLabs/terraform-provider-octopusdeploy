package schemas

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const TagResourceName = "tag"

type TagSchema struct{}

var _ EntitySchema = TagSchema{}

func (t TagSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

func (t TagSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages tags in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id": GetIdResourceSchema(),
			"canonical_tag_name": util.ResourceString().
				Computed().
				Description("The canonical name of the tag.").
				Build(),
			"color": util.ResourceString().
				Required().
				Description("The color of the tag.").
				Build(),
			"description": util.ResourceString().
				Optional().
				Description("The description of the tag.").
				Default("").
				Computed().
				Build(),
			"name": util.ResourceString().
				Required().
				Description("The name of the tag.").
				Build(),
			"sort_order": util.ResourceInt64().
				Optional().
				Computed().
				Description("The sort order of the tag.").
				Build(),
			"tag_set_id": util.ResourceString().
				Required().
				Description("The ID of the associated tag set.").
				Build(),
			"tag_set_space_id": util.ResourceString().
				Optional().
				Computed().
				Description("The Space ID of the associated tag set. Required if the tag set is not in the same space as what is configured on the provider.").
				Build(),
		},
	}
}

type TagResourceModel struct {
	CanonicalTagName types.String `tfsdk:"canonical_tag_name"`
	Color            types.String `tfsdk:"color"`
	Description      types.String `tfsdk:"description"`
	Name             types.String `tfsdk:"name"`
	SortOrder        types.Int64  `tfsdk:"sort_order"`
	TagSetId         types.String `tfsdk:"tag_set_id"`
	TagSetSpaceId    types.String `tfsdk:"tag_set_space_id"`
	ResourceModel
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
