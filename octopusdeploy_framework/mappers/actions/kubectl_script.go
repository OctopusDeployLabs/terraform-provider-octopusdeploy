package actions

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type KubectlScriptActionMapper struct{}

func (k KubectlScriptActionMapper) ToState(ctx context.Context, action *deployments.DeploymentAction, newAction map[string]attr.Value) diag.Diagnostics {
	diag := mapBaseDeploymentActionToState(ctx, action, newAction)
	if diag.HasError() {
		return diag
	}

	mapPropertyToStateString(action, newAction, "Octopus.Action.KubernetesContainers.Namespace", "namespace")

	return nil
}

func (k KubectlScriptActionMapper) ToDeploymentAction(actionAttribute attr.Value) *deployments.DeploymentAction {
	actionAttrs := GetActionAttributes(actionAttribute)
	if actionAttrs == nil {
		return nil
	}

	runscriptMapper := RunScriptActionMapper{}
	action := runscriptMapper.ToDeploymentAction(actionAttribute)
	if action == nil {
		return nil
	}
	action.ActionType = "Octopus.KubernetesRunScript"
	mapAttributeToProperty(action, actionAttrs, "namespace", "Octopus.Action.KubernetesContainers.Namespace")
	return action
}

var _ MappableAction = &KubectlScriptActionMapper{}
