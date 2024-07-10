package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const description = "space"

func GetSpaceResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"id":          util.GetIdResourceSchema(),
		"description": util.GetDescriptionResourceSchema(description),
		"name":        util.GetNameResourceSchema(true),
		"slug":        util.GetSlugResourceSchema(description),
		"space_managers_teams": resourceSchema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of team IDs designated to be managers of this space.",
			Optional:    true,
			Computed:    true,
		},
		"space_managers_team_members": resourceSchema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of user IDs designated to be managers of this space.",
			Optional:    true,
			Computed:    true,
		},
		"is_task_queue_stopped": resourceSchema.BoolAttribute{
			Description: "Specifies the status of the task queue for this space.",
			Optional:    true,
		},
		"is_default": resourceSchema.BoolAttribute{
			Description: "Specifies if this space is the default space in Octopus.",
			Optional:    true,
		},
	}
}

func GetSpaceDatasourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":          util.GetIdDatasourceSchema(),
		"description": util.GetDescriptionDatasourceSchema(description),
		"name":        util.GetNameDatasourceWithMaxLengthSchema(true, 20),
		"slug":        util.GetSlugDatasourceSchema(description),
		"space_managers_teams": datasourceSchema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of team IDs designated to be managers of this space.",
			Optional:    true,
			Computed:    true,
		},
		"space_managers_team_members": datasourceSchema.ListAttribute{
			ElementType: types.StringType,
			Description: "A list of user IDs designated to be managers of this space.",
			Optional:    true,
			Computed:    true,
		},
		"is_task_queue_stopped": datasourceSchema.BoolAttribute{
			Description: "Specifies the status of the task queue for this space.",
			Optional:    true,
		},
		"is_default": datasourceSchema.BoolAttribute{
			Description: "Specifies if this space is the default space in Octopus.",
			Optional:    true,
		},
	}
}
