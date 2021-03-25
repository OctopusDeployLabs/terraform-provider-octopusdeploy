package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandActionTemplateParameter(tfTemplate map[string]interface{}) *octopusdeploy.ActionTemplateParameter {
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

func expandActionTemplateParameters(actionTemplateParameters []interface{}) []*octopusdeploy.ActionTemplateParameter {
	if len(actionTemplateParameters) == 0 {
		return nil
	}

	expandedActionTemplateParameters := make([]*octopusdeploy.ActionTemplateParameter, len(actionTemplateParameters))
	for _, actionTemplateParameter := range actionTemplateParameters {
		actionTemplateParameterMap := actionTemplateParameter.(map[string]interface{})
		expandedActionTemplateParameters = append(expandedActionTemplateParameters, expandActionTemplateParameter(actionTemplateParameterMap))
	}
	return expandedActionTemplateParameters
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

func getActionTemplateParameterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_value": {
			Description: "A default value for the parameter, if applicable. This can be a hard-coded value or a variable reference.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"display_settings": {
			Description: "The display settings for the parameter.",
			Optional:    true,
			Type:        schema.TypeMap,
		},
		"help_text": {
			Description: "The help presented alongside the parameter input.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"id": getIDSchema(),
		"label": {
			Description: "The label shown beside the parameter when presented in the deployment process. Example: `Server name`.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"name": {
			Description:      "The name of the variable set by the parameter. The name can contain letters, digits, dashes and periods. Example: `ServerName`.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
	}
}
