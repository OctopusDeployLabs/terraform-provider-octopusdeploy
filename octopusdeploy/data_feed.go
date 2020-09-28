package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataFeed() *schema.Resource {
	return &schema.Resource{
		Read: dataFeedReadByName,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataFeedReadByName(d *schema.ResourceData, m interface{}) error {
	name := d.Get(constName).(string)

	apiClient := m.(*client.Client)
	resourceList, err := apiClient.Feeds.GetByPartialName(name)
	if err != nil {
		return createResourceOperationError(errorReadingFeed, name, err)
	}
	if len(resourceList) == 0 {
		// d.SetId(constEmptyString)
		return nil
	}

	logResource(constFeed, m)

	// NOTE: two or more feeds can have the same name in Octopus and
	// therefore, a better search criteria needs to be implemented below

	for _, resource := range resourceList {
		if resource.Name == name {
			logResource(constFeed, m)

			d.SetId(resource.ID)
			d.Set(constName, resource.Name)

			return nil
		}
	}

	return nil
}
