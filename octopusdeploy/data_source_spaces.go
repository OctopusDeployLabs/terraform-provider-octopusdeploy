package octopusdeploy

import (
	"context"
	"log"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSpaces() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing spaces.",
		ReadContext: dataSourceSpacesRead,
		Schema:      getSpaceDataSchema(),
	}
}

func dataSourceSpacesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*octopusdeploy.Client)

	flattenedSpaces := []interface{}{}

	space_name := d.Get("name").(string)
	if space_name != "" {
		space, err := client.Spaces.GetByName(space_name)
		if err != nil {
			return diag.Errorf("Unable to find space with name '%s'", space_name)
		}
		log.Printf("[INFO] Found space with name '%s', with ID '%s'", space.Name, space.ID)

		flattenedSpaces = append(flattenedSpaces, flattenSpace(space))
	} else {
		query := octopusdeploy.SpacesQuery{
			IDs:         expandArray(d.Get("ids").([]interface{})),
			PartialName: d.Get("partial_name").(string),
			Skip:        d.Get("skip").(int),
			Take:        d.Get("take").(int),
		}

		spaces, err := client.Spaces.Get(query)
		if err != nil {
			return diag.FromErr(err)
		}

		for _, space := range spaces.Items {
			flattenedSpaces = append(flattenedSpaces, flattenSpace(space))
		}
	}

	d.Set("spaces", flattenedSpaces)
	d.SetId("Spaces " + time.Now().UTC().String())

	return nil
}
