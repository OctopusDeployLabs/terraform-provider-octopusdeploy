package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserRoleCreate,
		DeleteContext: resourceUserRoleDelete,
		Description:   "This resource manages user roles in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceUserRoleRead,
		Schema:        getUserRoleSchema(),
		UpdateContext: resourceUserRoleUpdate,
	}
}

func resourceUserRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userRole := expandUserRole(d)

	log.Printf("[INFO] creating user role: %#v", userRole)

	client := m.(*octopusdeploy.Client)
	createdUserRole, err := client.UserRoles.Add(userRole)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setUserRole(ctx, d, createdUserRole); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdUserRole.GetID())

	log.Printf("[INFO] user role created (%s)", d.Id())
	return nil
}

func resourceUserRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting user role (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.UserRoles.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] user role deleted")
	return nil
}

func resourceUserRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading user role (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	userRole, err := client.UserRoles.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] user role (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if err := setUserRole(ctx, d, userRole); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] user role read (%s)", d.Id())
	return nil
}

func resourceUserRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating user role (%s)", d.Id())

	userRole := expandUserRole(d)
	client := m.(*octopusdeploy.Client)
	updatedUserRole, err := client.UserRoles.Update(userRole)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setUserRole(ctx, d, updatedUserRole); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] user role updated (%s)", d.Id())
	return nil
}
