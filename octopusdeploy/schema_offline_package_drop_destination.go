package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandOfflinePackageDropDestination(values interface{}) *machines.OfflinePackageDropDestination {
	if values == nil {
		return nil
	}

	flattenedValues := values.([]interface{})
	flattenedEndpoint := flattenedValues[0].(map[string]interface{})

	return &machines.OfflinePackageDropDestination{
		DestinationType: flattenedEndpoint["destination_type"].(string),
		DropFolderPath:  flattenedEndpoint["drop_folder_path"].(string),
	}
}

func flattenOfflinePackageDropDestination(offlineDropDestination *machines.OfflinePackageDropDestination) []interface{} {
	if offlineDropDestination == nil {
		return nil
	}

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
