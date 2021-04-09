package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func expandProperties(properties interface{}) map[string]octopusdeploy.PropertyValue {
	if properties == nil {
		return nil
	}

	expandedProperties := map[string]octopusdeploy.PropertyValue{}
	for k, v := range properties.(map[string]interface{}) {
		expandedProperties[k] = octopusdeploy.NewPropertyValue(v.(string), false)
	}
	return expandedProperties
}

func flattenProperties(properties map[string]octopusdeploy.PropertyValue) map[string]interface{} {
	if len(properties) == 0 {
		return nil
	}

	flattenedProperties := map[string]interface{}{}
	for k, v := range properties {
		flattenedProperties[k] = v.Value
	}
	return flattenedProperties
}
