package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserReadByName,
		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceUserReadByName(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*octopusdeploy.Client)
	username := d.Get(constUsername).(string)
	query := octopusdeploy.UsersQuery{
		Filter: username,
		Take:   1,
	}

	users, err := client.Users.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}
	if users == nil || len(users.Items) == 0 {
		return diag.FromErr(fmt.Errorf("Unabled to retrieve user (filter: %s)", username))
	}

	// NOTE: two or more users can have the same name in Octopus and
	// therefore, a better search criteria needs to be implemented below

	for _, user := range users.Items {
		if user.Username == constUsername {
			logResource(constUser, meta)

			d.SetId(user.GetID())
			d.Set(constCanPasswordBeEdited, user.CanPasswordBeEdited)
			d.Set(constDisplayName, user.DisplayName)
			d.Set(constEmailAddress, user.EmailAddress)
			d.Set(constIdentities, user.Identities)
			d.Set(constIsActive, user.IsActive)
			d.Set(constIsRequestor, user.IsRequestor)
			d.Set(constIsService, user.IsService)
			d.Set(constPassword, user.Password)
			d.Set(constUsername, user.Username)

			return nil
		}
	}

	return nil
}
