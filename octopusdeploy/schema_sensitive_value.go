package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandSensitiveValue(values interface{}) *octopusdeploy.SensitiveValue {
	if values == nil {
		return nil
	}

	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 {
		return nil
	}

	flattenedSensitiveValue := flattenedValues[0].(map[string]interface{})

	newValue := flattenedSensitiveValue["new_value"].(string)
	expandedSensitiveValue := octopusdeploy.NewSensitiveValue(newValue)

	if hint, ok := flattenedSensitiveValue["hint"]; ok {
		hintString := hint.(string)
		expandedSensitiveValue.Hint = &hintString
	}

	return expandedSensitiveValue
}

func flattenSensitiveValue(sensitiveValue *octopusdeploy.SensitiveValue) []interface{} {
	if sensitiveValue == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"hint":      sensitiveValue.Hint,
		"new_value": sensitiveValue.NewValue,
	}}
}

func getSensitiveValueSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"hint": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"new_value": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
	}
}
