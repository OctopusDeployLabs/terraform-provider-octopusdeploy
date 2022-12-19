package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandPromptedVariableSettings(v interface{}) *variables.VariablePromptOptions {
	if v == nil {
		return nil
	}
	s := v.([]interface{})
	tfPromptList := s[0].(map[string]interface{})

	newPrompt := variables.VariablePromptOptions{
		Description:     tfPromptList["description"].(string),
		Label:           tfPromptList["label"].(string),
		IsRequired:      tfPromptList["is_required"].(bool),
		DisplaySettings: expandPromptedDisplaySettings(tfPromptList["display_settings"]),
	}
	return &newPrompt
}

func expandPromptedDisplaySettings(v interface{}) *variables.DisplaySettings {
	if v == nil {
		return nil
	}
	s := v.([]interface{})
	tfPromptList := s[0].(map[string]interface{})

	controlType := variables.ControlType(tfPromptList["control_type"].(string))
	var selectOptions []*variables.SelectOption

	if controlType == variables.ControlTypeSelect {
		selectOptions = expandSelectOptions(tfPromptList["select_option"])
	}

	settings := variables.NewDisplaySettings(controlType, selectOptions)
	return settings
}

func expandSelectOptions(tfPromptList interface{}) []*variables.SelectOption {
	list := tfPromptList.([]interface{})
	var selectOptions []*variables.SelectOption

	for i := 0; i < len(list); i++ {
		item := list[i].(map[string]interface{})
		selectOptions = append(selectOptions, &variables.SelectOption{
			Value:       item["value"].(string),
			DisplayName: item["display_name"].(string),
		})
	}
	return selectOptions
}

func flattenPromptedVariableSettings(promptOptions *variables.VariablePromptOptions) []interface{} {
	if promptOptions == nil {
		return nil
	}
	flattenedPrompt := map[string]interface{}{}

	flattenedPrompt["description"] = promptOptions.Description
	flattenedPrompt["is_required"] = promptOptions.IsRequired
	flattenedPrompt["label"] = promptOptions.Label

	if promptOptions.DisplaySettings != nil {
		flattenedPrompt["display_settings"] = flattenPromptedVariableDisplaySettings(promptOptions.DisplaySettings)

	}

	return []interface{}{flattenedPrompt}
}

func flattenPromptedVariableDisplaySettings(displaySettings *variables.DisplaySettings) []interface{} {
	flattenedDisplaySettings := map[string]interface{}{}
	flattenedDisplaySettings["control_type"] = displaySettings.ControlType
	if displaySettings.ControlType == variables.ControlTypeSelect {
		flattenedDisplaySettings["select_option"] = flattenSelectOptions(displaySettings.SelectOptions)
	}
	return []interface{}{flattenedDisplaySettings}
}

func flattenSelectOptions(selectOptions []*variables.SelectOption) []map[string]interface{} {
	var options []map[string]interface{}
	for i := 0; i < len(selectOptions); i++ {
		options = append(options, map[string]interface{}{
			"value":        selectOptions[i].Value,
			"display_name": selectOptions[i].DisplayName,
		})
	}
	return options
}

func getVariablePromptOptionsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": getDescriptionSchema("variable prompt option"),
		"is_required": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"label": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"display_settings": {
			Elem:     &schema.Resource{Schema: getDisplaySettingsSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}

func getDisplaySettingsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"control_type": {
			Description: "The type of control for rendering this prompted variable. Valid types are `SingleLineText`, `MultiLineText`, `Checkbox`, `Select`.",
			Required:    true,
			Type:        schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"SingleLineText",
				"MultiLineText",
				"Checkbox",
				"Select",
			}, false)),
		},
		"select_option": {
			Elem:        &schema.Resource{Schema: getDisplaySettingsSelectOptionsSchema()},
			Description: "If the `control_type` is `Select`, then this value defines an option.",
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}

func getDisplaySettingsSelectOptionsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"value": {
			Description: "The select value",
			Required:    true,
			Type:        schema.TypeString,
		},
		"display_name": {
			Description: "The display name for the select value",
			Required:    true,
			Type:        schema.TypeString,
		},
	}
}
