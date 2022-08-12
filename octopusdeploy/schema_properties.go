package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
)

func expandProperties(propertyValues interface{}) map[string]core.PropertyValue {
	if propertyValues == nil {
		return nil
	}

	expandedPropertyValues := map[string]core.PropertyValue{}
	for k, v := range propertyValues.(map[string]interface{}) {
		expandedPropertyValues[k] = expandPropertyValue(v)
	}
	return expandedPropertyValues
}

func flattenProperties(propertyValues map[string]core.PropertyValue) map[string]interface{} {
	if len(propertyValues) == 0 {
		return nil
	}

	flattenedProperties := map[string]interface{}{}
	for i := range propertyValues {
		propertyValue := propertyValues[i]
		flattenedProperties[i] = flattenPropertyValue(&propertyValue)
	}
	return flattenedProperties
}
