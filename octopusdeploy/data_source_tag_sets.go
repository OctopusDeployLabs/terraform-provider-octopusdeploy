package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTagSets() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing tag sets.",
		ReadContext: dataSourceTagSetsRead,
		Schema:      getTagSetDataSchema(),
	}
}

func dataSourceTagSetsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.TagSetsQuery{
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	client := m.(*octopusdeploy.Client)
	tagSets, err := client.TagSets.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedTagSets := []interface{}{}
	for _, tagSet := range tagSets.Items {
		flattenedTagSets = append(flattenedTagSets, flattenTagSet(tagSet))
	}

	d.Set("tag_sets", flattenedTagSets)
	d.SetId("TagSets " + time.Now().UTC().String())

	return nil
}
