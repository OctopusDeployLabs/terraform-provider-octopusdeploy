package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	if d == nil {
		return createInvalidParameterError("resourceLibraryVariableSetCreate", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceLibraryVariableSetCreate", "m")
	}

	apiClient := m.(*client.Client)

	newLibraryVariableSet := buildLibraryVariableSetResource(d)

	createdLibraryVariableSet, err := apiClient.LibraryVariableSets.Add(newLibraryVariableSet)

	if err != nil {
		return createResourceOperationError(errorCreatingLibraryVariableSet, newLibraryVariableSet.Name, err)
	}

	d.SetId(createdLibraryVariableSet.ID)

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
	apiClient := m.(*client.Client)

	id := d.Id()
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
	resource := buildLibraryVariableSetResource(d)
	resource.ID = d.Id() // set libraryVariableSet struct ID so octopus knows which libraryVariableSet to update

	apiClient := m.(*client.Client)
	updatedResource, err := apiClient.LibraryVariableSets.Update(*resource)

	if err != nil {
		return createResourceOperationError(errorUpdatingLibraryVariableSet, d.Id(), err)
	}

	d.SetId(updatedResource.ID)

	return nil
}

func resourceLibraryVariableSetDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	id := d.Id()

	err := apiClient.LibraryVariableSets.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingLibraryVariableSet, id, err)
	}

	d.SetId(constEmptyString)
	return nil
}
