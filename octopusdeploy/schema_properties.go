package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func expandProperties(propertyValues interface{}) map[string]octopusdeploy.PropertyValue {
	if propertyValues == nil {
		return nil
	}

	expandedPropertyValues := map[string]octopusdeploy.PropertyValue{}
	for k, v := range propertyValues.(map[string]interface{}) {
		expandedPropertyValues[k] = *expandPropertyValue(v)
	}
	return expandedPropertyValues
}

func flattenProperties(propertyValues map[string]octopusdeploy.PropertyValue) map[string]interface{} {
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
