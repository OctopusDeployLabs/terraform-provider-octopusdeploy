package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandApplyTerraformTemplateAction(flattenedAction map[string]interface{}) octopusdeploy.DeploymentAction {
	action := expandAction(flattenedAction)
	action.ActionType = "Octopus.TerraformApply"

	if v, ok := flattenedAction["additional_init_params"]; ok {
		action.Properties["Octopus.Action.Terraform.AdditionalInitParams"] = octopusdeploy.NewPropertyValue(v.(string), false)
	}

	if v, ok := flattenedAction["allow_plugin_downloads"]; ok {
		allowPluginDownloads := v.(bool)
		action.Properties["Octopus.Action.Terraform.AllowPluginDownloads"] = octopusdeploy.NewPropertyValue(strconv.FormatBool(allowPluginDownloads), false)
	}

	if v, ok := flattenedAction["managed_account"]; ok {
		action.Properties["Octopus.Action.Terraform.ManagedAccount"] = octopusdeploy.NewPropertyValue(v.(string), false)
	}

	action.Properties["Octopus.Action.Script.ScriptSource"] = octopusdeploy.NewPropertyValue("Package", false)

	return action
}

func flattenApplyTerraformTemplateAction(action octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenAction(action)

	for propertyName, propertyValue := range action.Properties {
		switch propertyName {
		case "Octopus.Action.RunOnServer":
			runOnServer, _ := strconv.ParseBool(propertyValue.Value)
			flattenedAction["run_on_server"] = runOnServer
		case "Octopus.Action.Terraform.AdditionalInitParams":
			flattenedAction["additional_init_params"] = propertyValue.Value
		case "Octopus.Action.Terraform.AllowPluginDownloads":
			allowPluginDownloads, _ := strconv.ParseBool(propertyValue.Value)
			flattenedAction["allow_plugin_downloads"] = allowPluginDownloads
		case "Octopus.Action.Terraform.ManagedAccount":
			flattenedAction["managed_account"] = propertyValue.Value
		}
	}

	return flattenedAction
}

func getApplyTerraformTemplateActionSchema() *schema.Schema {
	actionSchema, element := getActionSchema()
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
