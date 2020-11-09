package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceChannel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceChannelReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"description": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"project_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"lifecycle_id": {
				Optional: true,
				Type:     schema.TypeString,
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
						constVersionRange: {
							Type:     schema.TypeString,
							Optional: true,
						},
						constTag: {
							Type:     schema.TypeString,
							Optional: true,
						},
						constActions: {
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

func dataSourceChannelReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	name := d.Get("name").(string)
	query := octopusdeploy.ChannelsQuery{
		PartialName: name,
		Take:        1,
	}

	channels, err := client.Channels.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}
	if channels == nil || len(channels.Items) == 0 {
		return diag.Errorf("unable to retrieve channel (partial name: %s)", name)
	}

	// NOTE: two or more channels can have the same name in Octopus and
	// therefore, a better search criteria needs to be implemented below

	for _, channel := range channels.Items {
		if channel.Name == name {
			logResource(constChannel, m)

			d.SetId(channel.ID)
			d.Set(constName, channel.Name)
			d.Set("description", channel.Description)

			return nil
		}
	}

	return nil
}
