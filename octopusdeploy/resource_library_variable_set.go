package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLibraryVariableSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceLibraryVariableSetCreate,
		Read:   resourceLibraryVariableSetRead,
		Update: resourceLibraryVariableSetUpdate,
		Delete: resourceLibraryVariableSetDelete,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constTemplates: getTemplatesSchema(),
		},
	}
}

func getTemplatesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				constName: {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func resourceLibraryVariableSetCreate(d *schema.ResourceData, m interface{}) error {
	libraryVariableSet := buildLibraryVariableSetResource(d)

	client := m.(*octopusdeploy.Client)
	resource, err := client.LibraryVariableSets.Add(libraryVariableSet)
	if err != nil {
		return createResourceOperationError(errorCreatingLibraryVariableSet, libraryVariableSet.Name, err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func buildLibraryVariableSetResource(d *schema.ResourceData) *octopusdeploy.LibraryVariableSet {
	name := d.Get(constName).(string)

	resource := octopusdeploy.NewLibraryVariableSet(name)

	if attr, ok := d.GetOk(constDescription); ok {
		resource.Description = attr.(string)
	}

	if attr, ok := d.GetOk(constTemplates); ok {
		tfTemplates := attr.([]interface{})

		for _, tfTemplate := range tfTemplates {
			template := buildTemplateResource(tfTemplate.(map[string]interface{}))
			resource.Templates = append(resource.Templates, &template)
		}
	}

	return resource
}

func buildTemplateResource(tfTemplate map[string]interface{}) octopusdeploy.ActionTemplateParameter {
	resource := octopusdeploy.ActionTemplateParameter{
		Name: tfTemplate[constName].(string),
		DisplaySettings: map[string]string{
			"Octopus.ControlType": "SingleLineText",
		},
	}

	return resource
}

func resourceLibraryVariableSetRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	resource, err := client.LibraryVariableSets.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingLibraryVariableSet, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constLibraryVariableSet, m)

	d.Set(constName, resource.Name)
	d.Set(constDescription, resource.Description)
	d.Set(constVariableSetID, resource.VariableSetID)

	return nil
}

func resourceLibraryVariableSetUpdate(d *schema.ResourceData, m interface{}) error {
	libraryVariableSet := buildLibraryVariableSetResource(d)
	libraryVariableSet.ID = d.Id() // set ID so Octopus API knows which library variable set to update

	client := m.(*octopusdeploy.Client)
	resource, err := client.LibraryVariableSets.Update(libraryVariableSet)
	if err != nil {
		return createResourceOperationError(errorUpdatingLibraryVariableSet, d.Id(), err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceLibraryVariableSetDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.LibraryVariableSets.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingLibraryVariableSet, id, err)
	}

	d.SetId(constEmptyString)

	return nil
}
