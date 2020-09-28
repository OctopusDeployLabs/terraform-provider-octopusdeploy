package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataMachinePolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataMachinePolicyReadByName,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"isdefault": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataMachinePolicyReadByName(d *schema.ResourceData, m interface{}) error {
	name := d.Get(constName).(string)

	apiClient := m.(*client.Client)
	resourceList, err := apiClient.MachinePolicies.GetAll()
	if err != nil {
		return createResourceOperationError(errorReadingMachinePolicy, name, err)
	}
	if len(resourceList) == 0 {
		return nil
	}

	logResource(constMachinePolicy, m)

	// NOTE: two or more machine policies could have the same name in Octopus
	// and therefore, a better search criteria needs to be implemented below

	for _, resource := range resourceList {
		if resource.Name == name {
			logResource(constMachinePolicy, m)

			d.SetId(resource.ID)
			d.Set(constDescription, resource.Description)
			d.Set("isdefault", resource.IsDefault)

			return nil
		}
	}

	return nil
}
