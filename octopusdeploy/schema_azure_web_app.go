package octopusdeploy

import "github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"

func expandAzureWebApp(flattenedMap map[string]interface{}) *machines.AzureWebAppEndpoint {
	endpoint := machines.NewAzureWebAppEndpoint()
	endpoint.AccountID = flattenedMap["account_id"].(string)
	endpoint.ID = flattenedMap["id"].(string)
	endpoint.ResourceGroupName = flattenedMap["resource_group_name"].(string)
	endpoint.WebAppName = flattenedMap["web_app_name"].(string)
	endpoint.WebAppSlotName = flattenedMap["web_app_slot_name"].(string)

	return endpoint
}
