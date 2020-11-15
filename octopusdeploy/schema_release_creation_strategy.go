package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandReleaseCreationStrategy(releaseCreationStrategy []interface{}) *octopusdeploy.ReleaseCreationStrategy {
	releaseCreationStrategyMap := releaseCreationStrategy[0].(map[string]interface{})
	return &octopusdeploy.ReleaseCreationStrategy{
		ChannelID:                    releaseCreationStrategyMap["channel_id"].(string),
		ReleaseCreationPackage:       expandDeploymentActionPackage(releaseCreationStrategyMap["release_creation_package"].([]interface{})),
		ReleaseCreationPackageStepID: releaseCreationStrategyMap["release_creation_package_step_id"].(*string),
	}
}

func flattenReleaseCreationStrategy(releaseCreationStrategy *octopusdeploy.ReleaseCreationStrategy) []interface{} {
	if releaseCreationStrategy == nil {
		return nil
	}

	flattenedReleaseCreationStrategy := make(map[string]interface{})
	flattenedReleaseCreationStrategy["channel_id"] = releaseCreationStrategy.ChannelID
	flattenedReleaseCreationStrategy["release_creation_package"] = flattenDeploymentActionPackage(releaseCreationStrategy.ReleaseCreationPackage)
	flattenedReleaseCreationStrategy["release_creation_package_step_id"] = releaseCreationStrategy.ReleaseCreationPackageStepID
	return []interface{}{flattenedReleaseCreationStrategy}
}

func getReleaseCreationStrategySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"channel_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"release_creation_package": {
			Computed: true,
			Optional: true,
			Elem:     &schema.Resource{Schema: getDeploymentActionPackageSchema()},
			MaxItems: 1,
			Type:     schema.TypeList,
		},
		"release_creation_package_step_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
