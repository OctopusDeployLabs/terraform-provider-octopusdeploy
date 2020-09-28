package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataLifecycle() *schema.Resource {
	return &schema.Resource{
		Read: dataLifecycleReadByName,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataLifecycleReadByName(d *schema.ResourceData, m interface{}) error {
	name := d.Get(constName).(string)

	apiClient := m.(*client.Client)
	resourceList, err := apiClient.Lifecycles.GetByPartialName(name)
	if err != nil {
		return createResourceOperationError(errorReadingLifecycle, name, err)
	}
	if len(resourceList) == 0 {
		// d.SetId(constEmptyString)
		return nil
	}

	logResource(constLifecycle, m)

	// NOTE: two or more lifecycles can have the same name in Octopus and
	// therefore, a better search criteria needs to be implemented below

	for _, resource := range resourceList {
		if resource.Name == name {
			logResource(constLifecycle, m)

			d.SetId(resource.ID)
			d.Set(constName, resource.Name)
			d.Set(constDescription, resource.Description)

			return nil
		}
	}

	return nil
}
