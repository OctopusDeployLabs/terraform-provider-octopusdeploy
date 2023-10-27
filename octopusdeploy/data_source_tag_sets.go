package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
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
	query := tagsets.TagSetsQuery{
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}
	spaceID := d.Get("space_id").(string)

	octopus := m.(*client.Client)
	existingTagSets, err := tagsets.Get(octopus, spaceID, query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedTagSets := []interface{}{}
	
	for _, tagSet := range existingTagSets.Items {
		flattenedTagSets = append(flattenedTagSets, flattenTagSet(tagSet))
	}

	d.Set("tag_sets", flattenedTagSets)
	d.SetId("TagSets " + time.Now().UTC().String())

	return nil
}
