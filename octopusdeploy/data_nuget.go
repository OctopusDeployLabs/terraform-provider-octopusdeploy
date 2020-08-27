package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataNuget() *schema.Resource {
	return &schema.Resource{
		Read: dataNugetReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"feed_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enhanced_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"download_attempts": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"download_retry_backoff_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func dataNugetReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	NugetName := d.Get("name")
	env, err := client.Feed.GetByName(NugetName.(string))

	if err == octopusdeploy.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Nuget with name %s: %s", NugetName, err.Error())
	}

	d.SetId(env.ID)

	d.Set("name", env.Name)

	return nil
}
