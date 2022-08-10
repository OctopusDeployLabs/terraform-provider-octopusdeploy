package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMavenFeed() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMavenFeedCreate,
		DeleteContext: resourceMavenFeedDelete,
		Description:   "This resource manages a Maven feed in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceMavenFeedRead,
		Schema:        getMavenFeedSchema(),
		UpdateContext: resourceMavenFeedUpdate,
	}
}

func resourceMavenFeedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mavenFeed, err := expandMavenFeed(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] creating Maven feed: %s", mavenFeed.GetName())

	client := m.(*client.Client)
	createdFeed, err := client.Feeds.Add(mavenFeed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setMavenFeed(ctx, d, createdFeed.(*feeds.MavenFeed)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdFeed.GetID())

	log.Printf("[INFO] Maven feed created (%s)", d.Id())
	return nil
}

func resourceMavenFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Maven feed (%s)", d.Id())

	client := m.(*client.Client)
	err := client.Feeds.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] Maven feed deleted")
	return nil
}

func resourceMavenFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading Maven feed (%s)", d.Id())

	client := m.(*client.Client)
	feedResource, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Maven feed")
	}

	feedResource, err = feeds.ToFeed(feedResource.(*feeds.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	mavenFeed := feedResource.(*feeds.MavenFeed)
	if err := setMavenFeed(ctx, d, mavenFeed); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Maven feed read (%s)", mavenFeed.GetID())
	return nil
}

func resourceMavenFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed, err := expandMavenFeed(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] updating Maven feed (%s)", feed.GetID())

	client := m.(*client.Client)
	updatedFeed, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	feedResource, err := feeds.ToFeed(updatedFeed.(*feeds.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setMavenFeed(ctx, d, feedResource.(*feeds.MavenFeed)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Maven feed updated (%s)", d.Id())
	return nil
}
