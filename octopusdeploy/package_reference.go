package octopusdeploy

import (
	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)


func getPrimaryPackageSchema() *schema.Schema {
	return &schema.Schema{
		Description: "The primary package for the action",
		Type:        schema.TypeSet,
		Optional:    true,
		MaxItems:	 1,
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
					Default: 	"feeds-builtin",
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

func getPackageSchema() *schema.Schema {
	s := getPrimaryPackageSchema();
	elementSchema := s.Elem.(*schema.Resource).Schema
	elementSchema["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The name of the package",
		Required:    true,
	}
	elementSchema["extract_during_deployment"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Whether to extract the package during deployment",
		Optional:    true,
	}
	return s
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

	return pkg;
}

