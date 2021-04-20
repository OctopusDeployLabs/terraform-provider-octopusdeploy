package octopusdeploy

import "github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"

func expandExtensionSettingsValues(extensionSettingsValues []interface{}) []octopusdeploy.ExtensionSettingsValues {
	expandedExtensionSettingsValues := make([]octopusdeploy.ExtensionSettingsValues, len(extensionSettingsValues))
	for _, extensionSettingsValue := range extensionSettingsValues {
		extensionSettingsValueMap := extensionSettingsValue.(map[string]interface{})
		expandedExtensionSettingsValues = append(expandedExtensionSettingsValues, octopusdeploy.ExtensionSettingsValues{
			ExtensionID: extensionSettingsValueMap["extension_id"].(string),
			Values:      extensionSettingsValueMap["values"].([]interface{}),
		})
	}
	return expandedExtensionSettingsValues
}
