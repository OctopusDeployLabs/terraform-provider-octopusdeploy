package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataLibraryVariableSet() *schema.Resource {
	return &schema.Resource{
		Read: dataLibraryVariableSetReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataLibraryVariableSetReadByName(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	name := d.Get("name")

	libraryVariableSet, err := apiClient.LibraryVariableSets.GetByName(name.(string))

	if err == client.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading libraryVariableSet name %s: %s", name, err.Error())
	}

	d.SetId(libraryVariableSet.ID)

	log.Printf("[DEBUG] libraryVariableSet: %v", m)
	d.Set("name", libraryVariableSet.Name)
	d.Set("description", libraryVariableSet.Description)
	d.Set("variable_set_id", libraryVariableSet.VariableSetID)
	return nil
}
