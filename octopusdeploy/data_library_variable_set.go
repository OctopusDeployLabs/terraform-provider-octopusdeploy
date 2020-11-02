package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataLibraryVariableSet() *schema.Resource {
	return &schema.Resource{
		Read: dataLibraryVariableSetReadByName,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataLibraryVariableSetReadByName(d *schema.ResourceData, m interface{}) error {
	name := d.Get(constName).(string)

	client := m.(*octopusdeploy.Client)
	resourceList, err := client.LibraryVariableSets.GetByPartialName(name)
	if err != nil {
		return createResourceOperationError(errorReadingLibraryVariableSet, name, err)
	}
	if len(resourceList) == 0 {
		return nil
	}

	logResource(constLibraryVariableSet, m)

	// NOTE: two or more library variables can have the same name in Octopus.
	// Therefore, a better search criteria needs to be implemented below.

	for _, resource := range resourceList {
		if resource.Name == name {
			logResource(constLibraryVariableSet, m)

			d.SetId(resource.GetID())
			d.Set(constName, resource.Name)
			d.Set(constDescription, resource.Description)
			d.Set(constVariableSetID, resource.VariableSetID)

			return nil
		}
	}

	return nil
}
