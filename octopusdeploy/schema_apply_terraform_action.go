package octopusdeploy

import (
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

func expandApplyTerraformAction(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	resource := expandDeploymentAction(tfAction)

	resource.ActionType = "Octopus.TerraformApply"
	resource.Properties["Octopus.Action.Terraform.AdditionalInitParams"] = tfAction["additional_init_params"].(string)
	resource.Properties["Octopus.Action.Terraform.AllowPluginDownloads"] = "True"
	resource.Properties["Octopus.Action.Terraform.ManagedAccount"] = "None"
	resource.Properties["Octopus.Action.Script.ScriptSource"] = "Package"

	return resource
}

func flattenApplyTerraformAction(deploymentAction octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedApplyTerraformAction := map[string]interface{}{
		"can_be_used_for_project_versioning": deploymentAction.CanBeUsedForProjectVersioning,
		"channels":                           deploymentAction.Channels,
		"condition":                          deploymentAction.Condition,
		"container":                          flattenDeploymentActionContainer(deploymentAction.Container),
		"environments":                       deploymentAction.Environments,
		"excluded_environments":              deploymentAction.ExcludedEnvironments,
		"id":                                 deploymentAction.ID,
		"is_disabled":                        deploymentAction.IsDisabled,
		"is_required":                        deploymentAction.IsRequired,
		"name":                               deploymentAction.Name,
		"notes":                              deploymentAction.Notes,
		"package":                            flattenPackageReferences(deploymentAction.Packages),
		"properties":                         deploymentAction.Properties,
		"tenant_tags":                        deploymentAction.TenantTags,
	}

	flattenApplyTerraformActionProperties(flattenedApplyTerraformAction, deploymentAction.Properties)

	return flattenedApplyTerraformAction
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
