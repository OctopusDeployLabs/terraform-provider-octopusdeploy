package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const projectGroupDescription = "project group"

func GetProjectGroupDatasourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":       util.GetIdResourceSchema(),
		"space_id": util.GetSpaceIdResourceSchema(projectGroupDescription),
		"name":     GetReadonlyNameDatasourceSchema(),
		"retention_policy_id": datasourceSchema.StringAttribute{
			Computed:    true,
			Optional:    true,
			Description: "The ID of the retention policy associated with this project group.",
		},
		"description": util.GetDescriptionResourceSchema(projectGroupDescription),
	}
}

func GetProjectGroupResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"id":       util.GetIdResourceSchema(),
		"space_id": util.GetSpaceIdResourceSchema(projectGroupDescription),
		"name":     util.GetNameResourceSchema(true),
		"retention_policy_id": resourceSchema.StringAttribute{
			Computed:    true,
			Optional:    true,
			Description: "The ID of the retention policy associated with this project group.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"description": util.GetDescriptionResourceSchema(projectGroupDescription),
	}
}

type ProjectGroupTypeResourceModel struct {
	Name              types.String `tfsdk:"name"`
	SpaceID           types.String `tfsdk:"space_id"`
	Description       types.String `tfsdk:"description"`
	RetentionPolicyID types.String `tfsdk:"retention_policy_id"`

	ResourceModel
}
