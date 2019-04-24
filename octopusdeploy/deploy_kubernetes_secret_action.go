package octopusdeploy

import (
	"encoding/json"
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func getDeployKubernetesSecretActionSchema() *schema.Schema {

	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	element.Schema["secret_name"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The name of the secret resource",
		Required:    true,
	}

	element.Schema["secret_values"] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:     schema.TypeString,
					Required: true,
				},
				"value": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}

	return actionSchema
}

func buildDeployKubernetesSecretActionResource(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	resource := buildDeploymentActionResource(tfAction)

	resource.ActionType = "Octopus.KubernetesDeploySecret"

	resource.Properties["Octopus.Action.KubernetesContainers.SecretName"] = tfAction["secret_name"].(string)

	if tfSecretValues, ok := tfAction["secret_values"]; ok {

		secretValues := make(map[string]string)

		for _, tfSecretValue := range tfSecretValues.([]interface{}) {
			tfSecretValueTyped := tfSecretValue.(map[string]interface{})
			secretValues[tfSecretValueTyped["key"].(string)] = tfSecretValueTyped["value"].(string)
		}

		j, _ := json.Marshal(secretValues)

		resource.Properties["Octopus.Action.KubernetesContainers.SecretValues"] = string(j)
	}

	return resource
}
