package octopusdeploy

import (
	"log"
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

func flattenPackageReference(packageReference octopusdeploy.PackageReference) map[string]interface{} {
	return map[string]interface{}{
		"acquisition_location": packageReference.AcquisitionLocation,
		"feed_id":              packageReference.FeedID,
		"id":                   packageReference.ID,
		"name":                 packageReference.Name,
		"package_id":           packageReference.PackageID,
		"properties":           packageReference.Properties,
	}
}

func flattenPackageReferences(packageReferences []octopusdeploy.PackageReference) []interface{} {
	if len(packageReferences) == 0 {
		return nil
	}

	flattenedPackageReferences := make([]interface{}, len(packageReferences))
	for key, packageReference := range packageReferences {
		flattenedPackageReferences[key] = flattenPackageReference(packageReference)
	}

	return flattenedPackageReferences
}

func getPackageSchema(required bool) *schema.Schema {
	return &schema.Schema{
		Computed:    !required,
		Description: "The primary package for the action",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"acquisition_location": {
					Default:     (string)(octopusdeploy.PackageAcquisitionLocationServer),
					Description: "Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression",
					Optional:    true,
					Type:        schema.TypeString,
				},
				"feed_id": {
					Default:     "feeds-builtin",
					Description: "The feed ID associated with this package reference.",
					Optional:    true,
					Type:        schema.TypeString,
				},
				"id":   getIDSchema(),
				"name": getNameSchema(false),
				"package_id": {
					Description: "The ID of the package",
					Required:    true,
					Type:        schema.TypeString,
				},
				"properties": {
					Computed: true,
					Optional: true,
					Type:     schema.TypeMap,
				},
			},
		},
		Optional: !required,
		Required: required,
		Type:     schema.TypeSet,
	}
}

func expandPackageReference(tfPkg map[string]interface{}) octopusdeploy.PackageReference {
	pkg := octopusdeploy.PackageReference{
		AcquisitionLocation: tfPkg["acquisition_location"].(string),
		FeedID:              tfPkg["feed_id"].(string),
		Name:                getStringOrEmpty(tfPkg["name"]),
		PackageID:           tfPkg["package_id"].(string),
		Properties:          map[string]string{},
	}

	if properties := tfPkg["properties"]; properties != nil {
		propertyMap := properties.(map[string]interface{})
		if len(propertyMap) > 0 {
			// TODO: convert propertyMap to map[string]string
			// pkg.Properties := convertedPropertyMap
			log.Println(propertyMap)
		}
	}

	extract := tfPkg["extract_during_deployment"]
	if extract != nil {
		pkg.Properties["Extract"] = strconv.FormatBool(extract.(bool))
	}

	return pkg
}
