package octopusdeploy

import (
	"encoding/json"

	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getDeployKubernetesSecretActionSchema() *schema.Schema {

	actionSchema, element := getCommonDeploymentActionSchema()
	addExecutionLocationSchema(element)
	element.Schema[constSecretName] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The name of the secret resource",
		Required:    true,
	}

	element.Schema[constSecretValues] = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				constKey: {
					Type:     schema.TypeString,
					Required: true,
				},
				constValue: {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}

	return actionSchema
}

func buildDeployKubernetesSecretActionResource(tfAction map[string]interface{}) model.DeploymentAction {
	resource := buildDeploymentActionResource(tfAction)

	resource.ActionType = "Octopus.KubernetesDeploySecret"

	resource.Properties["Octopus.Action.KubernetesContainers.SecretName"] = tfAction[constSecretValues].(string)

	if tfSecretValues, ok := tfAction[constSecretValues]; ok {

		secretValues := make(map[string]string)

		for _, tfSecretValue := range tfSecretValues.([]interface{}) {
			tfSecretValueTyped := tfSecretValue.(map[string]interface{})
			secretValues[tfSecretValueTyped[constKey].(string)] = tfSecretValueTyped[constValue].(string)
		}

		j, _ := json.Marshal(secretValues)

		resource.Properties["Octopus.Action.KubernetesContainers.SecretValues"] = string(j)
	}

	return resource
}
