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

func flattenAzureCloudService(endpoint *octopusdeploy.AzureCloudServiceEndpoint) []interface{} {
	if endpoint == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"account_id":                 endpoint.AccountID,
		"cloud_service_name":         endpoint.CloudServiceName,
		"default_worker_pool_id":     endpoint.DefaultWorkerPoolID,
		"id":                         endpoint.GetID(),
		"use_current_instance_count": endpoint.UseCurrentInstanceCount,
	}}
}
