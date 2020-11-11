package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNuGetFeed() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNuGetFeedCreate,
		DeleteContext: resourceNuGetFeedDelete,
		Importer:      getImporter(),
		ReadContext:   resourceNuGetFeedRead,
		Schema:        getNuGetFeedSchema(),
		UpdateContext: resourceNuGetFeedUpdate,
	}
}

func resourceNuGetFeedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed := expandNuGetFeed(d)

	client := m.(*octopusdeploy.Client)
	createdFeed, err := client.Feeds.Add(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	createdNuGetFeed := createdFeed.(*octopusdeploy.NuGetFeed)

	flattenNuGetFeed(ctx, d, createdNuGetFeed)
	return nil
}

func resourceNuGetFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	feedResource, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if feedResource == nil {
		d.SetId("")
		return nil
	}

	feedResource, err = octopusdeploy.ToFeed(feedResource.(*octopusdeploy.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	nuGetFeed := feedResource.(*octopusdeploy.NuGetFeed)

	flattenNuGetFeed(ctx, d, nuGetFeed)
	return nil
}

func resourceNuGetFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Feeds.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceNuGetFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed := expandNuGetFeed(d)

	client := m.(*octopusdeploy.Client)
	feedResource, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	feedResource, err = octopusdeploy.ToFeed(feedResource.(*octopusdeploy.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	updatedNuGetFeed := feedResource.(*octopusdeploy.NuGetFeed)

	flattenNuGetFeed(ctx, d, updatedNuGetFeed)
	return nil
}
