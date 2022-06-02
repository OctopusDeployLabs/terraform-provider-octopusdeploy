package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandNuGetFeed(d *schema.ResourceData) (*octopusdeploy.NuGetFeed, error) {
	name := d.Get("name").(string)
	feedURI := d.Get("feed_uri").(string)

	nuGetFeed, err := octopusdeploy.NewNuGetFeed(name, feedURI)
	if err != nil {
		return nil, err
	}

	nuGetFeed.ID = d.Id()

	if v, ok := d.GetOk("download_attempts"); ok {
		nuGetFeed.DownloadAttempts = v.(int)
	}

	if v, ok := d.GetOk("download_retry_backoff_seconds"); ok {
		nuGetFeed.DownloadRetryBackoffSeconds = v.(int)
	}

	if v, ok := d.GetOk("is_enhanced_mode"); ok {
		nuGetFeed.EnhancedMode = v.(bool)
	}

	if v, ok := d.GetOk("package_acquisition_location_options"); ok {
		nuGetFeed.PackageAcquisitionLocationOptions = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("password"); ok {
		nuGetFeed.Password = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("username"); ok {
		nuGetFeed.Username = v.(string)
	}

	return nuGetFeed, nil
}

func getNuGetFeedSchema() map[string]*schema.Schema {
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
			Description: "The feed URI can be a URL or a folder path.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"id": getIDSchema(),
		"is_enhanced_mode": {
			Default:     true,
			Description: "This will improve performance of the NuGet feed but may not be supported by some older feeds. Disable if the operation, Create Release does not return the latest version for a package.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
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

func setNuGetFeed(ctx context.Context, d *schema.ResourceData, nuGetFeed *octopusdeploy.NuGetFeed) error {
	d.Set("download_attempts", nuGetFeed.DownloadAttempts)
	d.Set("download_retry_backoff_seconds", nuGetFeed.DownloadRetryBackoffSeconds)
	d.Set("feed_uri", nuGetFeed.FeedURI)
	d.Set("is_enhanced_mode", nuGetFeed.EnhancedMode)
	d.Set("name", nuGetFeed.Name)
	d.Set("username", nuGetFeed.Username)

	if err := d.Set("package_acquisition_location_options", nuGetFeed.PackageAcquisitionLocationOptions); err != nil {
		return fmt.Errorf("error setting package_acquisition_location_options: %s", err)
	}

	d.SetId(nuGetFeed.GetID())

	return nil
}
