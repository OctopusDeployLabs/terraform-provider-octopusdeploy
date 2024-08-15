package actions

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KubectlScriptActionMapper struct{}

func (k KubectlScriptActionMapper) ToState(ctx context.Context, actionState attr.Value, action *deployments.DeploymentAction, newAction map[string]attr.Value) diag.Diagnostics {
	diag := mapBaseDeploymentActionToState(ctx, actionState, action, newAction)
	if diag.HasError() {
		return diag
	}

	mapPropertyToStateString(action, newAction, "Octopus.Action.KubernetesContainers.Namespace", "namespace")

	newAction["worker_pool_id"] = types.StringValue(action.WorkerPool)
	newAction["worker_pool_variable"] = types.StringValue(action.WorkerPoolVariable)

	mapPropertyToStateBool(action, newAction, "Octopus.Action.RunOnServer", "run_on_server", false)
	mapPropertyToStateString(action, newAction, "Octopus.Action.Script.ScriptBody", "script_body")
	mapPropertyToStateString(action, newAction, "Octopus.Action.Script.ScriptFileName", "script_file_name")
	mapPropertyToStateString(action, newAction, "Octopus.Action.Script.ScriptSource", "script_source")
	mapPropertyToStateString(action, newAction, "Octopus.Action.Script.ScriptParameters", "script_parameters")
	mapPropertyToStateString(action, newAction, "Octopus.Action.Script.Syntax", "script_syntax")
	mapPropertyToStateString(action, newAction, "Octopus.Action.Script.ScriptFileName", "script_file_name")
	mapPropertyToStateString(action, newAction, "Octopus.Action.SubstituteInFiles.TargetFiles", "variable_substitution_in_files")

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
