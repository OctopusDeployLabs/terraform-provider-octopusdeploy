package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
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
	client := m.(*octopusdeploy.Client)

	name := d.Get("name")

	libraryVariableSet, err := client.LibraryVariableSet.GetByName(name.(string))

	if err == octopusdeploy.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading libraryVariableSet name %s: %s", name, err.Error())
	}

	d.SetId(libraryVariableSet.ID)

	log.Printf("[DEBUG] libraryVariableSet: %v", m)
	d.Set("name", libraryVariableSet.Name)
	d.Set("description", libraryVariableSet.Description)
	d.Set("variable_set_id", libraryVariableSet.VariableSetId)
	return nil
}
