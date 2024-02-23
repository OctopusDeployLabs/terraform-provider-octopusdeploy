package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
)

func expandPromptedDisplaySettings(values interface{}) *resources.DisplaySettings {
	if values == nil {
		return nil
	}

	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 {
		return nil
	}

	promptedDisplaySettings := flattenedValues[0].(map[string]interface{})

	controlType := resources.ControlType(promptedDisplaySettings["control_type"].(string))

	var selectOptions []*resources.SelectOption
	if controlType == resources.ControlTypeSelect {
		selectOptions = expandSelectOptions(promptedDisplaySettings["select_option"])
	}

	return resources.NewDisplaySettings(controlType, selectOptions)
}

func expandPromptedVariableSettings(values interface{}) *variables.VariablePromptOptions {
	if values == nil {
		return nil
	}

	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 {
		return nil
	}

	promptedVariableSettings := flattenedValues[0].(map[string]interface{})
	return &variables.VariablePromptOptions{
		Description:     promptedVariableSettings["description"].(string),
		DisplaySettings: expandPromptedDisplaySettings(promptedVariableSettings["display_settings"]),
		IsRequired:      promptedVariableSettings["is_required"].(bool),
		Label:           promptedVariableSettings["label"].(string),
	}
}

func expandSelectOptions(values interface{}) []*resources.SelectOption {
	if values == nil {
		return nil
	}

	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 {
		return nil
	}

	selectOptions := make([]*resources.SelectOption, len(flattenedValues))

	for i := 0; i < len(flattenedValues); i++ {
		item := flattenedValues[i].(map[string]interface{})
		selectOptions[i] = &resources.SelectOption{
			DisplayName: item["display_name"].(string),
			Value:       item["value"].(string),
		}
	}

	return selectOptions
}

func flattenPromptedVariableDisplaySettings(displaySettings *resources.DisplaySettings) []interface{} {
	if displaySettings == nil {
		return nil
	}

	flattenedDisplaySettings := map[string]interface{}{}
	flattenedDisplaySettings["control_type"] = displaySettings.ControlType
	if displaySettings.ControlType == resources.ControlTypeSelect {
		flattenedDisplaySettings["select_option"] = flattenSelectOptions(displaySettings.SelectOptions)
	}
	return []interface{}{flattenedDisplaySettings}
}

func flattenPromptedVariableSettings(variablePromptOptions *variables.VariablePromptOptions) []interface{} {
	if variablePromptOptions == nil {
		return nil
	}

	flattenedPrompt := map[string]interface{}{}
	flattenedPrompt["description"] = variablePromptOptions.Description
	flattenedPrompt["is_required"] = variablePromptOptions.IsRequired
	flattenedPrompt["label"] = variablePromptOptions.Label

	if variablePromptOptions.DisplaySettings != nil {
		flattenedPrompt["display_settings"] = flattenPromptedVariableDisplaySettings(variablePromptOptions.DisplaySettings)

	}

	return []interface{}{flattenedPrompt}
}

func flattenSelectOptions(selectOptions []*resources.SelectOption) []map[string]interface{} {
	options := make([]map[string]interface{}, len(selectOptions))
	for i := 0; i < len(selectOptions); i++ {
		options[i] = map[string]interface{}{
			"value":        selectOptions[i].Value,
			"display_name": selectOptions[i].DisplayName,
		}
	}
	return options
}
