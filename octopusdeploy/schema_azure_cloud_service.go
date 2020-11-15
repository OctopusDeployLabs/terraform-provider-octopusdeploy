package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandAzureCloudService(d *schema.ResourceData) *octopusdeploy.AzureCloudServiceEndpoint {
	endpoint := octopusdeploy.NewAzureCloudServiceEndpoint()
	endpoint.ID = d.Id()

	if v, ok := d.GetOk("account_id"); ok {
		endpoint.AccountID = v.(string)
	}

	if v, ok := d.GetOk("cloud_service_name"); ok {
		endpoint.CloudServiceName = v.(string)
	}

	if v, ok := d.GetOk("default_worker_pool_id"); ok {
		endpoint.DefaultWorkerPoolID = v.(string)
	}

	if v, ok := d.GetOk("use_current_instance_count"); ok {
		endpoint.UseCurrentInstanceCount = v.(bool)
	}

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

func getAzureCloudServiceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"cloud_service_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"default_worker_pool_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"use_current_instance_count": {
			Optional: true,
			Type:     schema.TypeBool,
		},
	}
}
