package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandOfflineDropDestination(values interface{}) *octopusdeploy.OfflineDropDestination {
	if values == nil {
		return nil
	}

	flattenedValues := values.([]interface{})
	flattenedEndpoint := flattenedValues[0].(map[string]interface{})

	return &octopusdeploy.OfflineDropDestination{
		DestinationType: flattenedEndpoint["destination_type"].(string),
		DropFolderPath:  flattenedEndpoint["drop_folder_path"].(string),
	}
}

func flattenOfflineDropDestination(offlineDropDestination *octopusdeploy.OfflineDropDestination) []interface{} {
	if offlineDropDestination == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"destination_type": offlineDropDestination.DestinationType,
		"drop_folder_path": offlineDropDestination.DropFolderPath,
	}}
}

func getOfflineDropDestinationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"destination_type": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"drop_folder_path": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
