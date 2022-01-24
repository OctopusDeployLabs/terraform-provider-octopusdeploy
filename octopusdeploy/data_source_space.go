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
		Description: "Provides information about an exist space.",
		ReadContext: dataSourceSpaceRead,
		Schema:      getSpaceDataSourceSchema(),
	}
}

func dataSourceSpaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*octopusdeploy.Client)

	spaceName := d.Get("name").(string)
	space, err := client.Spaces.GetByName(spaceName)
	if err != nil {
		return diag.Errorf("Unable to find space with name '%s'", spaceName)
	}
	log.Printf("[INFO] Found space with name '%s', with ID '%s'", space.Name, space.ID)

	d.Set("id", space.ID)
	d.SetId(space.GetID())

	return nil
}
