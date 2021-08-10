package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func expandPropertyValue(value interface{}) octopusdeploy.PropertyValue {
	if v, ok := value.(string); ok {
		return octopusdeploy.NewPropertyValue(v, false)
	}

	if v, ok := value.([]interface{}); ok {
		if sensitiveValue, ok := v[0].(*octopusdeploy.SensitiveValue); ok {
			return octopusdeploy.PropertyValue{
				IsSensitive:    true,
				SensitiveValue: sensitiveValue,
			}
		}
	}

	panic("Invalid property value")
}

func flattenPropertyValue(propertyValue *octopusdeploy.PropertyValue) interface{} {
	if propertyValue == nil {
		return nil
	}

	if !propertyValue.IsSensitive {
		return propertyValue.Value
	}

	return nil
}
