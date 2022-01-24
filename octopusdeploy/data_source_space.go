package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSpace() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about an existing space.",
		ReadContext: dataSourceSpaceRead,
		Schema:      getSpaceDataSourceSchema(),
	}
}

func dataSourceSpaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*octopusdeploy.Client)

	spaceName := d.Get("name").(string)
	space, err := client.Spaces.GetByName(spaceName)
	if err != nil {
		return diag.Errorf("unable to find space with name '%s'", spaceName)
	}
	log.Printf("[INFO] Found space with name '%s', with ID '%s'", space.Name, space.ID)

	d.Set("id", space.ID)
	d.Set("description", space.Description)
	d.Set("is_default", space.IsDefault)
	d.Set("is_task_queue_stopped", space.TaskQueueStopped)
	d.Set("space_managers_team_members", space.SpaceManagersTeamMembers)
	d.Set("space_managers_teams", space.SpaceManagersTeams)
	d.SetId(space.GetID())

	return nil
}
