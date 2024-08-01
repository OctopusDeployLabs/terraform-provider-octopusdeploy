package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func flattenFeed(feed *feeds.FeedResource) map[string]interface{} {
	if feed == nil {
		return nil
	}

	return map[string]interface{}{
		"access_key":                            feed.AccessKey,
		"api_version":                           feed.APIVersion,
		"delete_unreleased_packages_after_days": feed.DeleteUnreleasedPackagesAfterDays,
		"download_attempts":                     feed.DownloadAttempts,
		"download_retry_backoff_seconds":        feed.DownloadRetryBackoffSeconds,
		"feed_type":                             feed.FeedType,
		"feed_uri":                              feed.FeedURI,
		"id":                                    feed.GetID(),
		"is_enhanced_mode":                      feed.EnhancedMode,
		"name":                                  feed.Name,
		"package_acquisition_location_options":  feed.PackageAcquisitionLocationOptions,
		"region":                                feed.Region,
		"registry_path":                         feed.RegistryPath,
		"space_id":                              feed.SpaceID,
		"username":                              feed.Username,
	}
}

func getFeedDataSchema() map[string]*schema.Schema {
	dataSchema := getFeedSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"feeds": {
			Computed:    true,
			Description: "A list of feeds that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    false,
			Type:        schema.TypeList,
		},
		"feed_type":    getQueryFeedType(),
		"ids":          getQueryIDs(),
		"name":         getQueryName(),
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
		"space_id":     getSpaceIDSchema(),
	}
}

func getFeedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_key": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"api_version": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"delete_unreleased_packages_after_days": {
			Optional: true,
			Type:     schema.TypeInt,
		},
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
		"feed_type": {
			Default:  "None",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"AwsElasticContainerRegistry",
				"BuiltIn",
				"Docker",
				"GitHub",
				"Helm",
				"Maven",
				"None",
				"NuGet",
				"OctopusProject",
			}, false)),
		},
		"feed_uri": {
			Required: true,
			Type:     schema.TypeString,
		},
		"id": getIDSchema(),
		"is_enhanced_mode": {
			Default:  true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"name": {
			Description:      "A short, memorable, unique name for this feed. Example: ACME Builds.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"password": getPasswordSchema(false),
		"package_acquisition_location_options": {
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"region": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"registry_path": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"secret_key": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"space_id": getSpaceIDSchema(),
		"username": getUsernameSchema(false),
	}
}
