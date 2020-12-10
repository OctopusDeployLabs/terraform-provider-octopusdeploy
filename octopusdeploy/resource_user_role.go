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
		DeleteContext: resourceUserRoleDelete,
		Description:   "This resource manages user roles in Octopus Deploy.",
		Importer:      getImporter(),
		Schema:        getUserRoleSchema(),
	}
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
