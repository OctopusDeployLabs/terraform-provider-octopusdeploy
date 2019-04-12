package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceLibraryVariableSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceLibraryVariableSetCreate,
		Read:   resourceLibraryVariableSetRead,
		Update: resourceLibraryVariableSetUpdate,
		Delete: resourceLibraryVariableSetDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": getTemplatesSchema(),
		},
	}
}

func getTemplatesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func resourceLibraryVariableSetCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newLibraryVariableSet := buildLibraryVariableSetResource(d)

	createdLibraryVariableSet, err := client.LibraryVariableSet.Add(newLibraryVariableSet)

	if err != nil {
		return fmt.Errorf("error creating project: %s", err.Error())
	}

	d.SetId(createdLibraryVariableSet.ID)

	return nil
}

func buildLibraryVariableSetResource(d *schema.ResourceData) *octopusdeploy.LibraryVariableSet {
	name := d.Get("name").(string)

	libraryVariableSet := octopusdeploy.NewLibraryVariableSet(name)

	if attr, ok := d.GetOk("description"); ok {
		libraryVariableSet.Description = attr.(string)
	}

	if attr, ok := d.GetOk("templates"); ok {
		tfTemplates := attr.([]interface{})

		for _, tfTemplate := range tfTemplates {
			template := buildTemplateResource(tfTemplate.(map[string]interface{}))
			libraryVariableSet.Templates = append(libraryVariableSet.Templates, template)
		}
	}

	return libraryVariableSet
}

func buildTemplateResource(tfTemplate map[string]interface{}) octopusdeploy.ActionTemplateParameter {
	template := octopusdeploy.ActionTemplateParameter{
		Name: tfTemplate["name"].(string),
		DisplaySettings: map[string]string{
			"Octopus.ControlType": "SingleLineText",
		},
	}

	return template
}

func resourceLibraryVariableSetRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	libraryVariableSetID := d.Id()

	libraryVariableSet, err := client.LibraryVariableSet.Get(libraryVariableSetID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading libraryVariableSet id %s: %s", libraryVariableSetID, err.Error())
	}

	log.Printf("[DEBUG] libraryVariableSet: %v", m)
	d.Set("name", libraryVariableSet.Name)
	d.Set("description", libraryVariableSet.Description)
	d.Set("variable_set_id", libraryVariableSet.VariableSetId)

	return nil
}

func resourceLibraryVariableSetUpdate(d *schema.ResourceData, m interface{}) error {
	libraryVariableSet := buildLibraryVariableSetResource(d)
	libraryVariableSet.ID = d.Id() // set libraryVariableSet struct ID so octopus knows which libraryVariableSet to update

	client := m.(*octopusdeploy.Client)

	libraryVariableSet, err := client.LibraryVariableSet.Update(libraryVariableSet)

	if err != nil {
		return fmt.Errorf("error updating libraryVariableSet id %s: %s", d.Id(), err.Error())
	}

	d.SetId(libraryVariableSet.ID)

	return nil
}

func resourceLibraryVariableSetDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	libraryVariableSetID := d.Id()

	err := client.LibraryVariableSet.Delete(libraryVariableSetID)

	if err != nil {
		return fmt.Errorf("error deleting libraryVariableSet id %s: %s", libraryVariableSetID, err.Error())
	}

	d.SetId("")
	return nil
}
