package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandOfflineDrop(flattenedMap map[string]interface{}) *octopusdeploy.OfflineDropEndpoint {
	endpoint := octopusdeploy.NewOfflineDropEndpoint()
	endpoint.ApplicationsDirectory = flattenedMap["applications_directory"].(string)
	endpoint.Destination = expandOfflineDropDestination(flattenedMap["destination"])
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.SensitiveVariablesEncryptionPassword = octopusdeploy.NewSensitiveValue(flattenedMap["sensitive_variables_encryption_password"].(string))
	endpoint.WorkingDirectory = flattenedMap["working_directory"].(string)

	return endpoint
}

func flattenOfflineDrop(endpoint *octopusdeploy.OfflineDropEndpoint) []interface{} {
	if endpoint == nil {
		return nil
	}

	rawEndpoint := map[string]interface{}{
		"applications_directory": endpoint.ApplicationsDirectory,
		"destination":            flattenOfflineDropDestination(endpoint.Destination),
		"id":                     endpoint.GetID(),
		"working_directory":      endpoint.WorkingDirectory,
	}

	return []interface{}{rawEndpoint}
}

func getOfflineDropSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"applications_directory": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"destination": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getOfflineDropDestinationSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"working_directory": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"sensitive_variables_encryption_password": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
	}
}
