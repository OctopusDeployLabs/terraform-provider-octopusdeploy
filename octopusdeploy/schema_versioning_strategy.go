package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandVersioningStrategy(values interface{}) *octopusdeploy.VersioningStrategy {
	versioningStrategyList := values.(*schema.Set).List()
	versioningStrategyMap := versioningStrategyList[0].(map[string]interface{})

	versioningStrategy := &octopusdeploy.VersioningStrategy{}

	if versioningStrategyMap["donor_package_step_id"] != nil {
		donorPackageStepID := versioningStrategyMap["donor_package_step_id"].(string)
		if len(donorPackageStepID) > 0 {
			versioningStrategy.DonorPackageStepID = &donorPackageStepID
		}
	}

	versioningStrategy.DonorPackage = expandDeploymentActionPackage(versioningStrategyMap["donor_package"])
	versioningStrategy.Template = versioningStrategyMap["template"].(string)

	return versioningStrategy
}

func flattenVersioningStrategy(versioningStrategy *octopusdeploy.VersioningStrategy) []interface{} {
	if versioningStrategy == nil {
		return nil
	}

	flattenedVersioningStrategy := make(map[string]interface{})
	flattenedVersioningStrategy["donor_package"] = flattenDeploymentActionPackage(versioningStrategy.DonorPackage)
	flattenedVersioningStrategy["donor_package_step_id"] = versioningStrategy.DonorPackageStepID
	flattenedVersioningStrategy["template"] = versioningStrategy.Template
	return []interface{}{flattenedVersioningStrategy}
}

func getVersionStrategySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"donor_package": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getDeploymentActionPackageSchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"donor_package_step_id": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"template": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
