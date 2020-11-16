package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandFeed(d *schema.ResourceData) *octopusdeploy.FeedResource {
	name := d.Get("name").(string)
	feedType := octopusdeploy.FeedType(d.Get("feed_type").(string))

	var feed = octopusdeploy.NewFeedResource(name, feedType)
	feed.ID = d.Id()

	if v, ok := d.GetOk("download_attempts"); ok {
		feed.DownloadAttempts = v.(int)
	}

	if v, ok := d.GetOk("download_retry_backoff_seconds"); ok {
		feed.DownloadRetryBackoffSeconds = v.(int)
	}

	if v, ok := d.GetOk("is_enhanced_mode"); ok {
		feed.EnhancedMode = v.(bool)
	}

	if v, ok := d.GetOk("feed_uri"); ok {
		feed.FeedURI = v.(string)
	}

	if v, ok := d.GetOk("username"); ok {
		feed.Username = v.(string)
	}

	if v, ok := d.GetOk("password"); ok {
		feed.Password = octopusdeploy.NewSensitiveValue(v.(string))
	}

	return feed
}

func flattenFeed(feed *octopusdeploy.FeedResource) map[string]interface{} {
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

func setFeed(ctx context.Context, d *schema.ResourceData, feed *octopusdeploy.FeedResource) {
	d.Set("access_key", feed.AccessKey)
	d.Set("api_version", feed.APIVersion)
	d.Set("delete_unreleased_packages_after_days", feed.DeleteUnreleasedPackagesAfterDays)
	d.Set("download_attempts", feed.DownloadAttempts)
	d.Set("download_retry_backoff_seconds", feed.DownloadRetryBackoffSeconds)
	d.Set("feed_type", feed.FeedType)
	d.Set("feed_uri", feed.FeedURI)
	d.Set("is_enhanced_mode", feed.EnhancedMode)
	d.Set("name", feed.Name)
	d.Set("package_acquisition_location_options", feed.PackageAcquisitionLocationOptions)
	d.Set("region", feed.Region)
	d.Set("registry_path", feed.RegistryPath)
	d.Set("space_id", feed.SpaceID)
	d.Set("username", feed.Username)

	d.SetId(feed.GetID())
}

func getFeedDataSchema() map[string]*schema.Schema {
	feedSchema := getFeedSchema()
	for _, field := range feedSchema {
		field.Computed = true
		field.Default = nil
		field.MaxItems = 0
		field.MinItems = 0
		field.Optional = false
		field.Required = false
		field.ValidateDiagFunc = nil
		field.ValidateFunc = nil
	}

	return map[string]*schema.Schema{
		"feeds": {
			Computed: true,
			Elem:     &schema.Resource{Schema: feedSchema},
			Type:     schema.TypeList,
		},
		"feed_type": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"ids": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"partial_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"skip": {
			Default:  0,
			Type:     schema.TypeInt,
			Optional: true,
		},
		"take": {
			Default:  1,
			Type:     schema.TypeInt,
			Optional: true,
		},
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
			Default:  5,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"download_retry_backoff_seconds": {
			Default:  10,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"feed_type": {
			Default:  "None",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validateValueFunc([]string{
				"AwsElasticContainerRegistry",
				"BuiltIn",
				"Docker",
				"GitHub",
				"Helm",
				"Maven",
				"None",
				"NuGet",
				"OctopusProject",
			}),
		},
		"feed_uri": {
			Required: true,
			Type:     schema.TypeString,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"is_enhanced_mode": {
			Default:  true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"password": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"package_acquisition_location_options": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Required: true,
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
		"space_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"username": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
	}
}
