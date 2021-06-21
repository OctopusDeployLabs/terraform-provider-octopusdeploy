package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandTeam(d *schema.ResourceData) *octopusdeploy.Team {
	name := d.Get("name").(string)

	team := octopusdeploy.NewTeam(name)
	team.ID = d.Id()

	if v, ok := d.GetOk("can_be_deleted"); ok {
		team.CanBeDeleted = v.(bool)
	}

	if v, ok := d.GetOk("can_be_renamed"); ok {
		team.CanBeRenamed = v.(bool)
	}

	if v, ok := d.GetOk("can_change_members"); ok {
		team.CanChangeMembers = v.(bool)
	}

	if v, ok := d.GetOk("can_change_roles"); ok {
		team.CanChangeRoles = v.(bool)
	}

	if v, ok := d.GetOk("description"); ok {
		team.Description = v.(string)
	}

	if v, ok := d.GetOk("external_security_group"); ok {
		team.ExternalSecurityGroups = expandExternalSecurityGroups(v.([]interface{}))
	}

	if v, ok := d.GetOk("space_id"); ok {
		team.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("users"); ok {
		team.MemberUserIDs = getSliceFromTerraformTypeList(v)
	}

	return team
}

func flattenTeam(team *octopusdeploy.Team) map[string]interface{} {
	if team == nil {
		return nil
	}

	return map[string]interface{}{
		"can_be_deleted":          team.CanBeDeleted,
		"can_be_renamed":          team.CanBeRenamed,
		"can_change_members":      team.CanChangeMembers,
		"can_change_roles":        team.CanChangeRoles,
		"description":             team.Description,
		"external_security_group": flattenExternalSecurityGroups(team.ExternalSecurityGroups),
		"id":                      team.GetID(),
		"name":                    team.Name,
		"space_id":                team.SpaceID,
		"users":                   team.MemberUserIDs,
	}
}

func getTeamDataSchema() map[string]*schema.Schema {
	dataSchema := getTeamSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"id":             getDataSchemaID(),
		"ids":            getQueryIDs(),
		"include_system": getQueryIncludeSystem(),
		"partial_name":   getQueryPartialName(),
		"skip":           getQuerySkip(),
		"spaces": {
			Description: "A list of spaces that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"take": getQueryTake(),
		"teams": {
			Computed:    true,
			Description: "A list of teams that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}

func getTeamSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"can_be_deleted": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"can_be_renamed": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"can_change_members": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"can_change_roles": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"description": {
			Description: "The user-friendly description of this team.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"external_security_group": {
			Optional: true,
			Elem:     &schema.Resource{Schema: getExternalSecurityGroupsSchema()},
			Type:     schema.TypeList,
		},
		"id": getIDSchema(),
		"name": {
			Description: "The name of this team.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"space_id": {
			Computed:    true,
			Description: "The space associated with this team.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"users": {
			Computed:    true,
			Description: "A list of user IDs designated to be members of this team.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}

func setTeam(ctx context.Context, d *schema.ResourceData, team *octopusdeploy.Team) error {
	d.Set("can_be_deleted", team.CanBeDeleted)
	d.Set("can_be_renamed", team.CanBeRenamed)
	d.Set("can_change_members", team.CanChangeMembers)
	d.Set("can_change_roles", team.CanChangeRoles)
	d.Set("description", team.Description)

	if err := d.Set("external_security_group", flattenExternalSecurityGroups(team.ExternalSecurityGroups)); err != nil {
		return fmt.Errorf("error setting external_security_group: %s", err)
	}

	d.Set("id", team.GetID())
	d.Set("name", team.Name)
	d.Set("space_id", team.SpaceID)

	if err := d.Set("users", team.MemberUserIDs); err != nil {
		return fmt.Errorf("error setting users: %s", err)
	}

	d.SetId(team.GetID())

	return nil
}
