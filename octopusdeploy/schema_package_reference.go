package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func addPrimaryPackageSchema(element *schema.Resource, required bool) error {
	if element == nil {
		return createInvalidParameterError("addPrimaryPackageSchema", "element")
	}

	element.Schema["primary_package"] = getPackageSchema(required)
	element.Schema["primary_package"].MaxItems = 1

	return nil
}

func addPackagesSchema(element *schema.Resource, primaryIsRequired bool) {
	addPrimaryPackageSchema(element, primaryIsRequired)

	element.Schema["package"] = getPackageSchema(false)

	packageElementSchema := element.Schema["package"].Elem.(*schema.Resource).Schema

	packageElementSchema["name"] = &schema.Schema{
		Description: "The name of the package",
		Required:    true,
		Type:        schema.TypeString,
	}

	packageElementSchema["extract_during_deployment"] = &schema.Schema{
		Default:     true,
		Description: "Whether to extract the package during deployment",
		Optional:    true,
		Type:        schema.TypeBool,
	}
}

func getPackageSchema(required bool) *schema.Schema {
	return &schema.Schema{
		Description: "The primary package for the action",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"package_id": {
					Description: "The ID of the package",
					Required:    true,
					Type:        schema.TypeString,
				},
				"feed_id": {
					Default:     "feeds-builtin",
					Description: "The feed to retrieve the package from",
					Optional:    true,
					Type:        schema.TypeString,
				},
				"acquisition_location": {
					Default:     (string)(octopusdeploy.PackageAcquisitionLocationServer),
					Description: "Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression",
					Optional:    true,
					Type:        schema.TypeString,
				},
				"property": getPropertySchema(),
			},
		},
		Optional: !required,
		Required: required,
		Type:     schema.TypeSet,
	}
}

func buildPackageReferenceResource(tfPkg map[string]interface{}) octopusdeploy.PackageReference {
	pkg := octopusdeploy.PackageReference{
		AcquisitionLocation: tfPkg["acquisition_location"].(string),
		FeedID:              tfPkg["feed_id"].(string),
		Name:                getStringOrEmpty(tfPkg["name"]),
		PackageID:           tfPkg["package_id"].(string),
		Properties:          buildPropertiesMap(tfPkg["property"]),
	}

	extract := tfPkg["extract_during_deployment"]
	if extract != nil {
		pkg.Properties["Extract"] = strconv.FormatBool(extract.(bool))
	}

	return pkg
}
