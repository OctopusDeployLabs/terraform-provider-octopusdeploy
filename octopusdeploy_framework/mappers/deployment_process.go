package mappers

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/gitdependencies"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"sort"
	"strconv"
	"strings"
)

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

func mapStepsToState(ctx context.Context, state *schemas.DeploymentProcessResourceModel, process *deployments.DeploymentProcess) diag.Diagnostics {
	if process.Steps == nil || len(process.Steps) == 0 {
		return nil
	}

	var steps []attr.Value
	for _, deploymentStep := range process.Steps {
		properties, diags := mapPropertiesToState(ctx, deploymentStep.Properties)
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

			newAction := map[string]attr.Value{
				"computed_sort_order": types.Int64Value(int64(i)),
			}

			sortOrder := getSortOrderStateValue(deploymentStep, a, state)
			if sortOrder != nil {
				newAction["sort_order"] = sortOrder
			}

			srcAction := deploymentStep.Actions[i]
			switch srcAction.ActionType {
			//case "Octopus.KubernetesDeploySecret":
			//	flatten_action_func("deploy_kubernetes_secret_action", i, flattenDeployKubernetesSecretAction)
			//case "Octopus.KubernetesRunScript":
			//	flatten_action_func("run_kubectl_script_action", i, flattenKubernetesRunScriptAction)
			//case "Octopus.Manual":
			//	flatten_action_func("manual_intervention_action", i, flattenManualInterventionAction)
			case "Octopus.Script":
				diag := mapRunScriptActionToState(ctx, a, newAction)
				if diag.HasError() {
					return diag
				}
				break
			//case "Octopus.TentaclePackage":
			//	flatten_action_func("deploy_package_action", i, flattenDeployPackageAction)
			//case "Octopus.TerraformApply":
			//	flatten_action_func("apply_terraform_template_action", i, flattenApplyTerraformTemplateAction)
			//case "Octopus.WindowsService":
			//	flatten_action_func("deploy_windows_service_action", i, flattenDeployWindowsServiceAction)
			default:
				diag := mapDeploymentActionToState(ctx, a, newAction)
				if diag.HasError() {
					return diag
				}
				break
			}

			terraformActionKeyName := getActionTypeTerraformKeyName(srcAction.ActionType)
			if _, ok := newActions[terraformActionKeyName]; !ok {
				newActions[terraformActionKeyName] = make([]attr.Value, 0)
			}
			newActions[terraformActionKeyName] = append(newActions[terraformActionKeyName], types.ObjectValueMust(getActionTypeAttrs(srcAction.ActionType), newAction))

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

func getSortOrderStateValue(step *deployments.DeploymentStep, a *deployments.DeploymentAction, state *schemas.DeploymentProcessResourceModel) attr.Value {
	for _, s := range state.Steps.Elements() {
		stepAttrs := s.(types.Object).Attributes()
		if step.Name == stepAttrs["name"].(types.String).ValueString() {
			for key, stepAttr := range stepAttrs {
				if isAction(key) {
					actions := stepAttr.(types.List)
					for _, action := range actions.Elements() {
						actionAttrs := action.(types.Object).Attributes()
						name := actionAttrs["name"].(types.String).ValueString()
						if a.Name == name {
							if v, ok := actionAttrs["sort_order"]; ok {
								return v
							} else {
								return nil
							}
						}
					}

				}
			}
		}
	}

	return nil
}

func isAction(key string) bool {
	for _, k := range schemas.ActionsAttributeToActionTypeMap {
		if k == key {
			return true
		}
	}

	return false
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

func getActionTypeTerraformKeyName(actionTypeName string) string {
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
		"git_dependency":                     types.SetType{ElemType: types.ObjectType{AttrTypes: getGitDependencyAttrTypes()}},
		"is_disabled":                        types.BoolType,
		"is_required":                        types.BoolType,
		"name":                               types.StringType,
		"notes":                              types.StringType,
		"primary_package":                    types.ListType{ElemType: types.ObjectType{AttrTypes: getPackageReferenceAttrTypes(true)}},
		"package":                            types.ListType{ElemType: types.ObjectType{AttrTypes: getPackageReferenceAttrTypes(false)}},
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
	case "Octopus.Script":
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
	default:
		attrs["action_type"] = types.StringType
		attrs["run_on_server"] = types.BoolType
		attrs["worker_pool_id"] = types.StringType
		attrs["worker_pool_variable"] = types.StringType
	}

	return attrs
}

func mapRunScriptActionToState(ctx context.Context, action *deployments.DeploymentAction, newAction map[string]attr.Value) diag.Diagnostics {
	diag := mapBaseDeploymentActionToState(ctx, action, newAction)
	if diag.HasError() {
		return diag
	}

	newAction["worker_pool_id"] = types.StringValue(action.WorkerPool)
	newAction["worker_pool_variable"] = types.StringValue(action.WorkerPoolVariable)

	mapPropertyToStateBool(action, newAction, "run_on_server", "Octopus.Action.RunOnServer", false)
	mapPropertyToStateString(action, newAction, "script_body", "Octopus.Action.Script.ScriptBody")
	mapPropertyToStateString(action, newAction, "script_file_name", "Octopus.Action.Script.ScriptFileName")
	mapPropertyToStateString(action, newAction, "script_source", "Octopus.Action.Script.ScriptSource")
	mapPropertyToStateString(action, newAction, "script_parameters", "Octopus.Action.Script.ScriptParameters")
	mapPropertyToStateString(action, newAction, "script_syntax", "Octopus.Action.Script.Syntax")
	mapPropertyToStateString(action, newAction, "script_file_name", "Octopus.Action.Script.ScriptFileName")
	mapPropertyToStateString(action, newAction, "variable_substitution_in_files", "Octopus.Action.SubstituteInFiles.TargetFiles")

	return nil
}

func mapPropertyToStateBool(action *deployments.DeploymentAction, actionState map[string]attr.Value, attrName string, propertyName string, defaultValue bool) {
	if v, ok := action.Properties[propertyName]; ok {
		parsedValue, _ := strconv.ParseBool(v.Value)
		actionState[attrName] = types.BoolValue(parsedValue)
	} else {
		actionState[attrName] = types.BoolValue(defaultValue)
	}
}

func mapPropertyToStateString(action *deployments.DeploymentAction, actionState map[string]attr.Value, attrName string, propertyName string) {
	if v, ok := action.Properties[propertyName]; ok {
		actionState[attrName] = types.StringValue(v.Value)
	} else {
		actionState[attrName] = types.StringValue("")
	}
}

func mapDeploymentActionToState(ctx context.Context, action *deployments.DeploymentAction, newAction map[string]attr.Value) diag.Diagnostics {
	diags := mapBaseDeploymentActionToState(ctx, action, newAction)
	if diags.HasError() {
		return diags
	}

	newAction["action_type"] = types.StringValue(action.ActionType)
	return nil
}

func mapBaseDeploymentActionToState(ctx context.Context, action *deployments.DeploymentAction, newAction map[string]attr.Value) diag.Diagnostics {
	newAction["can_be_used_for_project_versioning"] = types.BoolValue(action.CanBeUsedForProjectVersioning)
	newAction["is_disabled"] = types.BoolValue(action.IsDisabled)
	newAction["is_required"] = types.BoolValue(action.IsRequired)
	newAction["channels"] = util.FlattenStringList(action.Channels)
	newAction["condition"] = types.StringValue(action.Condition)
	newAction["container"] = mapContainerToState(action.Container)
	newAction["environments"] = util.FlattenStringList(action.Environments)
	newAction["excluded_environments"] = util.FlattenStringList(action.ExcludedEnvironments)
	newAction["id"] = types.StringValue(action.ID)
	newAction["name"] = types.StringValue(action.Name)
	newAction["slug"] = types.StringValue(action.Slug)
	newAction["notes"] = types.StringValue(action.Notes)

	updatedProperties, diags := mapPropertiesToState(ctx, action.Properties)
	if diags.HasError() {
		return diags
	}
	newAction["properties"] = updatedProperties
	//}

	//if len(action.TenantTags) > 0 {
	newAction["tenant_tags"] = util.FlattenStringList(action.TenantTags)
	//}

	if v, ok := action.Properties["Octopus.Action.EnabledFeatures"]; ok {
		newAction["features"] = util.FlattenStringList(strings.Split(v.Value, ","))
	} else {
		newAction["features"] = types.ListNull(types.StringType)
	}

	attrTypes := map[string]attr.Type{"id": types.StringType, "version": types.StringType}
	if v, ok := action.Properties["Octopus.Action.Template.Id"]; ok {
		actionTemplate := map[string]attr.Value{
			"id": types.StringValue(v.Value),
		}

		if v, ok := action.Properties["Octopus.Action.Template.Version"]; ok {
			actionTemplate["version"] = types.StringValue(v.Value)
		}

		list := make([]attr.Value, 1)
		list[0] = types.ObjectValueMust(attrTypes, actionTemplate)

		newAction["action_template"] = types.ListValueMust(types.ObjectType{AttrTypes: attrTypes}, list)
	} else {
		newAction["action_template"] = types.ListNull(types.ObjectType{AttrTypes: attrTypes})
	}

	hasPackageReference := false
	if len(action.Packages) > 0 {
		var packageReferences []attr.Value
		for _, packageReference := range action.Packages {
			packageReferenceAttribute, diags := mapPackageReferenceToState(ctx, packageReference)
			if diags.HasError() {
				return diags
			}
			if len(packageReference.Name) == 0 {

				newAction["primary_package"] = types.ListValueMust(types.ObjectType{AttrTypes: getPackageReferenceAttrTypes(true)}, []attr.Value{types.ObjectValueMust(getPackageReferenceAttrTypes(true), packageReferenceAttribute)})
				// TODO: consider these properties
				// actionProperties["Octopus.Action.Package.DownloadOnTentacle"] = packageReference.AcquisitionLocation
				// flattenedAction["properties"] = actionProperties
			} else {
				packageReferences = append(packageReferences, types.ObjectValueMust(getPackageReferenceAttrTypes(false), packageReferenceAttribute))
				newAction["package"] = types.ListValueMust(types.ObjectType{AttrTypes: getPackageReferenceAttrTypes(false)}, packageReferences)
				hasPackageReference = true
			}
		}
	} else {
		newAction["primary_package"] = types.ListNull(types.ObjectType{AttrTypes: getPackageReferenceAttrTypes(true)})
	}

	if !hasPackageReference {
		newAction["package"] = types.ListNull(types.ObjectType{AttrTypes: getPackageReferenceAttrTypes(false)})
	}

	if len(action.GitDependencies) > 0 {
		var gitDepenedencyList []attr.Value
		gitDepenedencyList = append(gitDepenedencyList, types.ObjectValueMust(getGitDependencyAttrTypes(), mapGitDependencyToState(action.GitDependencies[0])))
		newAction["git_dependency"] = types.SetValueMust(types.ObjectType{AttrTypes: getGitDependencyAttrTypes()}, gitDepenedencyList)
	} else {
		newAction["git_dependency"] = types.SetNull(types.ObjectType{AttrTypes: getGitDependencyAttrTypes()})
	}

	return nil
}

func mapPackageReferenceToState(ctx context.Context, packageReference *packages.PackageReference) (map[string]attr.Value, diag.Diagnostics) {
	properties, diags := types.MapValueFrom(ctx, types.StringType, packageReference.Properties)
	if diags.HasError() {
		return nil, diags
	}

	reference := map[string]attr.Value{
		"acquisition_location": types.StringValue(packageReference.AcquisitionLocation),
		"feed_id":              types.StringValue(packageReference.FeedID),
		"id":                   types.StringValue(packageReference.ID),
		"package_id":           types.StringValue(packageReference.PackageID),
		"properties":           properties,
	}

	if len(packageReference.Name) > 0 {
		if v, ok := packageReference.Properties["Extract"]; ok {
			extractDuringDeployment, _ := strconv.ParseBool(v)
			reference["extract_during_deployment"] = types.BoolValue(extractDuringDeployment)
		}
		reference["name"] = types.StringValue(packageReference.Name)
	}

	return reference, nil
}

func mapGitDependencyToState(gitDependency *gitdependencies.GitDependency) map[string]attr.Value {
	return map[string]attr.Value{
		"repository_uri":      types.StringValue(gitDependency.RepositoryUri),
		"default_branch":      types.StringValue(gitDependency.DefaultBranch),
		"git_credential_type": types.StringValue(gitDependency.GitCredentialType),
		"file_path_filters":   util.FlattenStringList(gitDependency.FilePathFilters),
		"git_credential_id":   types.StringValue(gitDependency.GitCredentialId),
	}
}

func getPackageReferenceAttrTypes(isPrimaryPackage bool) map[string]attr.Type {
	attrTypes := map[string]attr.Type{
		"acquisition_location": types.StringType,
		"feed_id":              types.StringType,
		"id":                   types.StringType,

		"package_id": types.StringType,

		"properties": types.MapType{types.StringType},
	}

	if !isPrimaryPackage {
		attrTypes["name"] = types.StringType
		attrTypes["extract_during_deployment"] = types.BoolType
	}

	return attrTypes
}

func getGitDependencyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"repository_uri":      types.StringType,
		"default_branch":      types.StringType,
		"git_credential_type": types.StringType,
		"file_path_filters":   types.ListType{ElemType: types.StringType},
		"git_credential_id":   types.StringType,
	}
}

func mapContainerToState(container *deployments.DeploymentActionContainer) types.List {
	attributeTypes := map[string]attr.Type{"feed_id": types.StringType, "image": types.StringType}
	if container == nil || (container.Image == "" && container.FeedID == "") {
		return types.ListNull(types.ObjectType{AttrTypes: attributeTypes})
	}

	list := make([]attr.Value, 0)
	containerAttributes := map[string]attr.Value{
		"feed_id": types.StringValue(container.FeedID),
		"image":   types.StringValue(container.Image),
	}

	list = append(list, types.ObjectValueMust(attributeTypes, containerAttributes))
	return types.ListValueMust(types.ObjectType{AttrTypes: attributeTypes}, list)
}

func mapPropertiesToState(ctx context.Context, properties map[string]core.PropertyValue) (types.Map, diag.Diagnostics) {

	if properties == nil || len(properties) == 0 {
		return types.MapNull(types.StringType), nil
	}

	stateMap := make(map[string]attr.Value)
	for key, value := range properties {
		if !value.IsSensitive {
			stateMap[key] = types.StringValue(value.Value)
		}
	}

	return types.MapValueFrom(ctx, types.StringType, stateMap)
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
				step.Actions = append(step.Actions, mappingFunc(attributes))
				actionAttrs := getActionAttributes(attributes)
				if posn, ok := actionAttrs["sort_order"].(types.Int64); ok && !posn.IsNull() && posn.ValueInt64() >= 0 {
					name := actionAttrs["name"].(types.String).ValueString()
					sort_order[name] = posn.ValueInt64()
				}
			}

			switch key {
			case schemas.DeploymentProcessAction:
				actionMapping(getBaseAction)
				break
			case schemas.DeploymentProcessRunScriptAction:
				actionMapping(mapRunScriptAction)
				break
			//case schemas.DeploymentProcessPackageAction:
			//	actionMapping(mapPackageAction)
			//	break
			case schemas.DeploymentProcessRunKubectlScriptAction:
				actionMapping(mapRunKubectlScriptAction)
				break
				//case schemas.DeploymentProcessApplyTerraformTemplateAction:
				//	actionMapping(mapTerraformTemplateAction)
				//	break
				//case schemas.DeploymentProcessApplyKubernetesSecretAction:
				//	actionMapping(mapKubernetesSecretAction)
				//	break
				//case schemas.DeploymentProcessWindowsServiceAction:
				//	actionMapping(mapWindowsServiceAction)
				//	break
				//case schemas.DeploymentProcessManualInterventionAction:
				//	actionMapping(mapManualInterventionAction)
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

func mapRunScriptAction(actionAttribute attr.Value) *deployments.DeploymentAction {
	actionAttrs := getActionAttributes(actionAttribute)
	if actionAttrs == nil {
		return nil
	}

	action := getBaseAction(actionAttribute)
	if action == nil {
		return nil
	}

	action.ActionType = "Octopus.Script"

	mapAttributeToProperty(action, actionAttrs, "script_file_name", "Octopus.Action.Script.ScriptFileName")
	mapAttributeToProperty(action, actionAttrs, "script_body", "Octopus.Action.Script.ScriptBody")
	mapAttributeToProperty(action, actionAttrs, "script_parameters", "Octopus.Action.Script.ScriptParameters")
	mapAttributeToProperty(action, actionAttrs, "script_source", "Octopus.Action.Script.ScriptSource")
	mapAttributeToProperty(action, actionAttrs, "script_syntax", "Octopus.Action.Script.Syntax")

	if variableSubstitutionInFiles, ok := actionAttrs["variable_substitution_in_files"]; ok {
		action.Properties["Octopus.Action.SubstituteInFiles.TargetFiles"] = core.NewPropertyValue(variableSubstitutionInFiles.(types.String).ValueString(), false)
		action.Properties["Octopus.Action.SubstituteInFiles.Enabled"] = core.NewPropertyValue(formatBoolForActionProperty(true), false)

		const substituteInFilesFeature = "Octopus.Features.SubstituteInFiles"
		const enabledFeatures = "Octopus.Action.EnabledFeatures"
		if len(action.Properties[enabledFeatures].Value) == 0 {
			action.Properties[enabledFeatures] = core.NewPropertyValue(substituteInFilesFeature, false)
		} else {
			// fixing https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/issues/641
			currentFeatures := action.Properties[enabledFeatures].Value
			if !strings.Contains(currentFeatures, substituteInFilesFeature) {
				action.Properties[enabledFeatures] = core.NewPropertyValue(currentFeatures+","+substituteInFilesFeature, false)
			}
		}
	}

	return action
}

func mapRunKubectlScriptAction(actionAttribute attr.Value) *deployments.DeploymentAction {
	actionAttrs := getActionAttributes(actionAttribute)
	if actionAttrs == nil {
		return nil
	}

	action := mapRunScriptAction(actionAttribute)
	if action == nil {
		return nil
	}
	action.ActionType = "Octopus.KubernetesRunScript"
	mapAttributeToProperty(action, actionAttrs, "namespace", "Octopus.Action.KubernetesContainers.Namespace")
	return action
}

func getActionAttributes(actionAttribute attr.Value) map[string]attr.Value {
	actionAttrsList := actionAttribute.(types.List)
	if actionAttrsList.IsNull() {
		return nil
	}

	actionAttrsElements := actionAttrsList.Elements()
	if len(actionAttrsElements) == 0 {
		return nil
	}

	return actionAttrsElements[0].(types.Object).Attributes()
}

func mapAttributeToProperty(action *deployments.DeploymentAction, attrs map[string]attr.Value, attributeName string, propertyName string) {
	var value string
	util.SetString(attrs, attributeName, &value)
	if value != "" {
		action.Properties[propertyName] = core.NewPropertyValue(value, false)
	}
}

func getBaseAction(actionAttribute attr.Value) *deployments.DeploymentAction {
	actionAttrs := getActionAttributes(actionAttribute)
	if actionAttrs == nil {
		return nil
	}

	var name string
	util.SetString(actionAttrs, "name", &name)

	var actionType string
	util.SetString(actionAttrs, "action_type", &actionType)

	action := deployments.NewDeploymentAction(name, actionType)
	util.SetString(actionAttrs, "id", &action.ID)
	util.SetString(actionAttrs, "condition", &action.Condition)
	util.SetBool(actionAttrs, "is_disabled", &action.IsDisabled)
	util.SetBool(actionAttrs, "is_required", &action.IsRequired)
	util.SetString(actionAttrs, "notes", &action.Notes)
	action.Channels = util.ExpandStringList(actionAttrs["channels"].(types.List))

	action.Container = getContainer(actionAttrs)

	action.Environments = getArray(actionAttrs, "environments")
	action.ExcludedEnvironments = getArray(actionAttrs, "excluded_environments")

	features := getArray(actionAttrs, "features")
	if features != nil {
		action.Properties["Octopus.Action.EnabledFeatures"] = core.NewPropertyValue(strings.Join(features, ","), false)
	}

	if v, ok := actionAttrs["run_on_server"]; ok {
		runOnServer := v.(types.Bool).ValueBool()
		action.Properties["Octopus.Action.RunOnServer"] = core.NewPropertyValue(formatBoolForActionProperty(runOnServer), false)
	}

	util.SetString(actionAttrs, "slug", &action.Slug)

	tenantTags := getArray(actionAttrs, "tenant_tags")
	if tenantTags != nil {
		action.TenantTags = tenantTags
	}

	util.SetString(actionAttrs, "worker_pool_id", &action.WorkerPool)
	util.SetString(actionAttrs, "worker_pool_variable", &action.WorkerPoolVariable)

	setActionTemplate(actionAttrs, action)
	setPrimaryPackage(actionAttrs, action)

	for key, attr := range actionAttrs {
		if key == "package" {
			for _, p := range attr.(types.List).Elements() {
				pkg := getPackageReference(p.(types.Object).Attributes())
				action.Packages = append(action.Packages, pkg)
			}
		}

		if key == "git_dependency" && len(attr.(types.Set).Elements()) > 0 {
			for _, gd := range attr.(types.Set).Elements() {
				gitDependency := getGitDependency(gd.(types.Object).Attributes())
				action.GitDependencies = append(action.GitDependencies, gitDependency)
			}

		}
	}

	// Polyfill the Kubernetes Object status check to default to true if not specified for Kubernetes steps
	switch actionType {
	case "Octopus.KubernetesDeployContainers":
		fallthrough
	case "Octopus.KubernetesDeployRawYaml":
		fallthrough
	case "Octopus.KubernetesDeployService":
		fallthrough
	case "Octopus.KubernetesDeployIngress":
		fallthrough
	case "Octopus.KubernetesDeployConfigMap":
		fallthrough
	case "Octopus.Kustomize":
		if _, exists := action.Properties["Octopus.Action.Kubernetes.ResourceStatusCheck"]; !exists {
			action.Properties["Octopus.Action.Kubernetes.ResourceStatusCheck"] = core.NewPropertyValue(formatBoolForActionProperty(true), false)
		}
		break
	}

	return action
}

func formatBoolForActionProperty(b bool) string {
	return cases.Title(language.Und, cases.NoLower).String(strconv.FormatBool(b))
}

func setPrimaryPackage(attrs map[string]attr.Value, action *deployments.DeploymentAction) {
	primaryPackageAttributes := getAttributesForSingleElementList(attrs, "primary_package")
	if primaryPackageAttributes == nil {
		return
	}

	primaryPackageReference := getPackageReference(primaryPackageAttributes)
	switch primaryPackageReference.AcquisitionLocation {
	case "Server":
		action.Properties["Octopus.Action.Package.DownloadOnTentacle"] = core.NewPropertyValue("False", false)
	default:
		action.Properties["Octopus.Action.Package.DownloadOnTentacle"] = core.NewPropertyValue(primaryPackageReference.AcquisitionLocation, false)
	}

	if len(primaryPackageReference.PackageID) > 0 {
		action.Properties["Octopus.Action.Package.PackageId"] = core.NewPropertyValue(primaryPackageReference.PackageID, false)
	}

	if len(primaryPackageReference.FeedID) > 0 {
		action.Properties["Octopus.Action.Package.FeedId"] = core.NewPropertyValue(primaryPackageReference.FeedID, false)
	}

	action.Packages = append(action.Packages, primaryPackageReference)
}

func getPackageReference(attrs map[string]attr.Value) *packages.PackageReference {
	pkg := &packages.PackageReference{Properties: map[string]string{}}
	util.SetString(attrs, "acquisition_location", &pkg.AcquisitionLocation)
	util.SetString(attrs, "feed_id", &pkg.FeedID)
	util.SetString(attrs, "name", &pkg.Name)
	util.SetString(attrs, "package_id", &pkg.PackageID)

	var extractDuringDeployment bool
	util.SetBool(attrs, "extract_during_deployment", &extractDuringDeployment)
	pkg.Properties["Extract"] = formatBoolForActionProperty(extractDuringDeployment)

	if properties := attrs["properties"]; properties != nil {
		propertyMap := properties.(types.Map).Elements()
		for k, v := range propertyMap {
			pkg.Properties[k] = v.(types.String).ValueString()
		}
	}

	return pkg
}

func setActionTemplate(attrs map[string]attr.Value, action *deployments.DeploymentAction) {
	actionTemplate := getAttributesForSingleElementList(attrs, "action_template")
	if actionTemplate != nil {
		if id, ok := actionTemplate["id"]; ok {
			action.Properties["Octopus.Action.Template.Id"] = core.NewPropertyValue(id.(types.String).ValueString(), false)
		}

		if v, ok := actionTemplate["version"]; ok {
			action.Properties["Octopus.Action.Template.Version"] = core.NewPropertyValue(v.(types.String).ValueString(), false)
		}
	}
}

func getAttributesForSingleElementList(attrs map[string]attr.Value, s string) map[string]attr.Value {
	if a, ok := attrs[s]; ok {
		list := a.(types.List)
		if len(list.Elements()) > 0 {
			return list.Elements()[0].(types.Object).Attributes()
		}
	}

	return nil
}

func getArray(attrs map[string]attr.Value, s string) []string {
	if a, ok := attrs[s]; ok {
		list := a.(types.List)
		return util.GetStringSlice(list)
	}

	return nil
}

func getContainer(attrs map[string]attr.Value) *deployments.DeploymentActionContainer {
	if c, ok := attrs["container"]; ok {
		if c == nil || c.IsNull() || c.IsUnknown() {
			return nil
		}

		containerAttrs := c.(types.List).Elements()[0].(types.Object).Attributes()
		actionContainer := &deployments.DeploymentActionContainer{}
		util.SetString(containerAttrs, "feed_id", &actionContainer.FeedID)
		util.SetString(containerAttrs, "image", &actionContainer.Image)
		return actionContainer
	}

	return nil
}

func getGitDependency(gitAttrs map[string]attr.Value) *gitdependencies.GitDependency {
	gitDependency := &gitdependencies.GitDependency{}
	util.SetString(gitAttrs, "repository_uri", &gitDependency.RepositoryUri)
	util.SetString(gitAttrs, "default_branch", &gitDependency.DefaultBranch)
	util.SetString(gitAttrs, "git_credential_type", &gitDependency.GitCredentialType)
	util.SetString(gitAttrs, "git_credential_id", &gitDependency.GitCredentialId)
	gitDependency.FilePathFilters = getArray(gitAttrs, "file_path_filters")
	return gitDependency
}
