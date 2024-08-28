package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const projectGroupDescription = "project group"

type ProjectGroupSchema struct{}

var _ EntitySchema = ProjectGroupSchema{}

func (p ProjectGroupSchema) GetDatasourceSchema() datasourceSchema.Schema {
	description := "project group"
	return datasourceSchema.Schema{
		Attributes: map[string]datasourceSchema.Attribute{
			// request
			"space_id":     GetSpaceIdDatasourceSchema(description, false),
			"ids":          GetQueryIDsDatasourceSchema(),
			"partial_name": GetQueryPartialNameDatasourceSchema(),
			"skip":         GetQuerySkipDatasourceSchema(),
			"take":         GetQueryTakeDatasourceSchema(),

			// response
			"id": GetIdDatasourceSchema(true),
			"project_groups": datasourceSchema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of project groups that match the filter(s).",
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"id":          GetIdDatasourceSchema(true),
						"space_id":    GetSpaceIdDatasourceSchema(description, true),
						"name":        GetReadonlyNameDatasourceSchema(),
						"description": GetDescriptionDatasourceSchema(projectGroupDescription),
					},
				},
			},
		},
	}
}

func (p ProjectGroupSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id":          GetIdResourceSchema(),
			"space_id":    GetSpaceIdResourceSchema(projectGroupDescription),
			"name":        GetNameResourceSchema(true),
			"description": GetDescriptionResourceSchema(projectGroupDescription),
		},
	}
}

type ProjectGroupTypeResourceModel struct {
	Name        types.String `tfsdk:"name"`
	SpaceID     types.String `tfsdk:"space_id"`
	Description types.String `tfsdk:"description"`

	ResourceModel
}
