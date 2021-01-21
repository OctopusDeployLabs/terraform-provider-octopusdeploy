package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenVersioningStrategy(versioningStrategy octopusdeploy.VersioningStrategy) []interface{} {
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
			Type:     schema.TypeList,
		},
		"donor_package_step_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"template": {
			Computed: true,
			Type:     schema.TypeString,
		},
	}
}
