package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataNuget() *schema.Resource {
	return &schema.Resource{
		Read: dataNugetReadByName,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constFeedURI: {
				Type:     schema.TypeString,
				Required: true,
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
				Type:     schema.TypeString,
				Optional: true,
			},
			constPassword: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func dataNugetReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)
	name := d.Get(constName).(string)
	query := octopusdeploy.FeedsQuery{
		PartialName: name,
		Take:        1,
	}

	feeds, err := client.Feeds.Get(query)
	if err != nil {
		return createResourceOperationError(errorReadingNuGetFeed, name, err)
	}
	if feeds == nil || len(feeds.Items) == 0 {
		return fmt.Errorf("Unabled to retrieve feed (partial name: %s)", name)
	}

	logResource(constFeed, m)

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
