package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNuGetFeed() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNuGetFeedCreate,
		DeleteContext: resourceNuGetFeedDelete,
		Description:   "This resource manages a NuGet feed in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceNuGetFeedRead,
		Schema:        getNuGetFeedSchema(),
		UpdateContext: resourceNuGetFeedUpdate,
	}
}

func resourceNuGetFeedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed, err := expandNuGetFeed(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] creating NuGet feed: %s", feed.GetName())

	client := m.(*octopusdeploy.Client)
	createdFeed, err := client.Feeds.Add(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setNuGetFeed(ctx, d, createdFeed.(*octopusdeploy.NuGetFeed)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdFeed.GetID())

	log.Printf("[INFO] NuGet feed created (%s)", d.Id())
	return nil
}

func resourceNuGetFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting NuGet feed (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	err := client.Feeds.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] NuGet feed deleted")
	return nil
}

func resourceNuGetFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading NuGet feed (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	feedResource, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] NuGet feed (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	feedResource, err = octopusdeploy.ToFeed(feedResource.(*octopusdeploy.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	nuGetFeed := feedResource.(*octopusdeploy.NuGetFeed)
	if err := setNuGetFeed(ctx, d, nuGetFeed); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] NuGet feed read (%s)", nuGetFeed.GetID())
	return nil
}

func resourceNuGetFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed, err := expandNuGetFeed(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] updating NuGet feed (%s)", feed.GetID())

	client := m.(*octopusdeploy.Client)
	updatedFeed, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	feedResource, err := octopusdeploy.ToFeed(updatedFeed.(*octopusdeploy.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setNuGetFeed(ctx, d, feedResource.(*octopusdeploy.NuGetFeed)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] NuGet feed updated (%s)", d.Id())
	return nil
}
