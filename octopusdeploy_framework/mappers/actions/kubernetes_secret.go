package actions

import (
	"context"
	"encoding/json"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KubernetesSecretActionMapper struct{}

var _ MappableAction = &KubernetesSecretActionMapper{}

func (k KubernetesSecretActionMapper) ToState(ctx context.Context, actionState attr.Value, action *deployments.DeploymentAction, newAction map[string]attr.Value) diag.Diagnostics {
	diag := mapBaseDeploymentActionToState(ctx, actionState, action, newAction)
	if diag.HasError() {
		return diag
	}

	mapPropertyToStateBool(action, newAction, "Octopus.Action.RunOnServer", "run_on_server", false)
	mapPropertyToStateString(action, newAction, "Octopus.Action.KubernetesContainers.SecretName", "secret_name")
	mapPropertyToStateBool(action, newAction, "Octopus.Action.Kubernetes.ResourceStatusCheck", "kubernetes_object_status_check_enabled", false)
	newAction["worker_pool_id"] = types.StringValue(action.WorkerPool)
	newAction["worker_pool_variable"] = types.StringValue(action.WorkerPoolVariable)

	if v, ok := action.Properties["Octopus.Action.KubernetesContainers.SecretValues"]; ok {
		var secretKeyValues map[string]string
		json.Unmarshal([]byte(v.Value), &secretKeyValues)
		mappedSecrets := make(map[string]attr.Value)
		for key, value := range secretKeyValues {
			mappedSecrets[key] = types.StringValue(value)
		}
		newAction["secret_values"], diag = types.MapValue(types.StringType, mappedSecrets)
		if diag.HasError() {
			return diag
		}
	} else {
		newAction["secret_values"] = types.MapNull(types.StringType)
	}

	return nil
}

func (k KubernetesSecretActionMapper) ToDeploymentAction(actionAttribute attr.Value) *deployments.DeploymentAction {
	actionAttrs := GetActionAttributes(actionAttribute)
	if actionAttrs == nil {
		return nil
	}

	action := GetBaseAction(actionAttribute)
	if action == nil {
		return nil
	}

	mapAttributeToProperty(action, actionAttrs, "secret_name", "Octopus.Action.KubernetesContainers.SecretName")
	mapBooleanAttributeToProperty(action, actionAttrs, "kubernetes_object_status_check_enabled", "Octopus.Action.Kubernetes.ResourceStatusCheck")

	if attrValue, ok := actionAttrs["secret_values"]; ok {
		secretValues := attrValue.(types.Map)
		mappedValues := make(map[string]string)
		for key, value := range secretValues.Elements() {
			mappedValues[key] = value.(types.String).ValueString()
		}

		j, _ := json.Marshal(mappedValues)
		action.Properties["Octopus.Action.KubernetesContainers.SecretValues"] = core.NewPropertyValue(string(j), false)
	}

	return action
}
