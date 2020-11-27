package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFeed() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFeedCreate,
		DeleteContext: resourceFeedDelete,
		Importer:      getImporter(),
		ReadContext:   resourceFeedRead,
		Schema:        getFeedSchema(),
		UpdateContext: resourceFeedUpdate,
	}
}

func resourceFeedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feedResource := expandFeed(d)

	client := m.(*octopusdeploy.Client)
	createdFeed, err := client.Feeds.Add(feedResource)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdFeed.GetID())
	return resourceFeedRead(ctx, d, m)
}

func resourceFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Feeds.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	feed, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	feedResource := feed.(*octopusdeploy.FeedResource)

	setFeed(ctx, d, feedResource)
	return nil
}

func resourceFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feedResource := expandFeed(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Feeds.Update(feedResource)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFeedRead(ctx, d, m)
}
