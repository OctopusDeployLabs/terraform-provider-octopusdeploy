package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const projectGroupDescription = "project group"

type 

func GetProjectGroupDatasourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":          GetIdResourceSchema(),
		"space_id":    GetSpaceIdResourceSchema(projectGroupDescription),
		"name":        GetReadonlyNameDatasourceSchema(),
		"description": GetDescriptionResourceSchema(projectGroupDescription),
	}
}

func GetProjectGroupResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"id":          GetIdResourceSchema(),
		"space_id":    GetSpaceIdResourceSchema(projectGroupDescription),
		"name":        GetNameResourceSchema(true),
		"description": GetDescriptionResourceSchema(projectGroupDescription),
	}
}

type ProjectGroupTypeResourceModel struct {
	Name        types.String `tfsdk:"name"`
	SpaceID     types.String `tfsdk:"space_id"`
	Description types.String `tfsdk:"description"`

	ResourceModel
}
