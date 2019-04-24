package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func addPrimaryPackageSchema(element *schema.Resource, required bool) {
	element.Schema["primary_package"] = getPackageSchema(required)
	element.Schema["primary_package"].MaxItems = 1
}

func addPackagesSchema(element *schema.Resource, primaryIsRequired bool) {
	addPrimaryPackageSchema(element, primaryIsRequired)

	element.Schema["package"] = getPackageSchema(false)

	packageElementSchema := element.Schema["package"].Elem.(*schema.Resource).Schema

	packageElementSchema["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The name of the package",
		Required:    true,
	}

	packageElementSchema["extract_during_deployment"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Whether to extract the package during deployment",
		Optional:    true,
		Default:     "true",
	}
}

func getPackageSchema(required bool) *schema.Schema {
	return &schema.Schema{
		Description: "The primary package for the action",
		Type:        schema.TypeSet,
		Required:    required,
		Optional:    !required,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"package_id": {
					Type:        schema.TypeString,
					Description: "The ID of the package",
					Required:    true,
				},
				"feed_id": {
					Type:        schema.TypeString,
					Description: "The feed to retrieve the package from",
					Default:     "feeds-builtin",
					Optional:    true,
				},
				"acquisition_location": {
					Type:        schema.TypeString,
					Description: "Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression",
					Default:     (string)(octopusdeploy.PackageAcquisitionLocation_Server),
					Optional:    true,
				},
				"property": getPropertySchema(),
			},
		},
	}
}

func buildPackageReferenceResource(tfPkg map[string]interface{}) octopusdeploy.PackageReference {
	pkg := octopusdeploy.PackageReference{
		Name:                getStringOrEmpty(tfPkg["name"]),
		PackageId:           tfPkg["package_id"].(string),
		FeedId:              tfPkg["feed_id"].(string),
		AcquisitionLocation: tfPkg["acquisition_location"].(string),
		Properties:          buildPropertiesMap(tfPkg["property"]),
	}

	extract := tfPkg["extract_during_deployment"]
	if extract != nil {
		pkg.Properties["Extract"] = extract.(string)
	}

	return pkg
}
