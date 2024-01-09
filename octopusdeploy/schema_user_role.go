package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/userroles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandUserRole(d *schema.ResourceData) *userroles.UserRole {
	userRole := &userroles.UserRole{}
	userRole.ID = d.Id()

	if v, ok := d.GetOk("can_be_deleted"); ok {
		userRole.CanBeDeleted = v.(bool)
	}

	if v, ok := d.GetOk("description"); ok {
		userRole.Description = v.(string)
	}

	if v, ok := d.GetOk("granted_space_permissions"); ok {
		userRole.GrantedSpacePermissions = getSliceFromTerraformTypeList(v)
	} else {
		userRole.GrantedSpacePermissions = []string{}
	}

	if v, ok := d.GetOk("granted_system_permissions"); ok {
		userRole.GrantedSystemPermissions = getSliceFromTerraformTypeList(v)
	} else {
		userRole.GrantedSystemPermissions = []string{}
	}

	if v, ok := d.GetOk("name"); ok {
		userRole.Name = v.(string)
	}

	if v, ok := d.GetOk("space_permission_descriptions"); ok {
		userRole.SpacePermissionDescriptions = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("supported_restrictions"); ok {
		userRole.SupportedRestrictions = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("system_permission_descriptions"); ok {
		userRole.SystemPermissionDescriptions = getSliceFromTerraformTypeList(v)
	}

	return userRole
}

func flattenUserRole(userRole *userroles.UserRole) map[string]interface{} {
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
		"space_id":     getQuerySpaceID(),
		"user_roles": {
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
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"description": getDescriptionSchema("user role"),
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
			Computed: true,
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
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}

func setUserRole(ctx context.Context, d *schema.ResourceData, userRole *userroles.UserRole) error {
	d.Set("can_be_deleted", userRole.CanBeDeleted)
	d.Set("description", userRole.Description)
	d.Set("id", userRole.GetID())

	if err := d.Set("granted_space_permissions", userRole.GrantedSpacePermissions); err != nil {
		return fmt.Errorf("error setting granted_space_permissions: %s", err)
	}

	if err := d.Set("granted_system_permissions", userRole.GrantedSystemPermissions); err != nil {
		return fmt.Errorf("error setting granted_system_permissions: %s", err)
	}

	d.Set("name", userRole.Name)

	if err := d.Set("space_permission_descriptions", userRole.SpacePermissionDescriptions); err != nil {
		return fmt.Errorf("error setting space_permission_descriptions: %s", err)
	}

	if err := d.Set("supported_restrictions", userRole.SupportedRestrictions); err != nil {
		return fmt.Errorf("error setting supported_restrictions: %s", err)
	}

	if err := d.Set("system_permission_descriptions", userRole.SystemPermissionDescriptions); err != nil {
		return fmt.Errorf("error setting system_permission_descriptions: %s", err)
	}

	d.SetId(userRole.GetID())

	return nil
}
