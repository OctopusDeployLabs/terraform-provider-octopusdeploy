package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func expandOfflinePackageDrop(flattenedMap map[string]interface{}) *octopusdeploy.OfflinePackageDropEndpoint {
	endpoint := octopusdeploy.NewOfflinePackageDropEndpoint()
	endpoint.ApplicationsDirectory = flattenedMap["applications_directory"].(string)
	endpoint.Destination = expandOfflinePackageDropDestination(flattenedMap["destination"])
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.SensitiveVariablesEncryptionPassword = octopusdeploy.NewSensitiveValue(flattenedMap["sensitive_variables_encryption_password"].(string))
	endpoint.WorkingDirectory = flattenedMap["working_directory"].(string)

	return endpoint
}

func flattenOfflinePackageDrop(endpoint *octopusdeploy.OfflinePackageDropEndpoint) []interface{} {
	rawEndpoint := map[string]interface{}{
		"applications_directory": endpoint.ApplicationsDirectory,
		"destination":            flattenOfflinePackageDropDestination(&endpoint.Destination),
		"id":                     endpoint.GetID(),
		"working_directory":      endpoint.WorkingDirectory,
	}

	return []interface{}{rawEndpoint}
}
