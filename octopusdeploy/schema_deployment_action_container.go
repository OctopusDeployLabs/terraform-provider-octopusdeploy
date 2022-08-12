package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandContainer(values interface{}) *deployments.DeploymentActionContainer {
	if values == nil {
		return nil
	}

	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 || flattenedValues[0] == nil {
		return nil
	}

	flattenedMap := flattenedValues[0].(map[string]interface{})
	if len(flattenedMap) == 0 {
		return nil
	}

	deploymentActionContainer := &deployments.DeploymentActionContainer{}

	if feedID := flattenedMap["feed_id"]; feedID != nil {
		deploymentActionContainer.FeedID = feedID.(string)
	}

	if image := flattenedMap["image"]; image != nil {
		deploymentActionContainer.Image = image.(string)
	}

	return deploymentActionContainer
}

func flattenContainer(deploymentActionContainer *deployments.DeploymentActionContainer) []interface{} {
	if deploymentActionContainer == nil {
		return nil
	}

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
