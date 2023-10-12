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

func resourceGitHubRepositoryFeed() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitHubRepositoryFeedCreate,
		DeleteContext: resourceGitHubRepositoryFeedDelete,
		Description:   "This resource manages a GitHub repository feed in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceGitHubRepositoryFeedRead,
		Schema:        getGitHubRepositoryFeedSchema(),
		UpdateContext: resourceGitHubRepositoryFeedUpdate,
	}
}

func resourceGitHubRepositoryFeedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed, err := expandGitHubRepositoryFeed(d)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("creating GitHub repository feed, %s", feed.GetName()))

	client := m.(*client.Client)
	createdGitHubRepositoryFeed, err := feeds.Add(client, feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setGitHubRepositoryFeed(ctx, d, createdGitHubRepositoryFeed.(*feeds.GitHubRepositoryFeed)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdGitHubRepositoryFeed.GetID())

	tflog.Info(ctx, fmt.Sprintf("GitHub repository feed created (%s)", d.Id()))
	return nil
}

func resourceGitHubRepositoryFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("deleting GitHub repository feed (%s)", d.Id()))

	client := m.(*client.Client)
	err := feeds.DeleteByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	tflog.Info(ctx, "GitHub repository feed deleted")
	return nil
}

func resourceGitHubRepositoryFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("reading GitHub repository feed (%s)", d.Id()))

	client := m.(*client.Client)
	feed, err := feeds.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "GitHub repository feed")
	}

	gitHubRepositoryFeed := feed.(*feeds.GitHubRepositoryFeed)
	if err := setGitHubRepositoryFeed(ctx, d, gitHubRepositoryFeed); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("GitHub repository feed read (%s)", gitHubRepositoryFeed.GetID()))
	return nil
}

func resourceGitHubRepositoryFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed, err := expandGitHubRepositoryFeed(d)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("updating GitHub repository feed (%s)", feed.GetID()))

	client := m.(*client.Client)
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setGitHubRepositoryFeed(ctx, d, updatedFeed.(*feeds.GitHubRepositoryFeed)); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("GitHub repository feed updated (%s)", d.Id()))
	return nil
}
