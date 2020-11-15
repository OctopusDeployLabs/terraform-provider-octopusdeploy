package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandCloudRegion(d *schema.ResourceData) *octopusdeploy.CloudRegionEndpoint {
	endpoint := octopusdeploy.NewCloudRegionEndpoint()
	endpoint.ID = d.Id()

	if v, ok := d.GetOk("default_worker_pool_id"); ok {
		endpoint.DefaultWorkerPoolID = v.(string)
	}

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

func getCloudRegionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"default_worker_pool_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
	}
}
