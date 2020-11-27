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
	spaceSchema := getSpaceSchema()
	for _, field := range spaceSchema {
		field.Computed = true
		field.Default = nil
		field.MaxItems = 0
		field.MinItems = 0
		field.Optional = false
		field.Required = false
		field.ValidateDiagFunc = nil
		field.ValidateFunc = nil
	}

	return map[string]*schema.Schema{
		"ids": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"partial_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"skip": {
			Default:  0,
			Type:     schema.TypeInt,
			Optional: true,
		},
		"spaces": {
			Computed: true,
			Elem:     &schema.Resource{Schema: spaceSchema},
			Type:     schema.TypeList,
		},
		"take": {
			Default:  1,
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
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
		"name": {
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
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"task_queue_stopped": {
			Optional: true,
			Type:     schema.TypeBool,
		},
	}
}

func setSpace(ctx context.Context, d *schema.ResourceData, space *octopusdeploy.Space) {
	d.Set("description", space.Description)
	d.Set("is_default", space.IsDefault)
	d.Set("name", space.Name)
	d.Set("space_managers_team_members", space.SpaceManagersTeamMembers)
	d.Set("space_managers_teams", space.SpaceManagersTeams)
	d.Set("task_queue_stopped", space.TaskQueueStopped)

	d.SetId(space.GetID())
}
