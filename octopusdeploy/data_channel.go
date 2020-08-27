package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataChannel() *schema.Resource {
	return &schema.Resource{
		Read: dataChannelReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lifecycle_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_range": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tag": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"actions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataChannelReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	ChannelName := d.Get("name")
	env, err := client.Channel.Get(ChannelName.(string))

	if err == octopusdeploy.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Channel with name %s: %s", ChannelName, err.Error())
	}

	d.SetId(env.ID)

	d.Set("name", env.Name)
	d.Set("description", env.Description)

	return nil
}
