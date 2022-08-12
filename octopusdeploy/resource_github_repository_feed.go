package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
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

	log.Printf("[INFO] creating GitHub repository feed, %s", feed.GetName())

	client := m.(*client.Client)
	createdGitHubRepositoryFeed, err := client.Feeds.Add(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setGitHubRepositoryFeed(ctx, d, createdGitHubRepositoryFeed.(*feeds.GitHubRepositoryFeed)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdGitHubRepositoryFeed.GetID())

	log.Printf("[INFO] GitHub repository feed created (%s)", d.Id())
	return nil
}

func resourceGitHubRepositoryFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting GitHub repository feed (%s)", d.Id())

	client := m.(*client.Client)
	err := client.Feeds.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] GitHub repository feed deleted")
	return nil
}

func resourceGitHubRepositoryFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading GitHub repository feed (%s)", d.Id())

	client := m.(*client.Client)
	feedResource, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "GitHub repository feed")
	}

	feedResource, err = feeds.ToFeed(feedResource.(*feeds.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	gitHubRepositoryFeed := feedResource.(*feeds.GitHubRepositoryFeed)
	if err := setGitHubRepositoryFeed(ctx, d, gitHubRepositoryFeed); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] GitHub repository feed read (%s)", gitHubRepositoryFeed.GetID())
	return nil
}

func resourceGitHubRepositoryFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed, err := expandGitHubRepositoryFeed(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] updating GitHub repository feed (%s)", feed.GetID())

	client := m.(*client.Client)
	updatedFeed, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	feedResource, err := feeds.ToFeed(updatedFeed.(*feeds.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setGitHubRepositoryFeed(ctx, d, feedResource.(*feeds.GitHubRepositoryFeed)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] GitHub repository feed updated (%s)", d.Id())
	return nil
}
