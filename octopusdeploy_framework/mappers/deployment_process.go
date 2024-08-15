package mappers

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/mappers/actions"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"sort"
	"strings"
)

var ActionMappers = map[string]actions.MappableAction{
	schemas.DeploymentProcessAction:                      &actions.Action{},
	schemas.DeploymentProcessRunScriptAction:             &actions.RunScriptActionMapper{},
	schemas.DeploymentProcessRunKubectlScriptAction:      &actions.KubectlScriptActionMapper{},
	schemas.DeploymentProcessPackageAction:               &actions.PackageActionMapper{},
	schemas.DeploymentProcessWindowsServiceAction:        &actions.WindowsServiceActionMapper{},
	schemas.DeploymentProcessManualInterventionAction:    &actions.ManualInterventionActionMapper{},
	schemas.DeploymentProcessApplyKubernetesSecretAction: &actions.KubernetesSecretActionMapper{},
	//schemas.DeploymentProcessApplyTerraformTemplateAction: "Octopus.TerraformApply",
}

func MapDeploymentProcessToState(ctx context.Context, deploymentProcess *deployments.DeploymentProcess, state *schemas.DeploymentProcessResourceModel) diag.Diagnostics {
	state.ID = types.StringValue(deploymentProcess.ID)
	state.Branch = types.StringValue(deploymentProcess.Branch)
	state.ProjectID = types.StringValue(deploymentProcess.ProjectID)
	state.SpaceID = types.StringValue(deploymentProcess.SpaceID)
	state.Version = types.StringValue(fmt.Sprintf("%d", deploymentProcess.Version))
	state.LastSnapshotID = types.StringValue(deploymentProcess.LastSnapshotID)

	return mapStepsToState(ctx, state, deploymentProcess)
}

func MapStateToDeploymentProcess(ctx context.Context, state *schemas.DeploymentProcessResourceModel, deploymentProcess *deployments.DeploymentProcess) diag.Diagnostics {
	// this should not map the version number from the schema
	deploymentProcess.Branch = state.Branch.ValueString()
	deploymentProcess.SpaceID = state.SpaceID.ValueString()
	deploymentProcess.LastSnapshotID = state.LastSnapshotID.ValueString()
	deploymentProcess.ProjectID = state.ProjectID.ValueString()
	mapStepsToDeploymentProcess(ctx, state.Steps, deploymentProcess)

	return nil
}

//func mapStepsToState(ctx context.Context, state *schemas.DeploymentProcessResourceModel, process *deployments.DeploymentProcess) diag.Diagnostics {
//	// this is supposed to be updating the in-memory state
//	if process.Steps == nil || len(process.Steps) == 0 {
//		//return types.ListNull(types.ObjectType{AttrTypes: getStepTypeAttrs()}), nil
//		return nil
//	}
//
//	//steps := make([]attr.Value, 0)
//
//	for _, step := range state.Steps.Elements() {
//		stepAttrs := step.(types.Object).Attributes()
//		name := stepAttrs["name"].(types.String).ValueString()
//		var currentStep *deployments.DeploymentStep
//		for _, processStep := range process.Steps {
//			if processStep.Name == name {
//				currentStep = processStep
//				break
//			}
//		}
//
//		stepAttrs["id"] = types.StringValue(currentStep.ID)
//
//		for actionKey, _ := range ActionMappers {
//			for _, action := range stepAttrs[actionKey].(types.List).Elements() {
//				actionAttrs := action.(types.Object).Attributes()
//				actionName := actionAttrs["name"].(types.String).ValueString()
//				var currentAction *deployments.DeploymentAction
//				for _, stepAction := range currentStep.Actions {
//					if stepAction.Name == actionName {
//						currentAction = stepAction
//						break
//					}
//				}
//
//				actionAttrs["id"] = types.StringValue(currentAction.ID)
//				actionAttrs["channels"] = types.ListValueMust(types.StringType, []attr.Value{})
//			}
//		}
//
//		//steps = append(steps, types.ObjectValueMust(getStepTypeAttrs(), stepAttrs))
//		tflog.Debug(ctx, name)
//	}
//
//	//return types.ListValueMust(types.ObjectType{AttrTypes: getStepTypeAttrs()}, steps), nil
//	return nil
//}

func mapStepsToState(ctx context.Context, state *schemas.DeploymentProcessResourceModel, process *deployments.DeploymentProcess) diag.Diagnostics {
	// this is creating a new set of nested resources, but only mapping directly to the attributes
	// the attribute updates seem to work but are not being maintained when the method returns back to the resource_deployment_process file

	if process.Steps == nil || len(process.Steps) == 0 {
		return nil
	}

	var steps []attr.Value
	for _, deploymentStep := range process.Steps {
		stepFromState := getStepFromState(deploymentStep, state)
		properties, diags := actions.MapPropertiesToState(ctx, deploymentStep.Properties)
		if diags.HasError() {
			return diags
		}
		newStep := map[string]attr.Value{
			"id":                  types.StringValue(deploymentStep.ID),
			"condition":           types.StringValue(string(deploymentStep.Condition)),
			"name":                types.StringValue(deploymentStep.Name),
			"package_requirement": types.StringValue(string(deploymentStep.PackageRequirement)),
			"properties":          properties,
			"start_trigger":       types.StringValue(string(deploymentStep.StartTrigger)),
		}

		newStep[schemas.DeploymentProcessWindowSize] = types.StringValue("")
		newStep[schemas.DeploymentProcessConditionExpression] = types.StringValue("")
		newStep[schemas.DeploymentProcessTargetRoles] = types.ListNull(types.StringType)
		for propertyName, propertyValue := range deploymentStep.Properties {
			switch propertyName {
			case "Octopus.Action.TargetRoles":
				newStep[schemas.DeploymentProcessTargetRoles] = util.FlattenStringList(strings.Split(propertyValue.Value, ","))
			case "Octopus.Action.MaxParallelism":
				newStep[schemas.DeploymentProcessWindowSize] = types.StringValue(propertyValue.Value)
			case "Octopus.Step.ConditionVariableExpression":
				newStep[schemas.DeploymentProcessConditionExpression] = types.StringValue(propertyValue.Value)
			}
		}

		newActions := make(map[string][]attr.Value)
		for i, a := range deploymentStep.Actions {
			actionFromState := getActionFromState(stepFromState, a)
			newAction := map[string]attr.Value{
				"computed_sort_order": types.Int64Value(int64(i)),
			}
			sortOrder := getSortOrderStateValue(deploymentStep, a, state)
			if sortOrder != nil {
				newAction["sort_order"] = sortOrder
			}

			srcAction := deploymentStep.Actions[i]
			terraformActionKeyName := getActionTypeTerraformAttributeName(srcAction.ActionType)

			d := ActionMappers[terraformActionKeyName].ToState(ctx, actionFromState, a, newAction)
			if d.HasError() {
				return d
			}
			//switch srcAction.ActionType {
			//	flatten_action_func("deploy_package_action", i, flattenDeployPackageAction)
			//case "Octopus.TerraformApply":
			//	flatten_action_func("apply_terraform_template_action", i, flattenApplyTerraformTemplateAction)

			//}

			if _, ok := newActions[terraformActionKeyName]; !ok {
				newActions[terraformActionKeyName] = make([]attr.Value, 0)
			}
			newActions[terraformActionKeyName] = append(newActions[terraformActionKeyName], types.ObjectValueMust(getActionTypeAttrs(terraformActionKeyName), newAction))

		}

		for actionAttributeName, _ := range schemas.ActionsAttributeToActionTypeMap {
			if len(newActions[actionAttributeName]) > 0 {
				newStep[actionAttributeName] = types.ListValueMust(types.ObjectType{AttrTypes: getActionTypeAttrs(actionAttributeName)}, newActions[actionAttributeName])
			} else {
				newStep[actionAttributeName] = types.ListNull(types.ObjectType{AttrTypes: getActionTypeAttrs(actionAttributeName)})
			}
		}

		mappedStep, diags := types.ObjectValue(getStepTypeAttrs(), newStep)
		if diags.HasError() {
			return diags
		}
		steps = append(steps, mappedStep)
	}

	updatedSteps, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: getStepTypeAttrs()}, steps)
	if diags.HasError() {
		return diags
	}

	state.Steps = updatedSteps
	return nil
}

func getStepFromState(step *deployments.DeploymentStep, state *schemas.DeploymentProcessResourceModel) attr.Value {
	for _, s := range state.Steps.Elements() {
		stepAttrs := s.(types.Object).Attributes()
		if step.Name == stepAttrs["name"].(types.String).ValueString() {
			return s
		}
	}

	return nil
}

func getSortOrderStateValue(step *deployments.DeploymentStep, a *deployments.DeploymentAction, state *schemas.DeploymentProcessResourceModel) attr.Value {
	s := getStepFromState(step, state)
	if s != nil {
		action := getActionFromState(s, a)
		actionAttrs := action.(types.Object).Attributes()
		if v, ok := actionAttrs["sort_order"]; ok {
			return v
		} else {
			return nil
		}
	}
	return nil

}

func getActionFromState(stepState attr.Value, a *deployments.DeploymentAction) attr.Value {
	stepAttrs := stepState.(types.Object).Attributes()
	for key, stepAttr := range stepAttrs {
		if isAction(key) {
			actions := stepAttr.(types.List)
			for _, action := range actions.Elements() {
				actionAttrs := action.(types.Object).Attributes()
				name := actionAttrs["name"].(types.String).ValueString()
				if a.Name == name {
					return action
				}
			}
		}
	}

	return nil
}

func getStepTypeAttrs() map[string]attr.Type {
	attrs := map[string]attr.Type{
		"id":                               types.StringType,
		"name":                             types.StringType,
		schemas.DeploymentProcessCondition: types.StringType,
		schemas.DeploymentProcessConditionExpression: types.StringType,
		schemas.DeploymentProcessPackageRequirement:  types.StringType,
		schemas.DeploymentProcessProperties:          types.MapType{ElemType: types.StringType},
		schemas.DeploymentProcessStartTrigger:        types.StringType,
		schemas.DeploymentProcessTargetRoles:         types.ListType{ElemType: types.StringType},
		schemas.DeploymentProcessWindowSize:          types.StringType,
	}

	for actionAttributeName, _ := range schemas.ActionsAttributeToActionTypeMap {
		attrs[actionAttributeName] = types.ListType{ElemType: types.ObjectType{AttrTypes: getActionTypeAttrs(actionAttributeName)}}
	}

	return attrs
}

func getActionTypeTerraformAttributeName(actionTypeName string) string {
	for actionAttributeName, actionType := range schemas.ActionsAttributeToActionTypeMap {
		if actionType == actionTypeName {
			return actionAttributeName
		}
	}

	return schemas.DeploymentProcessAction
}

func getActionTypeAttrs(actionType string) map[string]attr.Type {
	attrs := map[string]attr.Type{
		"can_be_used_for_project_versioning": types.BoolType,
		"channels":                           types.ListType{ElemType: types.StringType},
		"condition":                          types.StringType,
		"container":                          types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"feed_id": types.StringType, "image": types.StringType}}},
		"environments":                       types.ListType{ElemType: types.StringType},
		"excluded_environments":              types.ListType{ElemType: types.StringType},
		"features":                           types.ListType{ElemType: types.StringType},
		"action_template":                    types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"id": types.StringType, "version": types.StringType}}},
		"id":                                 types.StringType,
		"git_dependency":                     types.SetType{ElemType: types.ObjectType{AttrTypes: actions.GetGitDependencyAttrTypes()}},
		"is_disabled":                        types.BoolType,
		"is_required":                        types.BoolType,
		"name":                               types.StringType,
		"notes":                              types.StringType,
		"primary_package":                    types.ListType{ElemType: types.ObjectType{AttrTypes: actions.GetPackageReferenceAttrTypes(true)}},
		"package":                            types.ListType{ElemType: types.ObjectType{AttrTypes: actions.GetPackageReferenceAttrTypes(false)}},
		"properties":                         types.MapType{ElemType: types.StringType},
		"sort_order":                         types.Int64Type,
		"slug":                               types.StringType,
		"tenant_tags":                        types.ListType{ElemType: types.StringType},
		"computed_sort_order":                types.Int64Type,
	}
	switch actionType {
	case schemas.DeploymentProcessRunKubectlScriptAction:
		attrs["namespace"] = types.StringType
		attrs["script_body"] = types.StringType
		attrs["script_syntax"] = types.StringType
		attrs["worker_pool_id"] = types.StringType
		attrs["worker_pool_variable"] = types.StringType
		attrs["variable_substitution_in_files"] = types.StringType
		attrs["run_on_server"] = types.BoolType
		attrs["script_file_name"] = types.StringType
		attrs["script_parameters"] = types.StringType
		attrs["script_source"] = types.StringType
		attrs["script_file_name"] = types.StringType
		attrs["script_body"] = types.StringType
		attrs["script_syntax"] = types.StringType
	case schemas.DeploymentProcessRunScriptAction:
		attrs["run_on_server"] = types.BoolType
		attrs["script_file_name"] = types.StringType
		attrs["script_parameters"] = types.StringType
		attrs["script_source"] = types.StringType
		attrs["script_file_name"] = types.StringType
		attrs["script_body"] = types.StringType
		attrs["script_syntax"] = types.StringType
		attrs["worker_pool_id"] = types.StringType
		attrs["worker_pool_variable"] = types.StringType
		attrs["variable_substitution_in_files"] = types.StringType
		break
	case schemas.DeploymentProcessApplyTerraformTemplateAction:
		attrs["run_on_server"] = types.BoolType
		attrs["worker_pool_id"] = types.StringType
		attrs["worker_pool_variable"] = types.StringType
		attrs["advanced_options"] = types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"allow_additional_plugin_downloads": types.BoolType, "apply_parameters": types.StringType, "init_parameters": types.StringType, "plugin_cache_directory": types.StringType, "workspace": types.StringType}}}
		attrs["aws_account"] = types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"region": types.StringType, "variable": types.StringType, "use_instance_role": types.BoolType, "role": types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"arn": types.StringType, "external_id": types.StringType, "role_session_name": types.StringType, "session_duration": types.Int64Type}}}}}}
		attrs["azure_account"] = types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"variable": types.StringType}}}
		attrs["google_cloud_account"] = types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"variable": types.StringType, "use_vm_service_account": types.BoolType, "project": types.StringType, "region": types.StringType, "zone": types.StringType, "service_account_emails": types.StringType, "impersonate_service_account": types.BoolType}}}
		attrs["template"] = types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"additional_variable_files": types.StringType, "directory": types.StringType, "run_automatic_file_substitution": types.BoolType, "target_files": types.StringType}}}
		attrs["template_parameters"] = types.StringType
		attrs["inline_template"] = types.StringType
		break
	case schemas.DeploymentProcessApplyKubernetesSecretAction:
		attrs["run_on_server"] = types.BoolType
		attrs["worker_pool_id"] = types.StringType
		attrs["worker_pool_variable"] = types.StringType
		attrs["secret_name"] = types.StringType
		attrs["secret_values"] = types.MapType{ElemType: types.StringType}
		attrs["kubernetes_object_status_check_enabled"] = types.BoolType
		break
	case schemas.DeploymentProcessPackageAction:
		attrs["windows_service"] = types.SetType{ElemType: types.ObjectType{AttrTypes: getWindowsServiceAttrTypes()}}
		break
	case schemas.DeploymentProcessWindowsServiceAction:
		for k, v := range getWindowsServiceAttrTypes() {
			attrs[k] = v
		}
		break
	case schemas.DeploymentProcessManualInterventionAction:
		attrs["instructions"] = types.StringType
		attrs["responsible_teams"] = types.StringType
		break
	default:
		attrs["action_type"] = types.StringType
		attrs["run_on_server"] = types.BoolType
		attrs["worker_pool_id"] = types.StringType
		attrs["worker_pool_variable"] = types.StringType
	}

	return attrs
}

func getWindowsServiceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"arguments":                types.StringType,
		"create_or_update_service": types.BoolType,
		"custom_account_name":      types.StringType,
		"custom_account_password":  types.StringType,
		"dependencies":             types.StringType,
		"description":              types.StringType,
		"display_name":             types.StringType,
		"executable_path":          types.StringType,
		"service_account":          types.StringType,
		"service_name":             types.StringType,
		"start_mode":               types.StringType,
	}
}

func mapStepsToDeploymentProcess(ctx context.Context, steps types.List, current *deployments.DeploymentProcess) {
	if steps.IsNull() || steps.IsUnknown() {
		return
	}

	for _, s := range steps.Elements() {
		attrs := s.(types.Object).Attributes()
		step := deployments.NewDeploymentStep(attrs["name"].(types.String).ValueString())
		step.ID = attrs["id"].(types.String).ValueString()
		step.Condition = deployments.DeploymentStepConditionType(attrs[schemas.DeploymentProcessCondition].(types.String).ValueString())
		if conditionExpression, ok := attrs[schemas.DeploymentProcessConditionExpression]; ok {
			step.Properties["Octopus.Step.ConditionVariableExpression"] = core.NewPropertyValue(conditionExpression.(types.String).ValueString(), false)
		}
		if packageRequirement, ok := attrs["package_requirement"]; ok {
			step.PackageRequirement = deployments.DeploymentStepPackageRequirement(packageRequirement.(types.String).ValueString())
		}
		if startTrigger, ok := attrs["start_trigger"]; ok {
			step.StartTrigger = deployments.DeploymentStepStartTrigger(startTrigger.(types.String).ValueString())
		}

		if targetRoles, ok := attrs["target_roles"]; ok {
			roles := targetRoles.(types.List)
			step.Properties["Octopus.Action.TargetRoles"] = core.NewPropertyValue(strings.Join(util.ExpandStringList(roles), ","), false)
		}

		if windowSize, ok := attrs["window_size"]; ok {
			step.Properties["Octopus.Action.MaxParallelism"] = core.NewPropertyValue(windowSize.(types.String).ValueString(), false)
		}

		var sort_order map[string]int64 = make(map[string]int64)
		for key, attributes := range attrs {
			if attributes.IsNull() {
				continue
			}

			actionMapping := func(mappingFunc func(attr attr.Value) *deployments.DeploymentAction) {
				action := mappingFunc(attributes)
				if action.ActionType == "" {
					action.ActionType = schemas.ActionsAttributeToActionTypeMap[key]
				}

				step.Actions = append(step.Actions, action)
				actionAttrs := actions.GetActionAttributes(attributes)
				if posn, ok := actionAttrs["sort_order"].(types.Int64); ok && !posn.IsNull() && posn.ValueInt64() >= 0 {
					name := actionAttrs["name"].(types.String).ValueString()
					sort_order[name] = posn.ValueInt64()
				}
			}

			if isAction(key) {
				actionMapping(ActionMappers[key].ToDeploymentAction)
			}

			switch key {
			//case schemas.DeploymentProcessApplyTerraformTemplateAction:
			//	actionMapping(mapTerraformTemplateAction)
			//	break
			}

		}

		// Now that we have extracted all the steps off each of the properties into a single array, sort the array by the sort_order if provided
		if len(sort_order) > 0 {
			sort_order_entries := make(map[int64][]string)
			// Validate there are no duplicate sort_order entries
			for step_name, sort_order := range sort_order {
				sort_order_entries[sort_order] = append(sort_order_entries[sort_order], step_name)
			}
			for _, matching_names := range sort_order_entries {
				if len(matching_names) > 1 {
					tflog.Warn(ctx, fmt.Sprintf("The following actions have the same sort_order: %v", matching_names))
				}
			}

			// Validate that every step has a sort_order
			if len(sort_order) != len(step.Actions) {
				tflog.Warn(ctx, fmt.Sprintf("Not all actions on step '%s' have a `sort_order` parameter so they may be sorted in an unexpected order", step.Name))
			}

			sort.SliceStable(step.Actions, func(i, j int) bool {
				return sort_order[step.Actions[i].Name] < sort_order[step.Actions[j].Name]
			})
		}

		current.Steps = append(current.Steps, step)
	}
}

func isAction(key string) bool {
	for k, _ := range schemas.ActionsAttributeToActionTypeMap {
		if k == key {
			return true
		}
	}

	return false
}
