package actions

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RunScriptActionMapper struct{}

var _ MappableAction = &RunScriptActionMapper{}

func (r RunScriptActionMapper) ToState(ctx context.Context, actionState attr.Value, action *deployments.DeploymentAction, newAction map[string]attr.Value) diag.Diagnostics {
	diag := mapBaseDeploymentActionToState(ctx, actionState, action, newAction)
	if diag.HasError() {
		return diag
	}

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

func (r RunScriptActionMapper) ToDeploymentAction(actionAttribute attr.Value) *deployments.DeploymentAction {
	actionAttrs := GetActionAttributes(actionAttribute)
	if actionAttrs == nil {
		return nil
	}

	action := getBaseAction(actionAttribute)
	if action == nil {
		return nil
	}

	mapAttributeToProperty(action, actionAttrs, "script_file_name", "Octopus.Action.Script.ScriptFileName")
	mapAttributeToProperty(action, actionAttrs, "script_body", "Octopus.Action.Script.ScriptBody")
	mapAttributeToProperty(action, actionAttrs, "script_parameters", "Octopus.Action.Script.ScriptParameters")
	mapAttributeToProperty(action, actionAttrs, "script_source", "Octopus.Action.Script.ScriptSource")
	mapAttributeToProperty(action, actionAttrs, "script_syntax", "Octopus.Action.Script.Syntax")

	if variableSubstitutionInFiles, ok := actionAttrs["variable_substitution_in_files"]; ok {
		action.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"] = core.NewPropertyValue(variableSubstitutionInFiles.(types.String).ValueString(), false)
		action.Properties["Octopus.Action.SubstituteInFiles.Enabled"] = core.NewPropertyValue(formatBoolForActionProperty(true), false)

		ensureFeatureIsEnabled(action, "Octopus.Action.SubstituteInFiles")
	}

	return action
}
