package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandOfflinePackageDropDestination(values interface{}) octopusdeploy.OfflinePackageDropDestination {
	if values == nil {
		return octopusdeploy.OfflinePackageDropDestination{}
	}

	flattenedValues := values.([]interface{})
	flattenedEndpoint := flattenedValues[0].(map[string]interface{})

	return octopusdeploy.OfflinePackageDropDestination{
		DestinationType: flattenedEndpoint["destination_type"].(string),
		DropFolderPath:  flattenedEndpoint["drop_folder_path"].(string),
	}
}

func flattenOfflinePackageDropDestination(offlineDropDestination octopusdeploy.OfflinePackageDropDestination) []interface{} {
	return []interface{}{map[string]interface{}{
		"destination_type": offlineDropDestination.DestinationType,
		"drop_folder_path": offlineDropDestination.DropFolderPath,
	}}
}

func getOfflinePackageDropDestinationSchema() map[string]*schema.Schema {
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
