package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandActionTemplateParameters(actionTemplateParameters []interface{}) []*octopusdeploy.ActionTemplateParameter {
	expandedActionTemplateParameters := make([]*octopusdeploy.ActionTemplateParameter, len(actionTemplateParameters))
	for _, actionTemplateParameter := range actionTemplateParameters {
		actionTemplateParameterMap := actionTemplateParameter.(map[string]interface{})
		expandedActionTemplateParameters = append(expandedActionTemplateParameters, &octopusdeploy.ActionTemplateParameter{
			HelpText: actionTemplateParameterMap["help_text"].(string),
			Label:    actionTemplateParameterMap["label"].(string),
			Name:     actionTemplateParameterMap["name"].(string),
		})
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
			Type:     schema.TypeString,
		},
		"help_text": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"label": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"name": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
