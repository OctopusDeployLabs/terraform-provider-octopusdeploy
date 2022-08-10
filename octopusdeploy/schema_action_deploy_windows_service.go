package octopusdeploy

import (
	"strconv"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func getDeployWindowsServiceActionSchema() *schema.Schema {
	actionSchema, element := getActionSchema()
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
	parent.Schema["windows_service"] = &schema.Schema{
		Description: "Deploy a windows service feature",
		Elem:        element,
		MaxItems:    1,
		Optional:    true,
		Type:        schema.TypeSet,
	}
}

func addDeployWindowsServiceSchema(element *schema.Resource) {
	element.Schema["arguments"] = &schema.Schema{
		Description: "The command line arguments that will be passed to the service when it starts",
		Optional:    true,
		Type:        schema.TypeString,
	}
	element.Schema["create_or_update_service"] = &schema.Schema{
		Computed: true,
		Optional: true,
		Type:     schema.TypeBool,
	}
	element.Schema["custom_account_name"] = &schema.Schema{
		Description: "The Windows/domain account of the custom user that the service will run under",
		Optional:    true,
		Type:        schema.TypeString,
	}
	element.Schema["custom_account_password"] = &schema.Schema{
		Computed:    true,
		Description: "The password for the custom account",
		Optional:    true,
		Sensitive:   true,
		Type:        schema.TypeString,
	}
	element.Schema["dependencies"] = &schema.Schema{
		Description: "Any dependencies that the service has. Separate the names using forward slashes (/).",
		Optional:    true,
		Type:        schema.TypeString,
	}
	element.Schema["description"] = &schema.Schema{
		Description: "User-friendly description of the service (optional)",
		Optional:    true,
		Type:        schema.TypeString,
	}
	element.Schema["display_name"] = &schema.Schema{
		Description: "The display name of the service (optional)",
		Optional:    true,
		Type:        schema.TypeString,
	}
	element.Schema["executable_path"] = &schema.Schema{
		Description: "The path to the executable relative to the package installation directory",
		Required:    true,
		Type:        schema.TypeString,
	}
	element.Schema["service_account"] = &schema.Schema{
		Description: "Which built-in account will the service run under. Can be LocalSystem, NT Authority\\NetworkService, NT Authority\\LocalService, _CUSTOM or an expression",
		Default:     "LocalSystem",
		Optional:    true,
		Type:        schema.TypeString,
	}
	element.Schema["service_name"] = &schema.Schema{
		Description: "The name of the service",
		Required:    true,
		Type:        schema.TypeString,
	}
	element.Schema["start_mode"] = &schema.Schema{
		Default:     "auto",
		Description: "When will the service start. Can be auto, delayed-auto, manual, unchanged or an expression",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func expandDeployWindowsServiceAction(flattenedAction map[string]interface{}) *deployments.DeploymentAction {
	action := expandAction(flattenedAction)
	action.ActionType = "Octopus.WindowsService"

	addWindowsServiceToActionResource(flattenedAction, action)

	return action
}

func flattenWindowsService(properties map[string]core.PropertyValue) []interface{} {
	flattenedWindowsService := map[string]interface{}{}

	for propertyName, propertyValue := range properties {
		switch propertyName {
		case "Octopus.Action.WindowsService.Arguments":
			flattenedWindowsService["arguments"] = propertyValue.Value
		case "Octopus.Action.WindowsService.CreateOrUpdateService":
			createOrUpdateService, _ := strconv.ParseBool(propertyValue.Value)
			flattenedWindowsService["create_or_update_service"] = createOrUpdateService
		case "Octopus.Action.WindowsService.CustomAccountName":
			flattenedWindowsService["custom_account_name"] = propertyValue.Value
		case "Octopus.Action.WindowsService.CustomAccountPassword":
			flattenedWindowsService["custom_account_password"] = propertyValue.Value
		case "Octopus.Action.WindowsService.Dependencies":
			flattenedWindowsService["dependencies"] = propertyValue.Value
		case "Octopus.Action.WindowsService.Description":
			flattenedWindowsService["description"] = propertyValue.Value
		case "Octopus.Action.WindowsService.DisplayName":
			flattenedWindowsService["display_name"] = propertyValue.Value
		case "Octopus.Action.WindowsService.ExecutablePath":
			flattenedWindowsService["executable_path"] = propertyValue.Value
		case "Octopus.Action.WindowsService.ServiceAccount":
			flattenedWindowsService["service_account"] = propertyValue.Value
		case "Octopus.Action.WindowsService.ServiceName":
			flattenedWindowsService["service_name"] = propertyValue.Value
		case "Octopus.Action.WindowsService.StartMode":
			flattenedWindowsService["start_mode"] = propertyValue.Value
		}
	}

	return []interface{}{flattenedWindowsService}
}

func flattenDeployWindowsServiceAction(action *deployments.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenAction(action)

	for propertyName, propertyValue := range action.Properties {
		switch propertyName {
		case "Octopus.Action.WindowsService.Arguments":
			flattenedAction["arguments"] = propertyValue.Value
		case "Octopus.Action.WindowsService.CreateOrUpdateService":
			createOrUpdateService, _ := strconv.ParseBool(propertyValue.Value)
			flattenedAction["create_or_update_service"] = createOrUpdateService
		case "Octopus.Action.WindowsService.CustomAccountName":
			flattenedAction["custom_account_name"] = propertyValue.Value
		case "Octopus.Action.WindowsService.CustomAccountPassword":
			flattenedAction["custom_account_password"] = propertyValue.Value
		case "Octopus.Action.WindowsService.Dependencies":
			flattenedAction["dependencies"] = propertyValue.Value
		case "Octopus.Action.WindowsService.Description":
			flattenedAction["description"] = propertyValue.Value
		case "Octopus.Action.WindowsService.DisplayName":
			flattenedAction["display_name"] = propertyValue.Value
		case "Octopus.Action.WindowsService.ExecutablePath":
			flattenedAction["executable_path"] = propertyValue.Value
		case "Octopus.Action.WindowsService.ServiceAccount":
			flattenedAction["service_account"] = propertyValue.Value
		case "Octopus.Action.WindowsService.ServiceName":
			flattenedAction["service_name"] = propertyValue.Value
		case "Octopus.Action.WindowsService.StartMode":
			flattenedAction["start_mode"] = propertyValue.Value
		}
	}

	return flattenedAction
}

func addWindowsServiceFeatureToActionResource(tfAction map[string]interface{}, action *deployments.DeploymentAction) {
	if windowsServiceList, ok := tfAction["windows_service"]; ok {
		tfWindowsService := windowsServiceList.(*schema.Set).List()
		if len(tfWindowsService) > 0 {
			addWindowsServiceToActionResource(tfWindowsService[0].(map[string]interface{}), action)
		}
	}
}

func addWindowsServiceToActionResource(flattenedAction map[string]interface{}, action *deployments.DeploymentAction) {
	if len(action.Properties["Octopus.Action.EnabledFeatures"].Value) == 0 {
		action.Properties["Octopus.Action.EnabledFeatures"] = core.NewPropertyValue("Octopus.Features.WindowsService", false)
	} else if !strings.Contains(action.Properties["Octopus.Action.EnabledFeatures"].Value, "Octopus.Features.WindowsService") {
		actionPropertyValue := action.Properties["Octopus.Action.EnabledFeatures"].Value + ",Octopus.Features.WindowsService"
		action.Properties["Octopus.Action.EnabledFeatures"] = core.NewPropertyValue(actionPropertyValue, false)
	}

	if createOrUpdateService, ok := flattenedAction["create_or_update_service"]; ok {
		action.Properties["Octopus.Action.WindowsService.CreateOrUpdateService"] = core.NewPropertyValue(cases.Title(language.Und, cases.NoLower).String(strconv.FormatBool(createOrUpdateService.(bool))), false)
	}

	action.Properties["Octopus.Action.WindowsService.ServiceName"] = core.NewPropertyValue(flattenedAction["service_name"].(string), false)

	displayName := flattenedAction["display_name"]
	if displayName != nil {
		action.Properties["Octopus.Action.WindowsService.DisplayName"] = core.NewPropertyValue(displayName.(string), false)
	}

	description := flattenedAction["description"]
	if description != nil {
		action.Properties["Octopus.Action.WindowsService.Description"] = core.NewPropertyValue(description.(string), false)
	}

	action.Properties["Octopus.Action.WindowsService.ExecutablePath"] = core.NewPropertyValue(flattenedAction["executable_path"].(string), false)

	args := flattenedAction["arguments"]
	if args != nil {
		action.Properties["Octopus.Action.WindowsService.Arguments"] = core.NewPropertyValue(args.(string), false)
	}

	action.Properties["Octopus.Action.WindowsService.ServiceAccount"] = core.NewPropertyValue(flattenedAction["service_account"].(string), false)

	accountName := flattenedAction["custom_account_name"]
	if accountName != nil {
		action.Properties["Octopus.Action.WindowsService.CustomAccountName"] = core.NewPropertyValue(accountName.(string), false)
	}

	accountPassword := flattenedAction["custom_account_password"]
	if accountPassword != nil {
		action.Properties["Octopus.Action.WindowsService.CustomAccountPassword"] = core.NewPropertyValue(accountPassword.(string), false)
	}

	action.Properties["Octopus.Action.WindowsService.StartMode"] = core.NewPropertyValue(flattenedAction["start_mode"].(string), false)

	dependencies := flattenedAction["dependencies"]
	if dependencies != nil {
		action.Properties["Octopus.Action.WindowsService.Dependencies"] = core.NewPropertyValue(dependencies.(string), false)
	}
}
