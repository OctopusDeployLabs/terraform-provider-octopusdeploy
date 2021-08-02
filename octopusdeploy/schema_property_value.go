package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandPropertyValue(value interface{}) *octopusdeploy.PropertyValue {
	if value == nil {
		return nil
	}

	flattenedValue := value.([]interface{})
	if len(flattenedValue) == 0 {
		return nil
	}

	flattenedPropertyValue := flattenedValue[0].(map[string]interface{})

	isSensitive := flattenedPropertyValue["is_sensitive"].(bool)

	if !isSensitive {
		value := flattenedPropertyValue["value"].(string)
		return &octopusdeploy.PropertyValue{
			IsSensitive: isSensitive,
			Value:       value,
		}
	}

	return &octopusdeploy.PropertyValue{
		IsSensitive:    isSensitive,
		SensitiveValue: expandSensitiveValue(flattenedPropertyValue["sensitive_value"]),
	}
}

func flattenPropertyValue(propertyValue *octopusdeploy.PropertyValue) []interface{} {
	if propertyValue == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"is_sensitive":    propertyValue.IsSensitive,
		"sensitive_value": flattenSensitiveValue(propertyValue.SensitiveValue),
		"value":           propertyValue.Value,
	}}
}

func getPropertyValueSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"is_sensitive": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"sensitive_value": {
			Optional: true,
			Elem:     &schema.Resource{Schema: getSensitiveValueSchema()},
			MaxItems: 1,
			Type:     schema.TypeList,
		},
		"value": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
