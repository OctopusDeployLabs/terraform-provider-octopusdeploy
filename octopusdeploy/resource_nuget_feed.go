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
		Description:   "This resource manages a NuGet feed in Octopus Deploy.",
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

	d.SetId(createdFeed.GetID())
	return resourceNuGetFeedRead(ctx, d, m)
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

func resourceNuGetFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	feedResource, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	feedResource, err = octopusdeploy.ToFeed(feedResource.(*octopusdeploy.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	nuGetFeed := feedResource.(*octopusdeploy.NuGetFeed)

	setNuGetFeed(ctx, d, nuGetFeed)
	return nil
}

func resourceNuGetFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed := expandNuGetFeed(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNuGetFeedRead(ctx, d, m)
}
