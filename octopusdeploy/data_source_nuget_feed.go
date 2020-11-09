package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNuGetFeed() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNuGetFeedReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			constFeedURI: {
				Required: true,
				Type:     schema.TypeString,
			},
			constEnhancedMode: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			constDownloadAttempts: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			constDownloadRetryBackoffSeconds: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},
			constUsername: {
				Optional: true,
				Type:     schema.TypeString,
			},
			constPassword: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceNuGetFeedReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	name := d.Get("name").(string)
	query := octopusdeploy.FeedsQuery{
		PartialName: name,
		Take:        1,
	}

	feeds, err := client.Feeds.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}
	if feeds == nil || len(feeds.Items) == 0 {
		return diag.Errorf("unable to retrieve feed (partial name: %s)", name)
	}

	// NOTE: two or more feeds can have the same name in Octopus and
	// therefore, a better search criteria needs to be implemented below

	for _, feed := range feeds.Items {
		if feed.GetName() == name {
			logResource(constFeed, m)

			d.SetId(feed.GetID())
			d.Set(constName, feed.GetName())

			return nil
		}
	}

	return nil
}
