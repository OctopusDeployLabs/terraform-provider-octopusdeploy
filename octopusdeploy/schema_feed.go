package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandFeedResource(d *schema.ResourceData) *octopusdeploy.FeedResource {
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

func flattenFeedResource(ctx context.Context, d *schema.ResourceData, feed *octopusdeploy.FeedResource) {
	d.Set("download_attempts", feed.DownloadAttempts)
	d.Set("download_retry_backoff_seconds", feed.DownloadRetryBackoffSeconds)
	d.Set("feed_type", feed.FeedType)
	d.Set("feed_uri", feed.FeedURI)
	d.Set("is_enhanced_mode", feed.EnhancedMode)
	d.Set("name", feed.Name)
	// d.Set("password", feed.Password)
	d.Set("username", feed.Username)

	d.SetId(feed.GetID())
}

func getFeedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"is_enhanced_mode": {
			Default:  true,
			Optional: true,
			Type:     schema.TypeBool,
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
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"password": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"username": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
	}
}
