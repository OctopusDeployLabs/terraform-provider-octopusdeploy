package octopusdeploy

import (
	"encoding/json"
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandDeployKubernetesSecretAction(flattenedAction map[string]interface{}) octopusdeploy.DeploymentAction {
	action := expandDeploymentAction(flattenedAction)
	action.ActionType = "Octopus.KubernetesDeploySecret"

	action.Properties["Octopus.Action.KubernetesContainers.SecretName"] = flattenedAction["secret_name"].(string)

	if tfSecretValues, ok := flattenedAction["secret_values"]; ok {
		secretValues := make(map[string]string)

		for _, tfSecretValue := range tfSecretValues.([]interface{}) {
			tfSecretValueTyped := tfSecretValue.(map[string]interface{})
			secretValues[tfSecretValueTyped["key"].(string)] = tfSecretValueTyped["value"].(string)
		}

		j, _ := json.Marshal(secretValues)
		action.Properties["Octopus.Action.KubernetesContainers.SecretValues"] = string(j)
	}

	return action
}

func flattenDeployKubernetesSecretAction(action octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenCommonDeploymentAction(action)

	if v, ok := action.Properties["Octopus.Action.RunOnServer"]; ok {
		runOnServer, _ := strconv.ParseBool(v)
		flattenedAction["run_on_server"] = runOnServer
	}

	if v, ok := action.Properties["Octopus.Action.KubernetesContainers.SecretName"]; ok {
		flattenedAction["secret_name"] = v
	}

	if v, ok := action.Properties["Octopus.Action.KubernetesContainers.SecretValues"]; ok {
		var secretKeyValues map[string]string
		json.Unmarshal([]byte(v), &secretKeyValues)

		flattenedSecretKeyValues := []interface{}{}
		for secretKey, secretValue := range secretKeyValues {
			flattenedSecretKeyValues = append(flattenedSecretKeyValues, map[string]interface{}{
				"key":   secretKey,
				"value": secretValue,
			})
		}

		flattenedAction["secret_values"] = flattenedSecretKeyValues
	}

	return flattenedAction
}

func getDeployKubernetesSecretActionSchema() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	element.Schema["secret_name"] = &schema.Schema{
		Description: "The name of the secret resource",
		Required:    true,
		Type:        schema.TypeString,
	}

	element.Schema["secret_values"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Required: true,
					Type:     schema.TypeString,
				},
				"value": {
					Required: true,
					Type:     schema.TypeString,
				},
			},
		},
	}

	return actionSchema
}
