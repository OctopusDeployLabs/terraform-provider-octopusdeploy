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

func resourceArtifactoryGenericFeed() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArtifactoryGenericFeedCreate,
		DeleteContext: resourceArtifactoryGenericFeedDelete,
		Description:   "This resource manages a Artifactory Generic feed in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceArtifactoryGenericFeedRead,
		Schema:        getArtifactoryGenericFeedSchema(),
		UpdateContext: resourceArtifactoryGenericFeedUpdate,
	}
}

func resourceArtifactoryGenericFeedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	artifactoryGenericFeed, err := expandArtifactoryGenericFeed(d)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("creating Artifactory Generic feed: %s", artifactoryGenericFeed.GetName()))

	client := m.(*client.Client)
	createdFeed, err := feeds.Add(client, artifactoryGenericFeed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setArtifactoryGenericFeed(ctx, d, createdFeed.(*feeds.ArtifactoryGenericFeed)); err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, fmt.Sprintf("layout regex from created model: %s", createdFeed.(*feeds.ArtifactoryGenericFeed).LayoutRegex))
	d.SetId(createdFeed.GetID())

	tflog.Info(ctx, fmt.Sprintf("Artifactory Generic feed created (%s)", d.Id()))
	return nil
}

func resourceArtifactoryGenericFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("deleting Artifactory Generic feed (%s)", d.Id()))

	client := m.(*client.Client)
	err := feeds.DeleteByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	tflog.Info(ctx, "Artifactory Generic feed deleted")
	return nil
}

func resourceArtifactoryGenericFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("reading Artifactory Generic feed (%s)", d.Id()))

	client := m.(*client.Client)
	feed, err := feeds.GetByID(client, d.Get("space_id").(string), d.Id())

	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Artifactory Generic feed")
	}

	if err := setArtifactoryGenericFeed(ctx, d, feed.(*feeds.ArtifactoryGenericFeed)); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Artifactory Generic feed read (%s)", feed.GetID()))
	return nil
}

func resourceArtifactoryGenericFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed, err := expandArtifactoryGenericFeed(d)

	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("updating Artifactory Generic feed (%s)", feed.GetID()))

	client := m.(*client.Client)
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setArtifactoryGenericFeed(ctx, d, updatedFeed.(*feeds.ArtifactoryGenericFeed)); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Artifactory Generic feed updated (%s)", d.Id()))
	return nil
}
