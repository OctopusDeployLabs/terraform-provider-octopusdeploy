package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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
	dockerContainerRegistry := expandHelmFeed(d)

	log.Printf("[INFO] creating Helm feed: %#v", dockerContainerRegistry)

	client := m.(*octopusdeploy.Client)
	createdHelmFeed, err := client.Feeds.Add(dockerContainerRegistry)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setHelmFeed(ctx, d, createdHelmFeed.(*octopusdeploy.HelmFeed)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdHelmFeed.GetID())

	log.Printf("[INFO] Helm feed created (%s)", d.Id())
	return nil
}

func resourceHelmFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Helm feed (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
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

	client := m.(*octopusdeploy.Client)
	feedResource, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] Helm feed (%s) not found; deleting from state", d.Id())
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

	dockerContainerRegistry := feedResource.(*octopusdeploy.HelmFeed)
	if err := setHelmFeed(ctx, d, dockerContainerRegistry); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Helm feed read: %#v", dockerContainerRegistry)
	return nil
}

func resourceHelmFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed := expandHelmFeed(d)

	log.Printf("[INFO] updating Helm feed: %#v", feed)

	client := m.(*octopusdeploy.Client)
	updatedFeed, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	feedResource, err := octopusdeploy.ToFeed(updatedFeed.(*octopusdeploy.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setHelmFeed(ctx, d, feedResource.(*octopusdeploy.HelmFeed)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Helm feed updated (%s)", d.Id())
	return nil
}
