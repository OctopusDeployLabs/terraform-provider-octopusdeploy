package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/users"
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
	query := users.UsersQuery{
		Filter: d.Get("filter").(string),
		IDs:    expandArray(d.Get("ids").([]interface{})),
		Skip:   d.Get("skip").(int),
		Take:   d.Get("take").(int),
	}

	spaceID := d.Get("space_id").(string)

	client := meta.(*client.Client)
	existingUsers, err := users.Get(client, spaceID, query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedUsers := []interface{}{}
	for _, user := range existingUsers.Items {
		flattenedUsers = append(flattenedUsers, flattenUser(user))
	}

	d.Set("users", flattenedUsers)
	d.SetId("Users " + time.Now().UTC().String())

	return nil
}
