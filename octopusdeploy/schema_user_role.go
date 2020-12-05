package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandUserRole(d *schema.ResourceData) *octopusdeploy.UserRole {
	userRole := &octopusdeploy.UserRole{}
	userRole.ID = d.Id()

	if v, ok := d.GetOk("can_be_deleted"); ok {
		userRole.CanBeDeleted = v.(bool)
	}

	if v, ok := d.GetOk("description"); ok {
		userRole.Description = v.(string)
	}

	if v, ok := d.GetOk("granted_space_permissions"); ok {
		userRole.GrantedSpacePermissions = v.([]string)
	}

	if v, ok := d.GetOk("granted_system_permissions"); ok {
		userRole.GrantedSystemPermissions = v.([]string)
	}

	if v, ok := d.GetOk("name"); ok {
		userRole.Name = v.(string)
	}

	if v, ok := d.GetOk("space_permission_descriptions"); ok {
		userRole.SpacePermissionDescriptions = v.([]string)
	}

	if v, ok := d.GetOk("supported_restrictions"); ok {
		userRole.SupportedRestrictions = v.([]string)
	}

	if v, ok := d.GetOk("system_permission_descriptions"); ok {
		userRole.SystemPermissionDescriptions = v.([]string)
	}

	return userRole
}

func flattenUserRole(userRole *octopusdeploy.UserRole) map[string]interface{} {
	if userRole == nil {
		return nil
	}

	return map[string]interface{}{
		"can_be_deleted":                 userRole.CanBeDeleted,
		"description":                    userRole.Description,
		"granted_space_permissions":      userRole.GrantedSpacePermissions,
		"granted_system_permissions":     userRole.GrantedSystemPermissions,
		"id":                             userRole.GetID(),
		"name":                           userRole.Name,
		"space_permission_descriptions":  userRole.SpacePermissionDescriptions,
		"supported_restrictions":         userRole.SupportedRestrictions,
		"system_permission_descriptions": userRole.SystemPermissionDescriptions,
	}
}

func getUserRoleDataSchema() map[string]*schema.Schema {
	dataSchema := getUserRoleSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"id":           getDataSchemaID(),
		"ids":          getQueryIDs(),
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
		"user_role": {
			Computed:    true,
			Description: "A list of user roles that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}

func getUserRoleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"can_be_deleted": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"description": getDescriptionSchema(),
		"granted_space_permissions": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"granted_system_permissions": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"id":   getIDSchema(),
		"name": getNameSchema(true),
		"space_permission_descriptions": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"supported_restrictions": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"system_permission_descriptions": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}

func setUserRole(ctx context.Context, d *schema.ResourceData, userRole *octopusdeploy.UserRole) {
	d.Set("can_be_deleted", userRole.CanBeDeleted)
	d.Set("description", userRole.Description)
	d.Set("id", userRole.GetID())
	d.Set("granted_space_permissions", userRole.GrantedSpacePermissions)
	d.Set("granted_system_permissions", userRole.GrantedSystemPermissions)
	d.Set("name", userRole.Name)
	d.Set("space_permission_descriptions", userRole.SpacePermissionDescriptions)
	d.Set("supported_restrictions", userRole.SupportedRestrictions)
	d.Set("system_permission_descriptions", userRole.SystemPermissionDescriptions)

	d.SetId(userRole.GetID())
}
