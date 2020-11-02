package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataChannel() *schema.Resource {
	return &schema.Resource{
		Read: dataChannelReadByName,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constProjectID: {
				Type:     schema.TypeString,
				Required: true,
			},
			constLifecycleID: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constIsDefault: {
				Type:     schema.TypeBool,
				Optional: true,
			},
			constRule: {
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

func dataChannelReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)
	name := d.Get(constName).(string)
	query := octopusdeploy.ChannelsQuery{
		PartialName: name,
		Take:        1,
	}

	channels, err := client.Channels.Get(query)
	if err != nil {
		return createResourceOperationError(errorReadingChannel, name, err)
	}
	if channels == nil || len(channels.Items) == 0 {
		return fmt.Errorf("Unabled to retrieve channel (partial name: %s)", name)
	}

	logResource(constChannel, m)

	// NOTE: two or more channels can have the same name in Octopus and
	// therefore, a better search criteria needs to be implemented below

	for _, channel := range channels.Items {
		if channel.Name == name {
			logResource(constChannel, m)

			d.SetId(channel.ID)
			d.Set(constName, channel.Name)
			d.Set(constDescription, channel.Description)

			return nil
		}
	}

	return nil
}
