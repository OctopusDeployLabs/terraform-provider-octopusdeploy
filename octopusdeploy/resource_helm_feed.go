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

func resourceHelmFeed() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHelmFeedCreate,
		DeleteContext: resourceHelmFeedDelete,
		Description:   "This resource manages a Helm feed in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceHelmFeedRead,
		Schema:        getHelmFeedSchema(),
		UpdateContext: resourceHelmFeedUpdate,
	}
}

func resourceHelmFeedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed, err := expandHelmFeed(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] creating Helm feed, %s", feed.GetName())

	client := m.(*client.Client)
	createdFeed, err := client.Feeds.Add(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setHelmFeed(ctx, d, createdFeed.(*feeds.HelmFeed)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdFeed.GetID())

	log.Printf("[INFO] Helm feed created (%s)", d.Id())
	return nil
}

func resourceHelmFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Helm feed (%s)", d.Id())

	client := m.(*client.Client)
	err := client.Feeds.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] Helm feed deleted")
	return nil
}

func resourceHelmFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading Helm feed (%s)", d.Id())

	client := m.(*client.Client)
	feedResource, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Helm feed")
	}

	feedResource, err = feeds.ToFeed(feedResource.(*feeds.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	helmFeed := feedResource.(*feeds.HelmFeed)
	if err := setHelmFeed(ctx, d, helmFeed); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Helm feed read (%s)", helmFeed.GetID())
	return nil
}

func resourceHelmFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed, err := expandHelmFeed(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] updating Helm feed (%s)", feed.GetID())

	client := m.(*client.Client)
	updatedFeed, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setHelmFeed(ctx, d, updatedFeed.(*feeds.HelmFeed)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Helm feed updated (%s)", d.Id())
	return nil
}
