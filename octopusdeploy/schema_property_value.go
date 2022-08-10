package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
)

func expandPropertyValue(value interface{}) core.PropertyValue {
	if v, ok := value.(string); ok {
		return core.NewPropertyValue(v, false)
	}

	if v, ok := value.([]interface{}); ok {
		if sensitiveValue, ok := v[0].(*core.SensitiveValue); ok {
			return core.PropertyValue{
				IsSensitive:    true,
				SensitiveValue: sensitiveValue,
			}
		}
	}

	panic("Invalid property value")
}

func flattenPropertyValue(propertyValue *core.PropertyValue) interface{} {
	if propertyValue == nil {
		return nil
	}

	if !propertyValue.IsSensitive {
		return propertyValue.Value
	}

	return nil
}
