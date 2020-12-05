package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing users.",
		ReadContext: dataSourceUsersRead,
		Schema:      getUserDataSchema(),
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	query := octopusdeploy.UsersQuery{
		Filter: d.Get("filter").(string),
		IDs:    expandArray(d.Get("ids").([]interface{})),
		Skip:   d.Get("skip").(int),
		Take:   d.Get("take").(int),
	}

	client := meta.(*octopusdeploy.Client)
	users, err := client.Users.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedUsers := []interface{}{}
	for _, user := range users.Items {
		flattenedUsers = append(flattenedUsers, flattenUser(user))
	}

	d.Set("user", flattenedUsers)
	d.SetId("Users " + time.Now().UTC().String())

	return nil
}
