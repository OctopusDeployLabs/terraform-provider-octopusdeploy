package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const projectGroupDescription = "project group"

func GetProjectGroupDatasourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":          GetIdDatasourceSchema(true),
		"space_id":    GetSpaceIdDatasourceSchema(projectGroupDescription, true),
		"name":        GetReadonlyNameDatasourceSchema(),
		"description": GetReadonlyDescriptionDatasourceSchema(projectGroupDescription),
	}
}

func GetProjectGroupResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"id":          util.GetIdResourceSchema(),
		"space_id":    util.GetSpaceIdResourceSchema(projectGroupDescription),
		"name":        util.GetNameResourceSchema(true),
		"description": util.GetDescriptionResourceSchema(projectGroupDescription),
	}
}

type ProjectGroupTypeResourceModel struct {
	Name        types.String `tfsdk:"name"`
	SpaceID     types.String `tfsdk:"space_id"`
	Description types.String `tfsdk:"description"`

	ResourceModel
}
