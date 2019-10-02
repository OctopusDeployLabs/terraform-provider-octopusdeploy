package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataFeed() *schema.Resource {
	return &schema.Resource{
		Read: dataFeedReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataFeedReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	feedName := d.Get("name")

	feed, err := client.Feed.GetByName(feedName.(string))

	if err == octopusdeploy.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading feed name %s: %s", feedName, err.Error())
	}

	d.SetId(feed.ID)

	log.Printf("[DEBUG] feed: %v", m)
	d.Set("name", feed.Name)
	return nil
}
