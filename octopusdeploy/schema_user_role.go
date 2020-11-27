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
	userRoleSchema := getUserRoleSchema()
	for _, field := range userRoleSchema {
		field.Computed = true
		field.Default = nil
		field.MaxItems = 0
		field.MinItems = 0
		field.Optional = false
		field.Required = false
		field.ValidateDiagFunc = nil
	}

	return map[string]*schema.Schema{
		"ids": {
			Description: "Query and/or search by a list of IDs",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"partial_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"skip": {
			Default:     0,
			Description: "Indicates the number of items to skip in the response",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"take": {
			Default:     1,
			Description: "Indicates the number of items to take (or return) in the response",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"user_roles": {
			Computed: true,
			Elem:     &schema.Resource{Schema: userRoleSchema},
			Type:     schema.TypeList,
		},
	}
}

func getUserRoleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"can_be_deleted": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
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
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"name": {
			Optional: true,
			Type:     schema.TypeString,
		},
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
