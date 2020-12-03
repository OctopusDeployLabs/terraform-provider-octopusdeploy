package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandLibraryVariableSet(d *schema.ResourceData) *octopusdeploy.LibraryVariableSet {
	name := d.Get("name").(string)

	libraryVariableSet := octopusdeploy.NewLibraryVariableSet(name)
	libraryVariableSet.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		libraryVariableSet.Description = v.(string)
	}

	if attr, ok := d.GetOk("template"); ok {
		tfTemplates := attr.([]interface{})

		for _, tfTemplate := range tfTemplates {
			template := expandTemplate(tfTemplate.(map[string]interface{}))
			libraryVariableSet.Templates = append(libraryVariableSet.Templates, template)
		}
	}

	return libraryVariableSet
}

func expandTemplate(tfTemplate map[string]interface{}) *octopusdeploy.ActionTemplateParameter {
	actionTemplateParameter := octopusdeploy.NewActionTemplateParameter()

	propertyValue := octopusdeploy.PropertyValue(tfTemplate["default_value"].(string))
	actionTemplateParameter.DefaultValue = &octopusdeploy.PropertyValueResource{
		PropertyValue: &propertyValue,
	}
	actionTemplateParameter.DisplaySettings = flattenDisplaySettings(tfTemplate["display_settings"].(map[string]interface{}))
	actionTemplateParameter.HelpText = tfTemplate["help_text"].(string)
	actionTemplateParameter.ID = tfTemplate["id"].(string)
	actionTemplateParameter.Label = tfTemplate["label"].(string)
	actionTemplateParameter.Name = tfTemplate["name"].(string)

	return actionTemplateParameter
}

func flattenActionTemplateParameters(actionTemplateParameters []*octopusdeploy.ActionTemplateParameter) []interface{} {
	flattenedActionTemplateParameters := make([]interface{}, 0)
	for _, actionTemplateParameter := range actionTemplateParameters {
		a := make(map[string]interface{})
		a["default_value"] = actionTemplateParameter.DefaultValue.PropertyValue
		a["display_settings"] = actionTemplateParameter.DisplaySettings
		a["help_text"] = actionTemplateParameter.HelpText
		a["id"] = actionTemplateParameter.ID
		a["label"] = actionTemplateParameter.Label
		a["name"] = actionTemplateParameter.Name
		flattenedActionTemplateParameters = append(flattenedActionTemplateParameters, a)
	}
	return flattenedActionTemplateParameters
}

func flattenDisplaySettings(displaySettings map[string]interface{}) map[string]string {
	flattenedDisplaySettings := make(map[string]string, len(displaySettings))
	for key, displaySetting := range displaySettings {
		flattenedDisplaySettings[key] = displaySetting.(string)
	}
	return flattenedDisplaySettings
}

func setLibraryVariableSet(ctx context.Context, d *schema.ResourceData, libraryVariableSet *octopusdeploy.LibraryVariableSet) {
	d.Set("description", libraryVariableSet.Description)
	d.Set("name", libraryVariableSet.Name)
	d.Set("space_id", libraryVariableSet.SpaceID)
	d.Set("template", flattenActionTemplateParameters(libraryVariableSet.Templates))
	d.Set("variable_set_id", libraryVariableSet.VariableSetID)

	d.SetId(libraryVariableSet.GetID())
}

func getLibraryVariableSetDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": getNameSchema(true),
	}
}

func getLibraryVariableSetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": getDescriptionSchema(),
		"name":        getNameSchema(true),
		"space_id":    getSpaceIDSchema(),
		"template": {
			Optional: true,
			Elem:     &schema.Resource{Schema: getActionTemplateParameterSchema()},
			Type:     schema.TypeList,
		},
		"variable_set_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
	}
}
