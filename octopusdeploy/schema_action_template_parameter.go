package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandActionTemplateParameter(values interface{}) *octopusdeploy.ActionTemplateParameter {
	flattenedValues := values.([]interface{})
	flattenedMap := flattenedValues[0].(map[string]interface{})

	return &octopusdeploy.ActionTemplateParameter{
		HelpText: flattenedMap["help_text"].(string),
		Label:    flattenedMap["label"].(string),
		Name:     flattenedMap["name"].(string),
	}
}

func expandActionTemplateParameters(actionTemplateParameters []interface{}) []*octopusdeploy.ActionTemplateParameter {
	expandedActionTemplateParameters := make([]*octopusdeploy.ActionTemplateParameter, len(actionTemplateParameters))
	for _, actionTemplateParameter := range actionTemplateParameters {
		actionTemplateParameterMap := actionTemplateParameter.(map[string]interface{})
		expandedActionTemplateParameters = append(expandedActionTemplateParameters, expandActionTemplateParameter(actionTemplateParameterMap))
	}
	return expandedActionTemplateParameters
}

func getActionTemplateParameterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_value": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"display_settings": {
			Optional: true,
			Type:     schema.TypeMap,
		},
		"help_text": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": getIDSchema(),
		"label": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"name": getNameSchema(true),
	}
}
