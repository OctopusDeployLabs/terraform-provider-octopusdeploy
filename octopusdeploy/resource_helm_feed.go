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

	tflog.Info(ctx, fmt.Sprintf("creating Helm feed, %s", feed.GetName()))

	client := m.(*client.Client)
	createdFeed, err := client.Feeds.Add(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setHelmFeed(ctx, d, createdFeed.(*feeds.HelmFeed)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdFeed.GetID())

	tflog.Info(ctx, fmt.Sprintf("Helm feed created (%s)", d.Id()))
	return nil
}

func resourceHelmFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("deleting Helm feed (%s)", d.Id()))

	client := m.(*client.Client)
	err := client.Feeds.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	tflog.Info(ctx, "Helm feed deleted")
	return nil
}

func resourceHelmFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("reading Helm feed (%s)", d.Id()))

	client := m.(*client.Client)
	feed, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Helm feed")
	}

	helmFeed := feed.(*feeds.HelmFeed)
	if err := setHelmFeed(ctx, d, helmFeed); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Helm feed read (%s)", helmFeed.GetID()))
	return nil
}

func resourceHelmFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed, err := expandHelmFeed(d)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("updating Helm feed (%s)", feed.GetID()))

	client := m.(*client.Client)
	updatedFeed, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setHelmFeed(ctx, d, updatedFeed.(*feeds.HelmFeed)); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Helm feed updated (%s)", d.Id()))
	return nil
}
