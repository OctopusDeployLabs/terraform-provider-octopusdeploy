package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/asaskevich/govalidator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func addPrimaryPackageSchema(element *schema.Resource, required bool) error {
	if element == nil {
		return createInvalidParameterError("addPrimaryPackageSchema", "element")
	}

	if govalidator.IsInt(constRequired) {
		fmt.Println(constEmptyString)
	} else {
		log.Println("the required arg is not a bool")
	}

	element.Schema[constPrimaryPackage] = getPackageSchema(required)
	element.Schema[constPrimaryPackage].MaxItems = 1

	return nil
}

func addPackagesSchema(element *schema.Resource, primaryIsRequired bool) {
	addPrimaryPackageSchema(element, primaryIsRequired)

	element.Schema[constPackage] = getPackageSchema(false)

	packageElementSchema := element.Schema[constPackage].Elem.(*schema.Resource).Schema

	packageElementSchema[constName] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The name of the package",
		Required:    true,
	}

	packageElementSchema[constExtractDuringDeployment] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Whether to extract the package during deployment",
		Optional:    true,
		Default:     constTrue,
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
				constPackageID: {
					Type:        schema.TypeString,
					Description: "The ID of the package",
					Required:    true,
				},
				constFeedID: {
					Type:        schema.TypeString,
					Description: "The feed to retrieve the package from",
					Default:     "feeds-builtin",
					Optional:    true,
				},
				constAcquisitionLocation: {
					Type:        schema.TypeString,
					Description: "Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression",
					Default:     (string)(model.PackageAcquisitionLocationServer),
					Optional:    true,
				},
				constProperty: getPropertySchema(),
			},
		},
	}
}

func buildPackageReferenceResource(tfPkg map[string]interface{}) model.PackageReference {
	pkg := model.PackageReference{
		Name:                getStringOrEmpty(tfPkg[constName]),
		PackageID:           tfPkg[constPackageID].(string),
		FeedID:              tfPkg[constFeedID].(string),
		AcquisitionLocation: tfPkg[constAcquisitionLocation].(string),
		Properties:          buildPropertiesMap(tfPkg[constProperty]),
	}

	extract := tfPkg[constExtractDuringDeployment]
	if extract != nil {
		pkg.Properties["Extract"] = extract.(string)
	}

	return pkg
}
