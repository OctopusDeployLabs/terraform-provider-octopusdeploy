package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func getDeployWindowsServiceActionSchema() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()
	addPrimaryPackageSchema(element, true)
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
	element := &schema.Resource{
		Schema: map[string]*schema.Schema{},
	}
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
		Description: "Which built-in account will the service run under. Can be LocalSystem, NT Authority\\NetworkService, NT Authority\\LocalService, _CUSTOM or an expression",
		Optional:    true,
		Default:     "LocalSystem",
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
		Description: "When will the service start. Can be auto, delayed-auto, manual, unchanged or an expression",
		Optional:    true,
		Default:     "auto",
	}
	element.Schema["dependencies"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Any dependencies that the service has. Separate the names using forward slashes (/).",
		Optional:    true,
	}
}

func buildDeployWindowsServiceActionResource(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	resource := buildDeploymentActionResource(tfAction)
	resource.ActionType = "Octopus.WindowsService"
	addWindowsServiceToActionResource(tfAction, resource)
	return resource
}

func addWindowsServiceFeatureToActionResource(tfAction map[string]interface{}, action octopusdeploy.DeploymentAction) {
	if windowsServiceList, ok := tfAction["windows_service"]; ok {
		tfWindowsService := windowsServiceList.(*schema.Set).List()
		if len(tfWindowsService) > 0 {
			addWindowsServiceToActionResource(tfWindowsService[0].(map[string]interface{}), action)
		}
	}
}

func addWindowsServiceToActionResource(tfAction map[string]interface{}, action octopusdeploy.DeploymentAction) {
	action.Properties["Octopus.Action.WindowsService.CreateOrUpdateService"] = "True"
	action.Properties["Octopus.Action.WindowsService.ServiceName"] = tfAction["service_name"].(string)

	displayName := tfAction["display_name"]
	if displayName != nil {
		action.Properties["Octopus.Action.WindowsService.DisplayName"] = displayName.(string)
	}

	description := tfAction["description"]
	if description != nil {
		action.Properties["Octopus.Action.WindowsService.Description"] = description.(string)
	}

	action.Properties["Octopus.Action.WindowsService.ExecutablePath"] = tfAction["executable_path"].(string)

	args := tfAction["arguments"]
	if args != nil {
		action.Properties["Octopus.Action.WindowsService.Arguments"] = args.(string)
	}

	action.Properties["Octopus.Action.WindowsService.ServiceAccount"] = tfAction["service_account"].(string)

	accountName := tfAction["custom_account_name"]
	if accountName != nil {
		action.Properties["Octopus.Action.WindowsService.CustomAccountName"] = accountName.(string)
	}

	accountPassword := tfAction["custom_account_password"]
	if accountPassword != nil {
		action.Properties["Octopus.Action.WindowsService.CustomAccountPassword"] = accountPassword.(string)
	}

	action.Properties["Octopus.Action.WindowsService.StartMode"] = tfAction["start_mode"].(string)

	dependencies := tfAction["dependencies"]
	if dependencies != nil {
		action.Properties["Octopus.Action.WindowsService.Dependencies"] = dependencies.(string)
	}
}
