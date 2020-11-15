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
		ReadContext: dataSourceSpacesRead,
		Schema:      getSpaceDataSchema(),
	}
}

func dataSourceSpacesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ids := expandArray(d.Get("ids").([]interface{}))
	name := d.Get("name").(string)
	partialName := d.Get("partial_name").(string)
	skip := d.Get("skip").(int)
	take := d.Get("take").(int)

	query := octopusdeploy.SpacesQuery{
		IDs:         ids,
		Name:        name,
		PartialName: partialName,
		Skip:        skip,
		Take:        take,
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
