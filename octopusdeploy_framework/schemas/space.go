package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const spaceDescription = "space"

type SpaceModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Slug                     types.String `tfsdk:"slug"`
	Description              types.String `tfsdk:"description"`
	IsDefault                types.Bool   `tfsdk:"is_default"`
	SpaceManagersTeams       types.List   `tfsdk:"space_managers_teams"`
	SpaceManagersTeamMembers types.List   `tfsdk:"space_managers_team_members"`
	IsTaskQueueStopped       types.Bool   `tfsdk:"is_task_queue_stopped"`
}

func GetSpaceResourceSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		"id":          GetIdResourceSchema(),
		"description": GetDescriptionResourceSchema(spaceDescription),
		"name":        GetNameResourceSchema(true),
		"slug":        GetSlugResourceSchema(spaceDescription),
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
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		},
		"is_default": resourceSchema.BoolAttribute{
			Description: "Specifies if this space is the default space in Octopus.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		},
	}
}

func GetSpaceDatasourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":          GetIdDatasourceSchema(),
		"description": GetDescriptionDatasourceSchema(spaceDescription),
		"name":        GetNameDatasourceWithMaxLengthSchema(true, 20),
		"slug":        GetSlugDatasourceSchema(spaceDescription),
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
