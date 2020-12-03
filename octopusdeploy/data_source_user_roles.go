package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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
	query := octopusdeploy.UserRolesQuery{
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	client := meta.(*octopusdeploy.Client)
	users, err := client.UserRoles.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedUserRoles := []interface{}{}
	for _, user := range users.Items {
		flattenedUserRoles = append(flattenedUserRoles, flattenUserRole(user))
	}

	d.Set("user_roles", flattenedUserRoles)
	d.SetId("UserRoles " + time.Now().UTC().String())

	return nil
}
