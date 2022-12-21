package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandMavenFeed(d *schema.ResourceData) (*feeds.MavenFeed, error) {
	name := d.Get("name").(string)

	feed, err := feeds.NewMavenFeed(name)
	if err != nil {
		return nil, err
	}

	feed.ID = d.Id()

	if v, ok := d.GetOk("download_attempts"); ok {
		feed.DownloadAttempts = v.(int)
	}

	if v, ok := d.GetOk("download_retry_backoff_seconds"); ok {
		feed.DownloadRetryBackoffSeconds = v.(int)
	}

	if v, ok := d.GetOk("feed_uri"); ok {
		feed.FeedURI = v.(string)
	}

	if v, ok := d.GetOk("package_acquisition_location_options"); ok {
		feed.PackageAcquisitionLocationOptions = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("password"); ok {
		feed.Password = core.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("space_id"); ok {
		feed.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("username"); ok {
		feed.Username = v.(string)
	}

	return feed, nil
}

func getMavenFeedSchema() map[string]*schema.Schema {
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

func setMavenFeed(ctx context.Context, d *schema.ResourceData, feed *feeds.MavenFeed) error {
	d.Set("download_attempts", feed.DownloadAttempts)
	d.Set("download_retry_backoff_seconds", feed.DownloadRetryBackoffSeconds)
	d.Set("feed_uri", feed.FeedURI)
	d.Set("name", feed.Name)
	d.Set("space_id", feed.SpaceID)
	d.Set("username", feed.Username)

	if err := d.Set("package_acquisition_location_options", feed.PackageAcquisitionLocationOptions); err != nil {
		return fmt.Errorf("error setting package_acquisition_location_options: %s", err)
	}

	d.SetId(feed.GetID())

	return nil
}
