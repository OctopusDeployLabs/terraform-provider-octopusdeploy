package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandMavenFeed(d *schema.ResourceData) *octopusdeploy.MavenFeed {
	name := d.Get("name").(string)

	var mavenFeed = octopusdeploy.NewMavenFeed(name)
	mavenFeed.ID = d.Id()

	if v, ok := d.GetOk("download_attempts"); ok {
		mavenFeed.DownloadAttempts = v.(int)
	}

	if v, ok := d.GetOk("download_retry_backoff_seconds"); ok {
		mavenFeed.DownloadRetryBackoffSeconds = v.(int)
	}

	if v, ok := d.GetOk("feed_uri"); ok {
		mavenFeed.FeedURI = v.(string)
	}

	if v, ok := d.GetOk("package_acquisition_location_options"); ok {
		mavenFeed.PackageAcquisitionLocationOptions = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("password"); ok {
		mavenFeed.Password = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("username"); ok {
		mavenFeed.Username = v.(string)
	}

	return mavenFeed
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

func setMavenFeed(ctx context.Context, d *schema.ResourceData, mavenFeed *octopusdeploy.MavenFeed) error {
	d.Set("download_attempts", mavenFeed.DownloadAttempts)
	d.Set("download_retry_backoff_seconds", mavenFeed.DownloadRetryBackoffSeconds)
	d.Set("feed_uri", mavenFeed.FeedURI)
	d.Set("name", mavenFeed.Name)
	d.Set("space_id", mavenFeed.SpaceID)
	d.Set("username", mavenFeed.Username)

	if err := d.Set("package_acquisition_location_options", mavenFeed.PackageAcquisitionLocationOptions); err != nil {
		return fmt.Errorf("error setting package_acquisition_location_options: %s", err)
	}

	d.SetId(mavenFeed.GetID())

	return nil
}
