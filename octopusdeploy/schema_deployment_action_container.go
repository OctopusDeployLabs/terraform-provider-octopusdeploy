package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandDeploymentActionContainer(values interface{}) octopusdeploy.DeploymentActionContainer {
	flattenedValues := values.([]interface{})
	flattenedMap := flattenedValues[0].(map[string]interface{})

	return octopusdeploy.DeploymentActionContainer{
		FeedID: flattenedMap["feed_id"].(string),
		Image:  flattenedMap["image"].(string),
	}
}

func flattenDeploymentActionContainer(deploymentActionContainer octopusdeploy.DeploymentActionContainer) []interface{} {
	return []interface{}{map[string]interface{}{
		"feed_id": deploymentActionContainer.FeedID,
		"image":   deploymentActionContainer.Image,
	}}
}

func getDeploymentActionContainerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"feed_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"image": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
