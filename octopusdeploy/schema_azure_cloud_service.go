package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func expandAzureCloudService(flattenedMap map[string]interface{}) *octopusdeploy.AzureCloudServiceEndpoint {
	endpoint := octopusdeploy.NewAzureCloudServiceEndpoint()
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.AccountID = flattenedMap["account_id"].(string)
	endpoint.CloudServiceName = flattenedMap["cloud_service_name"].(string)
	endpoint.DefaultWorkerPoolID = flattenedMap["default_worker_pool_id"].(string)
	endpoint.UseCurrentInstanceCount = flattenedMap["use_current_instance_count"].(bool)

	return endpoint
}
