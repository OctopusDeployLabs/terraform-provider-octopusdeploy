package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	tflog.Info(ctx, fmt.Sprintf("creating NuGet feed: %s", feed.GetName()))

	client := m.(*client.Client)
	createdFeed, err := client.Feeds.Add(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setNuGetFeed(ctx, d, createdFeed.(*feeds.NuGetFeed)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdFeed.GetID())

	tflog.Info(ctx, fmt.Sprintf("NuGet feed created (%s)", d.Id()))
	return nil
}

func resourceNuGetFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("deleting NuGet feed (%s)", d.Id()))

	client := m.(*client.Client)
	err := client.Feeds.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	tflog.Info(ctx, "NuGet feed deleted")
	return nil
}

func resourceNuGetFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("reading NuGet feed (%s)", d.Id()))

	client := m.(*client.Client)
	feed, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "NuGet feed")
	}

	nuGetFeed := feed.(*feeds.NuGetFeed)
	if err := setNuGetFeed(ctx, d, nuGetFeed); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("NuGet feed read (%s)", nuGetFeed.GetID()))
	return nil
}

func resourceNuGetFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed, err := expandNuGetFeed(d)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("updating NuGet feed (%s)", feed.GetID()))

	client := m.(*client.Client)
	updatedFeed, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setNuGetFeed(ctx, d, updatedFeed.(*feeds.NuGetFeed)); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("NuGet feed updated (%s)", d.Id()))
	return nil
}
