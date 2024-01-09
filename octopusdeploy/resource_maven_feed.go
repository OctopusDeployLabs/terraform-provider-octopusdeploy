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
	mavenFeed, err := expandMavenFeed(d)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("creating Maven feed: %s", mavenFeed.GetName()))

	client := m.(*client.Client)
	createdFeed, err := feeds.Add(client, mavenFeed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setMavenFeed(ctx, d, createdFeed.(*feeds.MavenFeed)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdFeed.GetID())

	tflog.Info(ctx, fmt.Sprintf("Maven feed created (%s)", d.Id()))
	return nil
}

func resourceMavenFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("deleting Maven feed (%s)", d.Id()))

	client := m.(*client.Client)
	err := feeds.DeleteByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	tflog.Info(ctx, "Maven feed deleted")
	return nil
}

func resourceMavenFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("reading Maven feed (%s)", d.Id()))

	client := m.(*client.Client)
	feed, err := feeds.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Maven feed")
	}

	mavenFeed := feed.(*feeds.MavenFeed)
	if err := setMavenFeed(ctx, d, mavenFeed); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Maven feed read (%s)", mavenFeed.GetID()))
	return nil
}

func resourceMavenFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed, err := expandMavenFeed(d)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("updating Maven feed (%s)", feed.GetID()))

	client := m.(*client.Client)
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setMavenFeed(ctx, d, updatedFeed.(*feeds.MavenFeed)); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Maven feed updated (%s)", d.Id()))
	return nil
}
