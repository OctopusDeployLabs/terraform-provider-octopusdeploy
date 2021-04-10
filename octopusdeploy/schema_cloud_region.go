package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func expandCloudRegion(flattenedMap map[string]interface{}) *octopusdeploy.CloudRegionEndpoint {
	endpoint := octopusdeploy.NewCloudRegionEndpoint()
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.DefaultWorkerPoolID = flattenedMap["default_worker_pool_id"].(string)

	return endpoint
}

func flattenCloudRegion(endpoint *octopusdeploy.CloudRegionEndpoint) []interface{} {
	if endpoint == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"default_worker_pool_id": endpoint.DefaultWorkerPoolID,
		"id":                     endpoint.GetID(),
	}}
}
