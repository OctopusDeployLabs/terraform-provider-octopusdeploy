package octopusdeploy

import (
	"encoding/json"
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandDeployKubernetesSecretAction(flattenedAction map[string]interface{}) octopusdeploy.DeploymentAction {
	action := expandAction(flattenedAction)
	action.ActionType = "Octopus.KubernetesDeploySecret"

	action.Properties["Octopus.Action.KubernetesContainers.SecretName"] = octopusdeploy.NewPropertyValue(flattenedAction["secret_name"].(string), false)

	if tfSecretValues, ok := flattenedAction["secret_values"]; ok {
		secretValues := tfSecretValues.(map[string]interface{})

		j, _ := json.Marshal(secretValues)
		action.Properties["Octopus.Action.KubernetesContainers.SecretValues"] = octopusdeploy.NewPropertyValue(string(j), false)
	}

	return action
}

func flattenDeployKubernetesSecretAction(action octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenAction(action)

	if v, ok := action.Properties["Octopus.Action.RunOnServer"]; ok {
		runOnServer, _ := strconv.ParseBool(v.Value)
		flattenedAction["run_on_server"] = runOnServer
	}

	if v, ok := action.Properties["Octopus.Action.KubernetesContainers.SecretName"]; ok {
		flattenedAction["secret_name"] = v.Value
	}

	if v, ok := action.Properties["Octopus.Action.KubernetesContainers.SecretValues"]; ok {
		var secretKeyValues map[string]string
		json.Unmarshal([]byte(v.Value), &secretKeyValues)

		flattenedAction["secret_values"] = secretKeyValues
	}

	return flattenedAction
}

func getDeployKubernetesSecretActionSchema() *schema.Schema {
	actionSchema, element := getActionSchema()
	addExecutionLocationSchema(element)
	element.Schema["secret_name"] = &schema.Schema{
		Description: "The name of the secret resource",
		Required:    true,
		Type:        schema.TypeString,
	}

	element.Schema["secret_values"] = &schema.Schema{
		Elem:     &schema.Schema{Type: schema.TypeString},
		Required: true,
		Type:     schema.TypeMap,
	}

	return actionSchema
}
