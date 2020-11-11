package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandSpace(d *schema.ResourceData) *octopusdeploy.Space {
	name := d.Get("name").(string)

	space := octopusdeploy.NewSpace(name)
	space.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		space.Description = v.(string)
	}

	if v, ok := d.GetOk("is_default"); ok {
		space.IsDefault = v.(bool)
	}

	if v, ok := d.GetOk("space_managers_team_members"); ok {
		space.SpaceManagersTeamMembers = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("space_managers_teams"); ok {
		space.SpaceManagersTeams = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("task_queue_stopped"); ok {
		space.TaskQueueStopped = v.(bool)
	}

	return space
}

func flattenSpace(ctx context.Context, d *schema.ResourceData, space *octopusdeploy.Space) {
	d.Set("description", space.Description)
	d.Set("is_default", space.IsDefault)
	d.Set("name", space.Name)
	d.Set("space_managers_team_members", space.SpaceManagersTeamMembers)
	d.Set("space_managers_teams", space.SpaceManagersTeams)
	d.Set("task_queue_stopped", space.TaskQueueStopped)

	d.SetId(space.GetID())
}

func getSpaceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"is_default": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"name": &schema.Schema{
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"space_managers_team_members": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		"space_managers_teams": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		"task_queue_stopped": {
			Optional: true,
			Type:     schema.TypeBool,
		},
	}
}
