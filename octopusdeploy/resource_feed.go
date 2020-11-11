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
	feedResource := expandFeedResource(d)

	client := m.(*octopusdeploy.Client)
	feed, err := client.Feeds.Add(feedResource)
	if err != nil {
		return diag.FromErr(err)
	}

	feedResource, err = octopusdeploy.ToFeedResource(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenFeedResource(ctx, d, feedResource)
	return nil
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
		return diag.FromErr(err)
	}
	if feed == nil {
		d.SetId("")
		return nil
	}

	feedResource := feed.(*octopusdeploy.FeedResource)

	flattenFeedResource(ctx, d, feedResource)
	return nil
}

func resourceFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feedResource := expandFeedResource(d)

	client := m.(*octopusdeploy.Client)
	feed, err := client.Feeds.Update(feedResource)
	if err != nil {
		return diag.FromErr(err)
	}

	feedResource, err = octopusdeploy.ToFeedResource(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenFeedResource(ctx, d, feedResource)
	return nil
}
