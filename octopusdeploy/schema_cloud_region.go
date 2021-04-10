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
