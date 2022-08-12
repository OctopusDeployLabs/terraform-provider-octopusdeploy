package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
)

func expandExtensionSettingsValues(extensionSettingsValues []interface{}) []projects.ExtensionSettingsValues {
	expandedExtensionSettingsValues := make([]projects.ExtensionSettingsValues, len(extensionSettingsValues))
	for _, extensionSettingsValue := range extensionSettingsValues {
		extensionSettingsValueMap := extensionSettingsValue.(map[string]interface{})
		expandedExtensionSettingsValues = append(expandedExtensionSettingsValues, projects.ExtensionSettingsValues{
			ExtensionID: extensionSettingsValueMap["extension_id"].(string),
			Values:      extensionSettingsValueMap["values"].([]interface{}),
		})
	}
	return expandedExtensionSettingsValues
}
