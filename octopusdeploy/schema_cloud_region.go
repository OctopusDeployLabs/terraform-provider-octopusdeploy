package octopusdeploy

import "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"

func expandCloudRegion(flattenedMap map[string]interface{}) *machines.CloudRegionEndpoint {
	endpoint := machines.NewCloudRegionEndpoint()
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.DefaultWorkerPoolID = flattenedMap["default_worker_pool_id"].(string)

	return endpoint
}
