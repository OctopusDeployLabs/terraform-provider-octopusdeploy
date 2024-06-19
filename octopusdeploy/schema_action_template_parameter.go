package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandActionTemplateParameter(tfTemplate map[string]interface{}) actiontemplates.ActionTemplateParameter {
	actionTemplateParameter := actiontemplates.NewActionTemplateParameter()

	propertyValue := core.NewPropertyValue(tfTemplate["default_value"].(string), false)
	actionTemplateParameter.DefaultValue = &propertyValue
	actionTemplateParameter.DisplaySettings = flattenDisplaySettings(tfTemplate["display_settings"].(map[string]interface{}))
	actionTemplateParameter.HelpText = tfTemplate["help_text"].(string)
	actionTemplateParameter.ID = tfTemplate["id"].(string)
	actionTemplateParameter.Label = tfTemplate["label"].(string)
	actionTemplateParameter.Name = tfTemplate["name"].(string)

	return *actionTemplateParameter
}

func expandActionTemplateParameters(actionTemplateParameters []interface{}) []actiontemplates.ActionTemplateParameter {
	if len(actionTemplateParameters) == 0 {
		return nil
	}

	expandedActionTemplateParameters := []actiontemplates.ActionTemplateParameter{}
	for _, actionTemplateParameter := range actionTemplateParameters {
		actionTemplateParameterMap := actionTemplateParameter.(map[string]interface{})
		expandedActionTemplateParameters = append(expandedActionTemplateParameters, expandActionTemplateParameter(actionTemplateParameterMap))
	}
	return expandedActionTemplateParameters
}

func flattenActionTemplateParameters(actionTemplateParameters []actiontemplates.ActionTemplateParameter) []interface{} {
	flattenedActionTemplateParameters := make([]interface{}, 0)
	for _, actionTemplateParameter := range actionTemplateParameters {
		a := make(map[string]interface{})
		a["default_value"] = actionTemplateParameter.DefaultValue.Value
		a["display_settings"] = actionTemplateParameter.DisplaySettings
		a["help_text"] = actionTemplateParameter.HelpText
		a["id"] = actionTemplateParameter.ID
		a["label"] = actionTemplateParameter.Label
		a["name"] = actionTemplateParameter.Name
		flattenedActionTemplateParameters = append(flattenedActionTemplateParameters, a)
	}
	return flattenedActionTemplateParameters
}

func mapTemplateNamesToIds(actionTemplateParameters []actiontemplates.ActionTemplateParameter) map[string]string {
	templateNameIds := map[string]string{}
	for _, actionTemplateParameter := range actionTemplateParameters {
		templateNameIds[actionTemplateParameter.Name] = actionTemplateParameter.ID
	}
	return templateNameIds
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
