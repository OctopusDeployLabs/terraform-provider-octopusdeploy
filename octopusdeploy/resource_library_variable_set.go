package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
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

	apiClient := m.(*client.Client)
	resource, err := apiClient.LibraryVariableSets.Add(libraryVariableSet)
	if err != nil {
		return createResourceOperationError(errorCreatingLibraryVariableSet, libraryVariableSet.Name, err)
	}

	if isEmpty(resource.ID) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.ID)
	}

	return nil
}

func buildLibraryVariableSetResource(d *schema.ResourceData) *model.LibraryVariableSet {
	name := d.Get(constName).(string)

	resource := model.NewLibraryVariableSet(name)

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

func buildTemplateResource(tfTemplate map[string]interface{}) model.ActionTemplateParameter {
	resource := model.ActionTemplateParameter{
		Name: tfTemplate[constName].(string),
		DisplaySettings: map[string]string{
			"Octopus.ControlType": "SingleLineText",
		},
	}

	return resource
}

func resourceLibraryVariableSetRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.LibraryVariableSets.GetByID(id)
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

	apiClient := m.(*client.Client)
	resource, err := apiClient.LibraryVariableSets.Update(*libraryVariableSet)
	if err != nil {
		return createResourceOperationError(errorUpdatingLibraryVariableSet, d.Id(), err)
	}

	d.SetId(resource.ID)

	return nil
}

func resourceLibraryVariableSetDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	err := apiClient.LibraryVariableSets.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingLibraryVariableSet, id, err)
	}

	d.SetId(constEmptyString)

	return nil
}
