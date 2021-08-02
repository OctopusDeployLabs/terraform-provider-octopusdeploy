package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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
	dockerContainerRegistry := expandMavenFeed(d)

	log.Printf("[INFO] creating Maven feed: %#v", dockerContainerRegistry)

	client := m.(*octopusdeploy.Client)
	createdMavenFeed, err := client.Feeds.Add(dockerContainerRegistry)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setMavenFeed(ctx, d, createdMavenFeed.(*octopusdeploy.MavenFeed)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdMavenFeed.GetID())

	log.Printf("[INFO] Maven feed created (%s)", d.Id())
	return nil
}

func resourceMavenFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Maven feed (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
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

	client := m.(*octopusdeploy.Client)
	feedResource, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] Maven feed (%s) not found; deleting from state", d.Id())
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

	dockerContainerRegistry := feedResource.(*octopusdeploy.MavenFeed)
	if err := setMavenFeed(ctx, d, dockerContainerRegistry); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Maven feed read: %#v", dockerContainerRegistry)
	return nil
}

func resourceMavenFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed := expandMavenFeed(d)

	log.Printf("[INFO] updating Maven feed: %#v", feed)

	client := m.(*octopusdeploy.Client)
	updatedFeed, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	feedResource, err := octopusdeploy.ToFeed(updatedFeed.(*octopusdeploy.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setMavenFeed(ctx, d, feedResource.(*octopusdeploy.MavenFeed)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Maven feed updated (%s)", d.Id())
	return nil
}
