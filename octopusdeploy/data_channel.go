package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	apiClient := m.(*client.Client)

	name := d.Get(constName).(string)
	resourceList, err := apiClient.Channels.GetByPartialName(name)

	if err != nil {
		return createResourceOperationError(errorReadingChannel, name, err)
	}
	if len(resourceList) == 0 {
		// d.SetId(constEmptyString)
		return nil
	}

	// NOTE: two or more channels can have the same name in Octopus and
	// therefore, a better search criteria needs to be implemented below

	for _, resource := range resourceList {
		if resource.Name == name {
			logResource(constChannel, m)

			d.SetId(resource.ID)
			d.Set(constName, resource.Name)
			d.Set(constDescription, resource.Description)

			return nil
		}
	}

	return nil
}
