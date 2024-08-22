package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const TagSetDataSourceName = "tag_sets"
const TagSetResourceName = "tag_set"

func GetTagSetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages tag sets in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id": util.ResourceString().
				Optional().
				Computed().
				Description("The unique ID for this resource.").
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
			"name": util.ResourceString().
				Required().
				Description("The name of this resource.").
				Build(),
			"description": util.ResourceString().
				Optional().
				Computed().
				Description("The description of this tag set.").
				Build(),
			"sort_order": util.ResourceInt64().
				Optional().
				Computed().
				Description("The sort order associated with this resource.").
				Build(),
			"space_id": util.ResourceString().
				Optional().
				Computed().
				Description("The space ID associated with this resource.").
				PlanModifiers(stringplanmodifier.UseStateForUnknown()).
				Build(),
		},
	}
}

func GetTagSetDataSourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: "Provides information about existing tag sets.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id": util.DataSourceString().
				Computed().
				Description("The ID of this resource.").
				Build(),
			"space_id": util.DataSourceString().
				Optional().
				Description("The space ID associated with this resource.").
				Build(),
			"ids": util.DataSourceList(types.StringType).
				Optional().
				Description("A filter to search by a list of IDs.").
				Build(),
			"partial_name": util.DataSourceString().
				Optional().
				Description("A filter to search by the partial match of a name.").
				Build(),
			"skip": util.DataSourceInt64().
				Optional().
				Description("A filter to specify the number of items to skip in the response.").
				Build(),
			"take": util.DataSourceInt64().
				Optional().
				Description("A filter to specify the number of items to take (or return) in the response.").
				Build(),
			"tag_sets": datasourceSchema.ListNestedAttribute{
				Computed:    true,
				Optional:    false,
				Description: "A list of tag sets that match the filter(s).",
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"id": util.DataSourceString().
							Optional().
							Computed().
							Description("The unique ID for this resource.").
							Build(),
						"name": util.DataSourceString().
							Optional().
							Computed().
							Description("The name of this resource.").
							Build(),
						"description": util.DataSourceString().
							Optional().
							Computed().
							Description("The description of this tag set.").
							Build(),
						"sort_order": util.DataSourceInt64().
							Optional().
							Computed().
							Description("The sort order associated with this resource.").
							Build(),
						"space_id": util.DataSourceString().
							Optional().
							Computed().
							Description("The space ID associated with this resource.").
							Build(),
					},
				},
			},
		},
	}
}

func GetTagSetAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          types.StringType,
		"space_id":    types.StringType,
		"name":        types.StringType,
		"description": types.StringType,
		"sort_order":  types.Int64Type,
	}
}

type TagSetDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	SpaceID     types.String `tfsdk:"space_id"`
	IDs         types.List   `tfsdk:"ids"`
	PartialName types.String `tfsdk:"partial_name"`
	Skip        types.Int64  `tfsdk:"skip"`
	Take        types.Int64  `tfsdk:"take"`
	TagSets     types.List   `tfsdk:"tag_sets"`
}

type TagSetResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	SortOrder   types.Int64  `tfsdk:"sort_order"`
	SpaceID     types.String `tfsdk:"space_id"`

	ResourceModel
}
