package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSpace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSpaceReadByName,
		Schema: map[string]*schema.Schema{
			constName: {
				Required: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func dataSpaceReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	name := d.Get(constName).(string)

	space, err := client.Spaces.GetByName(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if space == nil {
		d.SetId(constEmptyString)
		return diag.FromErr(fmt.Errorf("Unable to retrieve space (name: %s)", name))
	}

	d.SetId(space.GetID())
	d.Set(constDescription, space.Description)
	d.Set(constIsDefault, space.IsDefault)
	d.Set(constName, space.Name)
	d.Set(constSpaceManagersTeamMembers, space.SpaceManagersTeamMembers)
	d.Set(constSpaceManagersTeams, space.SpaceManagersTeams)
	d.Set(constTaskQueueStopped, space.TaskQueueStopped)

	return nil
}
