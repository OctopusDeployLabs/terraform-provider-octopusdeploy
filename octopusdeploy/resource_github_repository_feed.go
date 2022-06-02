package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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

	client := m.(*octopusdeploy.Client)
	createdGitHubRepositoryFeed, err := client.Feeds.Add(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setGitHubRepositoryFeed(ctx, d, createdGitHubRepositoryFeed.(*octopusdeploy.GitHubRepositoryFeed)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdGitHubRepositoryFeed.GetID())

	log.Printf("[INFO] GitHub repository feed created (%s)", d.Id())
	return nil
}

func resourceGitHubRepositoryFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting GitHub repository feed (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
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

	client := m.(*octopusdeploy.Client)
	feedResource, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] GitHub repository feed (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	feedResource, err = octopusdeploy.ToFeed(feedResource.(*octopusdeploy.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	gitHubRepositoryFeed := feedResource.(*octopusdeploy.GitHubRepositoryFeed)
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

	client := m.(*octopusdeploy.Client)
	updatedFeed, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	feedResource, err := octopusdeploy.ToFeed(updatedFeed.(*octopusdeploy.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setGitHubRepositoryFeed(ctx, d, feedResource.(*octopusdeploy.GitHubRepositoryFeed)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] GitHub repository feed updated (%s)", d.Id())
	return nil
}
