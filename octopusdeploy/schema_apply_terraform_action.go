package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandApplyTerraformAction(flattenedAction map[string]interface{}) octopusdeploy.DeploymentAction {
	action := expandAction(flattenedAction)
	action.ActionType = "Octopus.TerraformApply"

	action.Properties["Octopus.Action.Terraform.AdditionalInitParams"] = flattenedAction["additional_init_params"].(string)
	action.Properties["Octopus.Action.Terraform.AllowPluginDownloads"] = "True"
	action.Properties["Octopus.Action.Terraform.ManagedAccount"] = "None"
	action.Properties["Octopus.Action.Script.ScriptSource"] = "Package"

	return action
}

func flattenApplyTerraformAction(action octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenAction(action)
	flattenApplyTerraformActionProperties(flattenedAction, action.Properties)

	return flattenedAction
}

func flattenApplyTerraformActionProperties(actionMap map[string]interface{}, properties map[string]string) {
	for propertyName, propertyValue := range properties {
		switch propertyName {
		case "Octopus.Action.RunOnServer":
			runOnServer, _ := strconv.ParseBool(propertyValue)
			actionMap["run_on_server"] = runOnServer
		case "Octopus.Action.Terraform.AdditionalInitParams":
			actionMap["additional_init_params"] = propertyValue
		case "Octopus.Action.Terraform.AllowPluginDownloads":
			allowPluginDownloads, _ := strconv.ParseBool(propertyValue)
			actionMap["allow_plugin_downloads"] = allowPluginDownloads
		case "Octopus.Action.Terraform.ManagedAccount":
			actionMap["managed_account"] = propertyValue
		}
	}
}

func getApplyTerraformActionSchema() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	addPrimaryPackageSchema(element, false)

	element.Schema["additional_init_params"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Additional parameters passed to the init command",
		Optional:    true,
	}

	element.Schema["allow_plugin_downloads"] = &schema.Schema{
		Computed: true,
		Type:     schema.TypeBool,
		Optional: true,
	}

	element.Schema["managed_account"] = &schema.Schema{
		Computed: true,
		Type:     schema.TypeString,
		Optional: true,
	}

	return actionSchema
}
