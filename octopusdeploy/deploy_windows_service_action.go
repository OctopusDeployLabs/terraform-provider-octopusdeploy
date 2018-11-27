package octopusdeploy

import (
	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func getDeployWindowsServiceActionSchema() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()
	addPrimaryPackageSchema(element)
	addDeployWindowsServiceSchema(element)
	//addCustomInstallationDirectoryFeature(element)
	//addCustomDeploymentScriptsFeature(element)
	//addJsonConfigurationVariablesFeature(element)
	//addConfigurationVariablesFeature(element)
	//addConfigurationTransformsFeature(element)
	//addSubstituteVariablesInFilesFeature(element)
	return actionSchema
}

func addWindowsServiceFeature(parent *schema.Resource) {
	element := &schema.Resource{}
	addDeployWindowsServiceSchema(element)
	parent.Schema["windows_service"] = &schema.Schema{
		Description: "Deploy a windows service feature",
		Type:        schema.TypeSet,
		Optional:    true,
		MaxItems:    1,
		Elem:        element,
	}
}

func addDeployWindowsServiceSchema(element *schema.Resource) {
	element.Schema["service_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The name of the service",
		Required:    true,
	}
	element.Schema["display_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The display name of the service (optional)",
		Optional:    true,
	}
	element.Schema["description"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "User-friendly description of the service (optional)",
		Optional:    true,
	}
	element.Schema["executable_path"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The path to the executable relative to the package installation directory",
		Required:    true,
	}
	element.Schema["arguments"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The command line arguments that will be passed to the service when it starts",
		Optional:    true,
	}
	element.Schema["service_account"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Which built-in account will the service run under",
		Required:    true,
	}
	element.Schema["custom_account_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The Windows/domain account of the custom user that the service will run under",
		Optional:    true,
	}
	element.Schema["custom_account_password"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The password for the custom account",
		Optional:    true,
	}
	element.Schema["start_mode"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "When will the service start. Can be Automatic, Automatic (delayed), Manual, Unchanged or an expression",
		Required:    true,
	}
	element.Schema["dependencies"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Any dependencies that the service has. Separate the names using forward slashes (/).",
		Optional:    true,
	}
}


func buildDeployWindowsServiceActionResource(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	resource := buildDeploymentActionResource(tfAction)
	addWindowsServiceToResource(tfAction, resource)
	return resource
}


func addWindowsServiceFeatureToResource(tfAction map[string]interface{}, resource octopusdeploy.DeploymentAction) {
	if windowsServiceList, ok := tfAction["windows_service"]; ok {
		tfWindowsService := windowsServiceList.(*schema.Set).List()
		if(len(tfWindowsService) > 0) {
			addWindowsServiceToResource(tfWindowsService[0].(map[string]interface{}), resource)
		}
	}
}

func addWindowsServiceToResource(tfAction map[string]interface{}, resource octopusdeploy.DeploymentAction) {
	// TODO
}
