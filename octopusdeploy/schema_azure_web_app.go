package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func expandAzureWebApp(flattenedMap map[string]interface{}) *octopusdeploy.AzureWebAppEndpoint {
	endpoint := octopusdeploy.NewAzureWebAppEndpoint()
	endpoint.AccountID = flattenedMap["account_id"].(string)
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.ResourceGroupName = flattenedMap["resource_group_name"].(string)
	endpoint.WebAppName = flattenedMap["web_app_name"].(string)
	endpoint.WebAppSlotName = flattenedMap["web_app_slot_name"].(string)

	return endpoint
}

func flattenAzureWebApp(endpoint *octopusdeploy.AzureWebAppEndpoint) []interface{} {
	if endpoint == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"account_id":          endpoint.AccountID,
		"id":                  endpoint.GetID(),
		"resource_group_name": endpoint.ResourceGroupName,
		"web_app_name":        endpoint.WebAppName,
		"web_app_slot_name":   endpoint.WebAppSlotName,
	}}
}
