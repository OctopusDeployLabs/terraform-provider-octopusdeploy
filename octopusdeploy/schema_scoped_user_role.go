package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/userroles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandScopedUserRole(d *schema.ResourceData) *userroles.ScopedUserRole {
	var scopedUserRole = userroles.NewScopedUserRole(d.Get("user_role_id").(string))
	scopedUserRole.ID = d.Get("id").(string)

	scopedUserRole.TeamID = d.Get("team_id").(string)
	scopedUserRole.SpaceID = d.Get("space_id").(string)
	if v, ok := d.GetOk("environment_ids"); ok {
		scopedUserRole.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}
	if v, ok := d.GetOk("project_ids"); ok {
		scopedUserRole.ProjectIDs = getSliceFromTerraformTypeList(v)

	}
	if v, ok := d.GetOk("project_group_ids"); ok {
		scopedUserRole.ProjectGroupIDs = getSliceFromTerraformTypeList(v)

	}
	if v, ok := d.GetOk("tenant_ids"); ok {
		scopedUserRole.TenantIDs = getSliceFromTerraformTypeList(v)
	}
	return scopedUserRole
}
func flattenScopedUserRoles(scopedUserRoles []*userroles.ScopedUserRole) []interface{} {
	if scopedUserRoles == nil {
		return nil
	}
	var flattenedScopedUserRoles = make([]interface{}, len(scopedUserRoles))
	for key, scopedUserRole := range scopedUserRoles {
		flattenedScopedUserRoles[key] = flattenScopedUserRole(scopedUserRole)
	}

	return flattenedScopedUserRoles
}

func flattenScopedUserRole(scopedUserRole *userroles.ScopedUserRole) map[string]interface{} {
	if scopedUserRole == nil {
		return nil
	}
	return map[string]interface{}{
		"environment_ids":   schema.NewSet(schema.HashString, flattenArray(scopedUserRole.EnvironmentIDs)),
		"id":                scopedUserRole.ID,
		"project_group_ids": schema.NewSet(schema.HashString, flattenArray(scopedUserRole.ProjectGroupIDs)),
		"project_ids":       schema.NewSet(schema.HashString, flattenArray(scopedUserRole.ProjectIDs)),
		"space_id":          scopedUserRole.SpaceID,
		"team_id":           scopedUserRole.TeamID,
		"tenant_ids":        schema.NewSet(schema.HashString, flattenArray(scopedUserRole.TenantIDs)),
		"user_role_id":      scopedUserRole.UserRoleID,
	}
}

func setScopedUserRole(ctx context.Context, d *schema.ResourceData, scopedUserRole *userroles.ScopedUserRole) error {
	if err := d.Set("environment_ids", scopedUserRole.EnvironmentIDs); err != nil {
		return fmt.Errorf("error setting environment_ids: %s", err)
	}

	d.Set("id", scopedUserRole.ID)
	d.SetId(scopedUserRole.GetID())

	if err := d.Set("project_group_ids", scopedUserRole.ProjectGroupIDs); err != nil {
		return fmt.Errorf("error setting project_group_ids: %s", err)
	}

	if err := d.Set("project_ids", scopedUserRole.ProjectIDs); err != nil {
		return fmt.Errorf("error setting project_ids: %s", err)
	}

	d.Set("team_id", scopedUserRole.TeamID)

	if err := d.Set("tenant_ids", scopedUserRole.TenantIDs); err != nil {
		return fmt.Errorf("error setting tenant_ids: %s", err)
	}

	d.Set("user_role_id", scopedUserRole.UserRoleID)

	return nil
}

func getScopedUserRoleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"environment_ids": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Set:      schema.HashString,
			Type:     schema.TypeSet,
		},
		"id": getIDSchema(),
		"project_group_ids": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Set:      schema.HashString,
			Type:     schema.TypeSet,
		},
		"project_ids": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Set:      schema.HashString,
			Type:     schema.TypeSet,
		},
		"space_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"team_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"tenant_ids": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Set:      schema.HashString,
			Type:     schema.TypeSet,
		},
		"user_role_id": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}
