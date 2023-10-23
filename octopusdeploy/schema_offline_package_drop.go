package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
)

func expandOfflinePackageDrop(flattenedMap map[string]interface{}) *machines.OfflinePackageDropEndpoint {
	endpoint := machines.NewOfflinePackageDropEndpoint()
	endpoint.ApplicationsDirectory = flattenedMap["applications_directory"].(string)
	endpoint.Destination = expandOfflinePackageDropDestination(flattenedMap["destination"])
	endpoint.ID = flattenedMap["id"].(string)

	if encryptedPassword, ok := flattenedMap["sensitive_variables_encryption_password"]; ok {
		endpoint.SensitiveVariablesEncryptionPassword = core.NewSensitiveValue(encryptedPassword.(string))
	}

	endpoint.WorkingDirectory = flattenedMap["working_directory"].(string)

	return endpoint
}
