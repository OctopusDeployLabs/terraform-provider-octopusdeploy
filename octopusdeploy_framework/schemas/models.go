package schemas

import "github.com/hashicorp/terraform-plugin-framework/types"

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
