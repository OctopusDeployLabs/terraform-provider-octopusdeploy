package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func flattenSpace(space *octopusdeploy.Space) map[string]interface{} {
	if space == nil {
		return nil
	}

	return map[string]interface{}{
		"description":                 space.Description,
		"id":                          space.GetID(),
		"is_default":                  space.IsDefault,
		"name":                        space.Name,
		"space_managers_team_members": space.SpaceManagersTeamMembers,
		"space_managers_teams":        space.SpaceManagersTeams,
		"task_queue_stopped":          space.TaskQueueStopped,
	}
}

func getSpaceDataSchema() map[string]*schema.Schema {
	dataSchema := getSpaceSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"id":           getIDDataSchema(),
		"ids":          getIDsQuery(),
		"name":         getNameQuery(),
		"partial_name": getPartialNameQuery(),
		"skip":         getSkipQuery(),
		"spaces": {
			Computed:    true,
			Description: "A list of spaces that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"take": getTakeQuery(),
	}
}

func getSpaceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": getDescriptionSchema(),
		"id":          getIDSchema(),
		"is_default": {
			Description: "Specifies if this space is the default space in Octopus.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"name": getNameSchema(true),
		"space_managers_team_members": {
			Computed:    true,
			Description: "A list of user IDs designated to be managers of this space.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"space_managers_teams": {
			Computed:    true,
			Description: "A list of team IDs designated to be managers of this space.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"task_queue_stopped": {
			Description: "Specifies the status of the task queue for this space.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
	}
}

func setSpace(ctx context.Context, d *schema.ResourceData, space *octopusdeploy.Space) {
	d.Set("description", space.Description)
	d.Set("id", space.GetID())
	d.Set("is_default", space.IsDefault)
	d.Set("name", space.Name)
	d.Set("space_managers_team_members", space.SpaceManagersTeamMembers)
	d.Set("space_managers_teams", space.SpaceManagersTeams)
	d.Set("task_queue_stopped", space.TaskQueueStopped)
}
