package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getDeployWindowsServiceActionSchema() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()
	addPrimaryPackageSchema(element, true)
	addDeployWindowsServiceSchema(element)
	// addCustomInstallationDirectoryFeature(element)
	// addCustomDeploymentScriptsFeature(element)
	// addJsonConfigurationVariablesFeature(element)
	// addConfigurationVariablesFeature(element)
	// addConfigurationTransformsFeature(element)
	// addSubstituteVariablesInFilesFeature(element)
	return actionSchema
}

func addWindowsServiceFeature(parent *schema.Resource) {
	element := &schema.Resource{
		Schema: map[string]*schema.Schema{},
	}
	addDeployWindowsServiceSchema(element)
	parent.Schema[constWindowsService] = &schema.Schema{
		Description: "Deploy a windows service feature",
		Type:        schema.TypeSet,
		Optional:    true,
		MaxItems:    1,
		Elem:        element,
	}
}

func addDeployWindowsServiceSchema(element *schema.Resource) {
	element.Schema[constServiceName] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The name of the service",
		Required:    true,
	}
	element.Schema[constDisplayName] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The display name of the service (optional)",
		Optional:    true,
	}
	element.Schema[constDescription] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "User-friendly description of the service (optional)",
		Optional:    true,
	}
	element.Schema[constExecutablePath] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The path to the executable relative to the package installation directory",
		Required:    true,
	}
	element.Schema[constArguments] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The command line arguments that will be passed to the service when it starts",
		Optional:    true,
	}
	element.Schema[constServiceAccount] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Which built-in account will the service run under. Can be LocalSystem, NT Authority\\NetworkService, NT Authority\\LocalService, _CUSTOM or an expression",
		Optional:    true,
		Default:     "LocalSystem",
	}
	element.Schema[constCustomAccountName] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The Windows/domain account of the custom user that the service will run under",
		Optional:    true,
	}
	element.Schema[constCustomAccountPassword] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The password for the custom account",
		Optional:    true,
	}
	element.Schema[constStartMode] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "When will the service start. Can be auto, delayed-auto, manual, unchanged or an expression",
		Optional:    true,
		Default:     "auto",
	}
	element.Schema[constDependencies] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Any dependencies that the service has. Separate the names using forward slashes (/).",
		Optional:    true,
	}
}

func buildDeployWindowsServiceActionResource(tfAction map[string]interface{}) model.DeploymentAction {
	resource := buildDeploymentActionResource(tfAction)
	resource.ActionType = "Octopus.WindowsService"
	addWindowsServiceToActionResource(tfAction, resource)
	return resource
}

func addWindowsServiceFeatureToActionResource(tfAction map[string]interface{}, action model.DeploymentAction) {
	if windowsServiceList, ok := tfAction[constWindowsService]; ok {
		tfWindowsService := windowsServiceList.(*schema.Set).List()
		if len(tfWindowsService) > 0 {
			addWindowsServiceToActionResource(tfWindowsService[0].(map[string]interface{}), action)
		}
	}
}

func addWindowsServiceToActionResource(tfAction map[string]interface{}, action model.DeploymentAction) {
	action.Properties["Octopus.Action.WindowsService.CreateOrUpdateService"] = "True"
	action.Properties["Octopus.Action.WindowsService.ServiceName"] = tfAction[constServiceName].(string)

	displayName := tfAction[constDisplayName]
	if displayName != nil {
		action.Properties["Octopus.Action.WindowsService.DisplayName"] = displayName.(string)
	}

	description := tfAction[constDescription]
	if description != nil {
		action.Properties["Octopus.Action.WindowsService.Description"] = description.(string)
	}

	action.Properties["Octopus.Action.WindowsService.ExecutablePath"] = tfAction[constExecutablePath].(string)

	args := tfAction[constArguments]
	if args != nil {
		action.Properties["Octopus.Action.WindowsService.Arguments"] = args.(string)
	}

	action.Properties["Octopus.Action.WindowsService.ServiceAccount"] = tfAction[constServiceAccount].(string)

	accountName := tfAction[constCustomAccountName]
	if accountName != nil {
		action.Properties["Octopus.Action.WindowsService.CustomAccountName"] = accountName.(string)
	}

	accountPassword := tfAction[constCustomAccountPassword]
	if accountPassword != nil {
		action.Properties["Octopus.Action.WindowsService.CustomAccountPassword"] = accountPassword.(string)
	}

	action.Properties["Octopus.Action.WindowsService.StartMode"] = tfAction[constStartMode].(string)

	dependencies := tfAction[constDependencies]
	if dependencies != nil {
		action.Properties["Octopus.Action.WindowsService.Dependencies"] = dependencies.(string)
	}
}
