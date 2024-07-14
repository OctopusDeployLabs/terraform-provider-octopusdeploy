package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const description = "project group"

func GetProjectGroupDatasourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":       util.GetIdResourceSchema(),
		"space_id": util.GetSpaceIdResourceSchema(description),
		"name":     util.GetNameResourceSchema(true),
		"retention_policy_id": datasourceSchema.StringAttribute{
			Computed:    true,
			Optional:    true,
			Description: "The ID of the retention policy associated with this project group.",
		},
		"description": util.GetDescriptionResourceSchema(description),
	}
}

func GetProjectGroupResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"id":       util.GetIdResourceSchema(),
		"space_id": util.GetSpaceIdResourceSchema(description),
		"name":     util.GetNameResourceSchema(true),
		"retention_policy_id": resourceSchema.StringAttribute{
			Computed:    true,
			Optional:    true,
			Description: "The ID of the retention policy associated with this project group.",
		},
		"description": util.GetDescriptionResourceSchema(description),
	}
}

type ProjectGroupTypeResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	SpaceID           types.String `tfsdk:"space_id"`
	Description       types.String `tfsdk:"description"`
	RetentionPolicyID types.String `tfsdk:"retention_policy_id"`
}
