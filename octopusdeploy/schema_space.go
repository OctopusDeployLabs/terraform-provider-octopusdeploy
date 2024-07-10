package octopusdeploy

import (
	"context"
	"fmt"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const spaceManagersTeamIDPrefix = "teams-spacemanagers-"

func expandSpace(d *schema.ResourceData) *spaces.Space {
	name := d.Get("name").(string)

	space := spaces.NewSpace(name)
	space.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		space.Description = v.(string)
	}

	if v, ok := d.GetOk("is_default"); ok {
		space.IsDefault = v.(bool)
	}

	if v, ok := d.GetOk("slug"); ok {
		space.Slug = v.(string)
	}

	if v, ok := d.GetOk("space_managers_team_members"); ok {
		space.SpaceManagersTeamMembers = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("space_managers_teams"); ok {
		space.SpaceManagersTeams = addSpaceManagers(space.GetID(), getSliceFromTerraformTypeList(v))
	}

	if v, ok := d.GetOk("is_task_queue_stopped"); ok {
		space.TaskQueueStopped = v.(bool)
	}

	return space
}

func addSpaceManagers(spaceID string, teamIDs []string) []string {
	var newSlice []string
	if getStringOrEmpty(spaceID) != "" {
		newSlice = append(newSlice, spaceManagersTeamIDPrefix+spaceID)
	}
	newSlice = append(newSlice, teamIDs...)
	return newSlice
}

func getSpaceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": getDescriptionSchema("space"),
		"id":          getIDSchema(),
		"is_default": {
			Description: "Specifies if this space is the default space in Octopus.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"name": getNameSchemaWithMaxLength(true, 20),
		"slug": {
			Computed:    true,
			Description: "The unique slug of this space.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"space_managers_team_members": {
			Computed:    true,
			Description: "A list of user IDs designated to be managers of this space.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeSet,
		},
		"space_managers_teams": {
			Computed:    true,
			Description: "A list of team IDs designated to be managers of this space.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeSet,
		},
		"is_task_queue_stopped": {
			Description: "Specifies the status of the task queue for this space.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
	}
}

func setSpace(ctx context.Context, d *schema.ResourceData, space *spaces.Space) error {
	d.Set("description", space.Description)
	d.Set("id", space.GetID())
	d.Set("is_default", space.IsDefault)
	d.Set("name", space.Name)
	d.Set("slug", space.Slug)

	if err := d.Set("space_managers_team_members", space.SpaceManagersTeamMembers); err != nil {
		return fmt.Errorf("error setting space_managers_team_members: %s", err)
	}

	if err := d.Set("space_managers_teams", removeSpaceManagers(space.SpaceManagersTeams)); err != nil {
		return fmt.Errorf("error setting space_managers_teams: %s", err)
	}

	d.Set("is_task_queue_stopped", space.TaskQueueStopped)

	return nil
}

func removeSpaceManagers(teamIDs []string) []string {
	if len(teamIDs) == 0 {
		return teamIDs
	}
	var newSlice []string
	for _, v := range teamIDs {
		if !strings.Contains(v, spaceManagersTeamIDPrefix) {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}
