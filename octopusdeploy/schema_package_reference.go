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

	// TODO: update name in schema to always be empty string

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
		Computed:    true,
		Description: "Whether to extract the package during deployment",
		Optional:    true,
		Type:        schema.TypeBool,
	}
}

func flattenPackageReference(packageReference octopusdeploy.PackageReference) map[string]interface{} {
	flattenedPackageReference := map[string]interface{}{
		"acquisition_location": packageReference.AcquisitionLocation,
		"feed_id":              packageReference.FeedID,
		"id":                   packageReference.ID,
		"name":                 packageReference.Name,
		"package_id":           packageReference.PackageID,
		"properties":           packageReference.Properties,
	}

	if v, ok := packageReference.Properties["Extract"]; ok {
		flattenedPackageReference["extract_during_deployment"] = v
	}

	return flattenedPackageReference
}

func getPackageSchema(required bool) *schema.Schema {
	return &schema.Schema{
		Computed:    !required,
		Description: "The package assocated with this action.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"acquisition_location": {
					Default:     "Server",
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
					Description: "The ID of the package.",
					Required:    true,
					Type:        schema.TypeString,
				},
				"properties": {
					Computed:    true,
					Description: "A list of properties associated with this package.",
					Elem:        &schema.Schema{Type: schema.TypeString},
					Optional:    true,
					Type:        schema.TypeMap,
				},
			},
		},
		Optional: !required,
		Required: required,
		Type:     schema.TypeList,
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

	if v, ok := tfPkg["extract_during_deployment"]; ok {
		pkg.Properties["Extract"] = strconv.FormatBool(v.(bool))
	}

	if properties := tfPkg["properties"]; properties != nil {
		propertyMap := properties.(map[string]interface{})
		for k, v := range propertyMap {
			pkg.Properties[k] = v.(string)
		}
	}

	return pkg
}
