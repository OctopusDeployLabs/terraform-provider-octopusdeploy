package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getApplyTerraformActionSchema() *schema.Schema {

	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	addPrimaryPackageSchema(element, false)

	element.Schema[constAdditionalInitParams] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Additional parameters passed to the init command",
		Optional:    true,
	}

	return actionSchema
}

func buildApplyTerraformActionResource(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	resource := expandDeploymentAction(tfAction)

	resource.ActionType = "Octopus.TerraformApply"
	resource.Properties["Octopus.Action.Terraform.AdditionalInitParams"] = tfAction[constAdditionalInitParams].(string)
	resource.Properties["Octopus.Action.Terraform.AllowPluginDownloads"] = "True"
	resource.Properties["Octopus.Action.Terraform.ManagedAccount"] = "None"

	resource.Properties["Octopus.Action.Script.ScriptSource"] = "Package"

	return resource
}
