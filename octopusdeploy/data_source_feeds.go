package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFeeds() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFeedsRead,
		Schema:      getFeedDataSchema(),
	}
}

func dataSourceFeedsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.FeedsQuery{
		FeedType:    d.Get("feed_type").(string),
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	client := m.(*octopusdeploy.Client)
	feeds, err := client.Feeds.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedFeeds := []interface{}{}
	for _, feed := range feeds.Items {
		feedResource, err := octopusdeploy.ToFeedResource(feed)
		if err != nil {
			return diag.FromErr(err)
		}

		flattenedFeeds = append(flattenedFeeds, flattenFeed(feedResource))
	}

	d.Set("feeds", flattenedFeeds)
	d.SetId("Feeds " + time.Now().UTC().String())

	return nil
}
