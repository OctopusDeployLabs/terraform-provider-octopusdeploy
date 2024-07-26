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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strconv"
	"strings"
)

func MapDeploymentProcessToSchema(ctx context.Context, deploymentProcess *deployments.DeploymentProcess, state *schemas.DeploymentProcessResourceModel) error {
	state.ID = types.StringValue(deploymentProcess.ID)
	state.Branch = types.StringValue(deploymentProcess.Branch)
	state.ProjectID = types.StringValue(deploymentProcess.ProjectID)
	state.SpaceID = types.StringValue(deploymentProcess.SpaceID)
	state.Version = types.StringValue(fmt.Sprintf("%d", deploymentProcess.Version))
	state.LastSnapshotID = types.StringValue(deploymentProcess.LastSnapshotID)

	mapStepsToState(state, deploymentProcess)

	return nil
}

func mapStepsToState(ctx context.Context, state *schemas.DeploymentProcessResourceModel, process *deployments.DeploymentProcess) diag.Diagnostics {
	if process.Steps == nil || len(process.Steps) == 0 {
		return nil
	}

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

		for propertyName, propertyValue := range deploymentStep.Properties {
			switch propertyName {
			case "Octopus.Action.TargetRoles":
				newStep["target_roles"] = util.FlattenStringList(strings.Split(propertyValue.Value, ","))
			case "Octopus.Action.MaxParallelism":
				newStep["window_size"] = types.StringValue(propertyValue.Value)
			case "Octopus.Step.ConditionVariableExpression":
				newStep["condition_expression"] = types.StringValue(propertyValue.Value)
			}
		}

		for a := range deploymentStep.Actions {
			switch deploymentStep.Actions[a].ActionType {
			case "Octopus.KubernetesDeploySecret":
				flatten_action_func("deploy_kubernetes_secret_action", i, flattenDeployKubernetesSecretAction)
			case "Octopus.KubernetesRunScript":
				flatten_action_func("run_kubectl_script_action", i, flattenKubernetesRunScriptAction)
			case "Octopus.Manual":
				flatten_action_func("manual_intervention_action", i, flattenManualInterventionAction)
			case "Octopus.Script":
				flatten_action_func("run_script_action", i, flattenRunScriptAction)
			case "Octopus.TentaclePackage":
				flatten_action_func("deploy_package_action", i, flattenDeployPackageAction)
			case "Octopus.TerraformApply":
				flatten_action_func("apply_terraform_template_action", i, flattenApplyTerraformTemplateAction)
			case "Octopus.WindowsService":
				flatten_action_func("deploy_windows_service_action", i, flattenDeployWindowsServiceAction)
			default:
				flatten_action_func("action", i, flattenDeploymentAction)
			}
			}
		}
	}
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

func MapSchemaToDeploymentProcess(plan schemas.DeploymentProcessResourceModel, deploymentProcess *deployments.DeploymentProcess) error {
	deploymentProcess.Branch = plan.Branch.ValueString()
	deploymentProcess.SpaceID = plan.SpaceID.ValueString()
	deploymentProcess.LastSnapshotID = plan.LastSnapshotID.ValueString()
	version, err := strconv.Atoi(plan.Version.ValueString())
	if err != nil {
		return err
	}
	deploymentProcess.Version = int32(version)

	deploymentProcess.ProjectID = plan.ProjectID.ValueString()
	mapStepsToDeploymentProcess(plan.Steps, deploymentProcess)

	return nil
}

func mapStepsToDeploymentProcess(steps types.List, current *deployments.DeploymentProcess) {
	if steps.IsNull() || steps.IsUnknown() {
		return
	}

	for _, s := range steps.Elements() {
		attrs := s.(types.Object).Attributes()
		step := deployments.NewDeploymentStep(attrs["name"].(types.String).String())
		step.Name = attrs["name"].(types.String).String()
		step.ID = attrs["id"].(types.String).String()
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

		for key, attr := range attrs {
			switch key {
			case "action":
				step.Actions = append(step.Actions, getBaseAction(attr))
				break
			case "run_script_action":
				step.Actions = append(step.Actions, mapRunScriptAction(attr))
				break
			}
		}
		//if actionsAttribute, ok := attrs["actions"]; ok {
		//	actions := actionsAttribute.(types.List)
		//	for _, action := range actions.Elements() {
		//		actionAttribute := action.(types.Object).Attributes()
		//
		//	}
		//}

		current.Steps = append(current.Steps, step)
	}
}

func mapRunScriptAction(actionAttribute attr.Value) *deployments.DeploymentAction {
	actionAttrs := actionAttribute.(types.Object).Attributes()

	action := getBaseAction(actionAttribute)
	if action == nil {
		return nil
	}

	action.ActionType = "Octopus.Script"

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

func mapAttributeToProperty(action *deployments.DeploymentAction, attrs map[string]attr.Value, attributeName string, propertyName string) {
	var value string
	util.SetString(attrs, attributeName, &value)
	if value != "" {
		action.Properties[propertyName] = core.NewPropertyValue(value, false)
	}
}

func getBaseAction(actionAttribute attr.Value) *deployments.DeploymentAction {
	actionAttrs := actionAttribute.(types.Object).Attributes()
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
	templateListAttributes := getAttributesForSingleElementList(attrs, "template_list")
	if templateListAttributes == nil {
		if id, ok := templateListAttributes["id"]; ok {
			action.Properties["Octopus.Action.Template.Id"] = core.NewPropertyValue(id.(types.String).ValueString(), false)
		}

		if v, ok := templateListAttributes["version"]; ok {
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
