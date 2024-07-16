package octopusdeploy

import (
	"encoding/json"
	"strconv"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandDeployKubernetesSecretAction(flattenedAction map[string]interface{}) *deployments.DeploymentAction {
	action := expandAction(flattenedAction)
	action.ActionType = "Octopus.KubernetesDeploySecret"

	action.Properties["Octopus.Action.KubernetesContainers.SecretName"] = core.NewPropertyValue(flattenedAction["secret_name"].(string), false)

	if tfSecretValues, ok := flattenedAction["secret_values"]; ok {
		secretValues := tfSecretValues.(map[string]interface{})

		j, _ := json.Marshal(secretValues)
		action.Properties["Octopus.Action.KubernetesContainers.SecretValues"] = core.NewPropertyValue(string(j), false)
	}

	if v, ok := flattenedAction["kubernetes_object_status_check_enabled"]; ok {
		action.Properties["Octopus.Action.Kubernetes.ResourceStatusCheck"] = core.NewPropertyValue(formatBoolForActionProperty(v.(bool)), false)
	}

	return action
}

func flattenDeployKubernetesSecretAction(action *deployments.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenAction(action)

	if v, ok := action.Properties["Octopus.Action.RunOnServer"]; ok {
		runOnServer, _ := strconv.ParseBool(v.Value)
		flattenedAction["run_on_server"] = runOnServer
	}

	if v, ok := action.Properties["Octopus.Action.KubernetesContainers.SecretName"]; ok {
		flattenedAction["secret_name"] = v.Value
	}

	if v, ok := action.Properties["Octopus.Action.Kubernetes.ResourceStatusCheck"]; ok {
		statusCheckEnabled, _ := strconv.ParseBool(v.Value)
		flattenedAction["kubernetes_object_status_check_enabled"] = statusCheckEnabled
	}

	if len(action.WorkerPool) > 0 {
		flattenedAction["worker_pool_id"] = action.WorkerPool
	}

	if len(action.WorkerPoolVariable) > 0 {
		flattenedAction["worker_pool_variable"] = action.WorkerPoolVariable
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
	addWorkerPoolSchema(element)
	addWorkerPoolVariableSchema(element)
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

	element.Schema["kubernetes_object_status_check_enabled"] = &schema.Schema{
		Optional:    true,
		Default:     true,
		Type:        schema.TypeBool,
		Description: "Indicates the status of the Kubernetes Object Status feature",
	}

	return actionSchema
}
