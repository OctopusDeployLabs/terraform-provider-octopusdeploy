package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFeeds() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing feeds.",
		ReadContext: dataSourceFeedsRead,
		Schema:      getFeedDataSchema(),
	}
}

func dataSourceFeedsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := feeds.FeedsQuery{
		FeedType:    d.Get("feed_type").(string),
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	client := m.(*client.Client)
	existingFeeds, err := client.Feeds.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedFeeds := []interface{}{}
	for _, feed := range existingFeeds.Items {
		feedResource, err := feeds.ToFeedResource(feed)
		if err != nil {
			return diag.FromErr(err)
		}

		flattenedFeeds = append(flattenedFeeds, flattenFeed(feedResource))
	}

	d.Set("feeds", flattenedFeeds)
	d.SetId("Feeds " + time.Now().UTC().String())

	return nil
}
