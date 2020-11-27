package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandAzureWebApp(flattenedMap map[string]interface{}) *octopusdeploy.AzureWebAppEndpoint {
	endpoint := octopusdeploy.NewAzureWebAppEndpoint()
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.ResourceGroupName = flattenedMap["resource_group_name"].(string)
	endpoint.WebAppName = flattenedMap["web_app_name"].(string)
	endpoint.WebAppSlotName = flattenedMap["web_app_slot_name"].(int)

	return endpoint
}

func flattenAzureWebApp(endpoint *octopusdeploy.AzureWebAppEndpoint) []interface{} {
	if endpoint == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"id":                  endpoint.GetID(),
		"resource_group_name": endpoint.ResourceGroupName,
		"web_app_name":        endpoint.WebAppName,
		"web_app_slot_name":   endpoint.WebAppSlotName,
	}}
}

func getAzureWebAppSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"resource_group_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"web_app_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"web_app_slot_name": {
			Optional: true,
			Type:     schema.TypeInt,
		},
	}
}
