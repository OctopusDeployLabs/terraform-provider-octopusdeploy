package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/userroles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUserRoles() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing user roles.",
		ReadContext: dataSourceUserRolesRead,
		Schema:      getUserRoleDataSchema(),
	}
}

func dataSourceUserRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	query := userroles.UserRolesQuery{
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	spaceID := d.Get("space_id").(string)

	client := meta.(*client.Client)
	existingUserRoles, err := userroles.Get(client, spaceID, query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedUserRoles := []interface{}{}
	for _, userRole := range existingUserRoles.Items {
		flattenedUserRoles = append(flattenedUserRoles, flattenUserRole(userRole))
	}

	d.Set("user_roles", flattenedUserRoles)
	d.SetId("UserRoles " + time.Now().UTC().String())

	return nil
}
