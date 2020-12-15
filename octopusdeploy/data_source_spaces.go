package octopusdeploy

import (
	"context"
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
	query := octopusdeploy.SpacesQuery{
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	client := m.(*octopusdeploy.Client)
	spaces, err := client.Spaces.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedSpaces := []interface{}{}
	for _, space := range spaces.Items {
		flattenedSpaces = append(flattenedSpaces, flattenSpace(space))
	}

	d.Set("spaces", flattenedSpaces)
	d.SetId("Spaces " + time.Now().UTC().String())

	return nil
}
