package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
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

	client := m.(*client.Client)

	spaceName := d.Get("name").(string)
	existingSpace, err := client.Spaces.GetByName(spaceName)
	if err != nil {
		return diag.Errorf("unable to find space with name '%s'", spaceName)
	}
	log.Printf("[INFO] Found space with name '%s', with ID '%s'", existingSpace.Name, existingSpace.ID)

	d.Set("id", existingSpace.ID)
	d.Set("description", existingSpace.Description)
	d.Set("is_default", existingSpace.IsDefault)
	d.Set("is_task_queue_stopped", existingSpace.TaskQueueStopped)
	d.Set("space_managers_team_members", existingSpace.SpaceManagersTeamMembers)
	d.Set("space_managers_teams", existingSpace.SpaceManagersTeams)
	d.SetId(existingSpace.GetID())

	return nil
}
