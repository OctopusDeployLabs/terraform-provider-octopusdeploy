package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandGitHubRepositoryFeed(d *schema.ResourceData) *octopusdeploy.GitHubRepositoryFeed {
	name := d.Get("name").(string)

	var gitHubRepositoryFeed = octopusdeploy.NewGitHubRepositoryFeed(name)
	gitHubRepositoryFeed.ID = d.Id()

	if v, ok := d.GetOk("download_attempts"); ok {
		gitHubRepositoryFeed.DownloadAttempts = v.(int)
	}

	if v, ok := d.GetOk("download_retry_backoff_seconds"); ok {
		gitHubRepositoryFeed.DownloadRetryBackoffSeconds = v.(int)
	}

	if v, ok := d.GetOk("feed_uri"); ok {
		gitHubRepositoryFeed.FeedURI = v.(string)
	}

	if v, ok := d.GetOk("package_acquisition_location_options"); ok {
		gitHubRepositoryFeed.PackageAcquisitionLocationOptions = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("password"); ok {
		gitHubRepositoryFeed.Password = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("username"); ok {
		gitHubRepositoryFeed.Username = v.(string)
	}

	return gitHubRepositoryFeed
}

func getGitHubRepositoryFeedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"download_attempts": {
			Default:     5,
			Description: "The number of times a deployment should attempt to download a package from this feed before failing.",
			Optional:    true,
			Type:        schema.TypeInt,
		},
		"download_retry_backoff_seconds": {
			Default:     10,
			Description: "The number of seconds to apply as a linear back off between download attempts.",
			Optional:    true,
			Type:        schema.TypeInt,
		},
		"feed_uri": {
			Required: true,
			Type:     schema.TypeString,
		},
		"id": getIDSchema(),
		"name": {
			Description:      "A short, memorable, unique name for this feed. Example: ACME Builds.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"package_acquisition_location_options": {
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"password": getPasswordSchema(false),
		"space_id": getSpaceIDSchema(),
		"username": getUsernameSchema(false),
	}
}

func setGitHubRepositoryFeed(ctx context.Context, d *schema.ResourceData, githubRepositoryFeed *octopusdeploy.GitHubRepositoryFeed) error {
	d.Set("download_attempts", githubRepositoryFeed.DownloadAttempts)
	d.Set("download_retry_backoff_seconds", githubRepositoryFeed.DownloadRetryBackoffSeconds)
	d.Set("feed_uri", githubRepositoryFeed.FeedURI)
	d.Set("name", githubRepositoryFeed.Name)
	d.Set("space_id", githubRepositoryFeed.SpaceID)
	d.Set("username", githubRepositoryFeed.Username)

	if err := d.Set("package_acquisition_location_options", githubRepositoryFeed.PackageAcquisitionLocationOptions); err != nil {
		return fmt.Errorf("error setting package_acquisition_location_options: %s", err)
	}

	d.SetId(githubRepositoryFeed.GetID())

	return nil
}
