package octopusdeploy

func flattenDisplaySettings(displaySettings map[string]interface{}) map[string]string {
	flattenedDisplaySettings := make(map[string]string, len(displaySettings))
	for key, displaySetting := range displaySettings {
		flattenedDisplaySettings[key] = displaySetting.(string)
	}
	return flattenedDisplaySettings
}