package octopusdeploy_framework

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/extensions"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func flattenProject(ctx context.Context, project *projects.Project, state *projectResourceModel) (*projectResourceModel, diag.Diagnostics) {
	if project == nil {
		return nil, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Error flattening project",
				"The project is nil",
			),
		}
	}

	model := &projectResourceModel{
		SpaceID:                         types.StringValue(project.SpaceID),
		Name:                            types.StringValue(project.Name),
		Description:                     types.StringValue(project.Description),
		LifecycleID:                     types.StringValue(project.LifecycleID),
		ProjectGroupID:                  types.StringValue(project.ProjectGroupID),
		IsDisabled:                      types.BoolValue(project.IsDisabled),
		AutoCreateRelease:               types.BoolValue(project.AutoCreateRelease),
		DefaultGuidedFailureMode:        types.StringValue(project.DefaultGuidedFailureMode),
		DefaultToSkipIfAlreadyInstalled: types.BoolValue(project.DefaultToSkipIfAlreadyInstalled),
		DeploymentChangesTemplate:       types.StringValue(project.DeploymentChangesTemplate),
		DeploymentProcessID:             types.StringValue(project.DeploymentProcessID),
		DiscreteChannelRelease:          types.BoolValue(project.IsDiscreteChannelRelease),
		IsDiscreteChannelRelease:        types.BoolValue(project.IsDiscreteChannelRelease),
		IsVersionControlled:             types.BoolValue(project.IsVersionControlled),
		TenantedDeploymentParticipation: types.StringValue(string(project.TenantedDeploymentMode)),
		VariableSetID:                   types.StringValue(project.VariableSetID),
		ReleaseNotesTemplate:            util.StringOrNull(project.ReleaseNotesTemplate),
		Slug:                            types.StringValue(project.Slug),
		ClonedFromProjectID:             util.StringOrNull(project.ClonedFromProjectID),
	}

	model.ID = types.StringValue(project.GetID())

	if len(project.IncludedLibraryVariableSets) > 0 {
		model.IncludedLibraryVariableSets = util.FlattenStringList(project.IncludedLibraryVariableSets)
	} else {
		model.IncludedLibraryVariableSets = types.ListValueMust(types.StringType, []attr.Value{})
	}

	model.AutoDeployReleaseOverrides = flattenAutoDeployReleaseOverrides(project.AutoDeployReleaseOverrides)

	if state.ConnectivityPolicy.IsNull() {
		model.ConnectivityPolicy = types.ListNull(types.ObjectType{AttrTypes: getConnectivityPolicyAttrTypes()})
	} else {
		model.ConnectivityPolicy = flattenConnectivityPolicy(project.ConnectivityPolicy)
	}

	if state.ReleaseCreationStrategy.IsNull() {
		model.ReleaseCreationStrategy = types.ListNull(types.ObjectType{AttrTypes: getReleaseCreationStrategyAttrTypes()})
	} else {
		model.ReleaseCreationStrategy = flattenReleaseCreationStrategy(project.ReleaseCreationStrategy)
	}

	if state.VersioningStrategy.IsNull() {
		model.VersioningStrategy = types.ListNull(types.ObjectType{AttrTypes: getVersioningStrategyAttrTypes()})
	} else {
		model.VersioningStrategy = flattenVersioningStrategy(project.VersioningStrategy)
	}

	model.Template = flattenTemplates(project.Templates)

	diags := processPersistenceSettings(ctx, project, model)

	if diags.HasError() {
		return model, diags
	}

	// Extension Settings
	model.JiraServiceManagementExtensionSettings = flattenJiraServiceManagementExtensionSettings(nil)
	model.ServiceNowExtensionSettings = flattenServiceNowExtensionSettings(nil)

	for _, extensionSetting := range project.ExtensionSettings {
		switch extensionSetting.ExtensionID() {
		case extensions.JiraServiceManagementExtensionID:
			if jsmSettings, ok := extensionSetting.(*projects.JiraServiceManagementExtensionSettings); ok {
				model.JiraServiceManagementExtensionSettings = flattenJiraServiceManagementExtensionSettings(jsmSettings)
			}
		case extensions.ServiceNowExtensionID:
			if snowSettings, ok := extensionSetting.(*projects.ServiceNowExtensionSettings); ok {
				model.ServiceNowExtensionSettings = flattenServiceNowExtensionSettings(snowSettings)
			}
		}
	}

	return model, diags
}

func processPersistenceSettings(ctx context.Context, project *projects.Project, model *projectResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if project.PersistenceSettings != nil {
		if project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
			gitSettings := project.PersistenceSettings.(projects.GitPersistenceSettings)
			gitCredentialType := gitSettings.Credential().Type()
			model.IsVersionControlled = types.BoolValue(true)
			switch gitCredentialType {
			case credentials.GitCredentialTypeReference:
				model.GitLibraryPersistenceSettings, diags = flattenGitPersistenceSettings(ctx, gitSettings)
				model.GitUsernamePasswordPersistenceSettings = types.ListNull(types.ObjectType{AttrTypes: getGitUsernamePasswordPersistenceSettingsAttrTypes()})
				model.GitAnonymousPersistenceSettings = types.ListNull(types.ObjectType{AttrTypes: getGitAnonymousPersistenceSettingsAttrTypes()})
			case credentials.GitCredentialTypeUsernamePassword:
				model.GitLibraryPersistenceSettings = types.ListNull(types.ObjectType{AttrTypes: getGitLibraryPersistenceSettingsAttrTypes()})
				model.GitUsernamePasswordPersistenceSettings, diags = flattenGitPersistenceSettings(ctx, gitSettings)
				model.GitAnonymousPersistenceSettings = types.ListNull(types.ObjectType{AttrTypes: getGitAnonymousPersistenceSettingsAttrTypes()})
			case credentials.GitCredentialTypeAnonymous:
				model.GitLibraryPersistenceSettings = types.ListNull(types.ObjectType{AttrTypes: getGitLibraryPersistenceSettingsAttrTypes()})
				model.GitUsernamePasswordPersistenceSettings = types.ListNull(types.ObjectType{AttrTypes: getGitUsernamePasswordPersistenceSettingsAttrTypes()})
				model.GitAnonymousPersistenceSettings, diags = flattenGitPersistenceSettings(ctx, gitSettings)
			}
		} else {
			model.GitLibraryPersistenceSettings = types.ListNull(types.ObjectType{AttrTypes: getGitLibraryPersistenceSettingsAttrTypes()})
			model.GitUsernamePasswordPersistenceSettings = types.ListNull(types.ObjectType{AttrTypes: getGitUsernamePasswordPersistenceSettingsAttrTypes()})
			model.GitAnonymousPersistenceSettings = types.ListNull(types.ObjectType{AttrTypes: getGitAnonymousPersistenceSettingsAttrTypes()})
		}
	} else {
		model.GitLibraryPersistenceSettings = types.ListNull(types.ObjectType{AttrTypes: getGitLibraryPersistenceSettingsAttrTypes()})
		model.GitUsernamePasswordPersistenceSettings = types.ListNull(types.ObjectType{AttrTypes: getGitUsernamePasswordPersistenceSettingsAttrTypes()})
		model.GitAnonymousPersistenceSettings = types.ListNull(types.ObjectType{AttrTypes: getGitAnonymousPersistenceSettingsAttrTypes()})
		model.IsVersionControlled = types.BoolValue(false)
	}
	return diags
}

func flattenConnectivityPolicy(policy *core.ConnectivityPolicy) types.List {
	if policy == nil {
		return types.ListNull(types.ObjectType{AttrTypes: getConnectivityPolicyAttrTypes()})
	}

	obj := types.ObjectValueMust(getConnectivityPolicyAttrTypes(), map[string]attr.Value{
		"allow_deployments_to_no_targets": types.BoolValue(policy.AllowDeploymentsToNoTargets),
		"exclude_unhealthy_targets":       types.BoolValue(policy.ExcludeUnhealthyTargets),
		"skip_machine_behavior":           types.StringValue(string(policy.SkipMachineBehavior)),
		"target_roles":                    util.FlattenStringList(policy.TargetRoles),
	})

	return types.ListValueMust(types.ObjectType{AttrTypes: getConnectivityPolicyAttrTypes()}, []attr.Value{obj})
}

func flattenVersioningStrategy(strategy *projects.VersioningStrategy) types.List {
	if strategy == nil {
		return types.ListNull(types.ObjectType{AttrTypes: getVersioningStrategyAttrTypes()})
	}
	obj := types.ObjectValueMust(getVersioningStrategyAttrTypes(), map[string]attr.Value{
		"donor_package":         flattenDeploymentActionPackage(strategy.DonorPackage),
		"donor_package_step_id": types.StringPointerValue(strategy.DonorPackageStepID),
		"template":              types.StringValue(strategy.Template),
	})

	return types.ListValueMust(types.ObjectType{AttrTypes: getVersioningStrategyAttrTypes()}, []attr.Value{obj})
}

func flattenGitPersistenceSettings(ctx context.Context, persistenceSettings projects.PersistenceSettings) (types.List, diag.Diagnostics) {
	if persistenceSettings == nil || persistenceSettings.Type() == projects.PersistenceSettingsTypeDatabase {
		return types.ListNull(types.ObjectType{AttrTypes: getGitAnonymousPersistenceSettingsAttrTypes()}), nil
	}

	gitPersistenceSettings := persistenceSettings.(projects.GitPersistenceSettings)

	baseAttrValues := map[string]attr.Value{
		"base_path":      types.StringValue(gitPersistenceSettings.BasePath()),
		"default_branch": types.StringValue(gitPersistenceSettings.DefaultBranch()),
		"url":            types.StringValue(gitPersistenceSettings.URL().String()),
	}

	protectedBranches, diags := types.SetValueFrom(ctx, types.StringType, gitPersistenceSettings.ProtectedBranchNamePatterns())
	if diags.HasError() {
		return types.ListNull(types.ObjectType{}), diags
	}
	baseAttrValues["protected_branches"] = protectedBranches

	var attrTypes map[string]attr.Type
	var attrValues map[string]attr.Value

	credential := gitPersistenceSettings.Credential()
	switch credential.Type() {
	case credentials.GitCredentialTypeReference:
		attrTypes = getGitLibraryPersistenceSettingsAttrTypes()
		attrValues = baseAttrValues
		attrValues["git_credential_id"] = types.StringValue(credential.(*credentials.Reference).ID)
	case credentials.GitCredentialTypeUsernamePassword:
		attrTypes = getGitUsernamePasswordPersistenceSettingsAttrTypes()
		attrValues = baseAttrValues
		attrValues["username"] = types.StringValue(credential.(*credentials.UsernamePassword).Username)
		attrValues["password"] = types.StringValue(*credential.(*credentials.UsernamePassword).Password.NewValue)
	case credentials.GitCredentialTypeAnonymous:
		attrTypes = getGitAnonymousPersistenceSettingsAttrTypes()
		attrValues = baseAttrValues
	default:
		return types.ListNull(types.ObjectType{}), diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Unsupported Git Credential Type",
				fmt.Sprintf("Git credential type %v is not supported", credential.Type()),
			),
		}
	}

	objValue, diags := types.ObjectValue(attrTypes, attrValues)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: attrTypes}), diags
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: attrTypes}, []attr.Value{objValue}), nil
}

func flattenJiraServiceManagementExtensionSettings(settings *projects.JiraServiceManagementExtensionSettings) types.List {
	if settings == nil {
		return types.ListValueMust(types.ObjectType{AttrTypes: getJSMExtensionSettingsAttrTypes()}, []attr.Value{})
	}

	obj := types.ObjectValueMust(getJSMExtensionSettingsAttrTypes(), map[string]attr.Value{
		"connection_id":             types.StringValue(settings.ConnectionID()),
		"is_enabled":                types.BoolValue(settings.IsChangeControlled()),
		"service_desk_project_name": types.StringValue(settings.ServiceDeskProjectName),
	})

	return types.ListValueMust(types.ObjectType{AttrTypes: getJSMExtensionSettingsAttrTypes()}, []attr.Value{obj})
}

func flattenServiceNowExtensionSettings(settings *projects.ServiceNowExtensionSettings) types.List {
	if settings == nil {
		return types.ListValueMust(types.ObjectType{AttrTypes: getServiceNowExtensionSettingsAttrTypes()}, []attr.Value{})
	}

	obj := types.ObjectValueMust(getServiceNowExtensionSettingsAttrTypes(), map[string]attr.Value{
		"connection_id":                       types.StringValue(settings.ConnectionID()),
		"is_enabled":                          types.BoolValue(settings.IsChangeControlled()),
		"is_state_automatically_transitioned": types.BoolValue(settings.IsStateAutomaticallyTransitioned),
		"standard_change_template_name":       util.StringOrNull(settings.StandardChangeTemplateName),
	})

	return types.ListValueMust(types.ObjectType{AttrTypes: getServiceNowExtensionSettingsAttrTypes()}, []attr.Value{obj})
}

func flattenTemplates(templates []actiontemplates.ActionTemplateParameter) types.List {
	if len(templates) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: getTemplateAttrTypes()})
	}

	templateList := make([]attr.Value, 0, len(templates))
	for _, template := range templates {

		obj := types.ObjectValueMust(getTemplateAttrTypes(), map[string]attr.Value{
			"id":            types.StringValue(template.Resource.ID),
			"name":          types.StringValue(template.Name),
			"label":         util.StringOrNull(template.Label),
			"help_text":     util.StringOrNull(template.HelpText),
			"default_value": util.StringOrNull(template.DefaultValue.Value),
			"display_settings": types.MapValueMust(
				types.StringType,
				convertMapStringToMapAttrValue(template.DisplaySettings),
			),
		})

		templateList = append(templateList, obj)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: getTemplateAttrTypes()}, templateList)
}

func flattenAutoDeployReleaseOverrides(overrides []projects.AutoDeployReleaseOverride) types.List {
	if len(overrides) == 0 {
		return types.ListValueMust(types.ObjectType{AttrTypes: getAutoDeployReleaseOverrideAttrTypes()}, []attr.Value{})
	}

	overrideList := make([]attr.Value, 0, len(overrides))
	for _, override := range overrides {
		obj := types.ObjectValueMust(getAutoDeployReleaseOverrideAttrTypes(), map[string]attr.Value{
			"environment_id": types.StringValue(override.EnvironmentID),
			"release_id":     types.StringValue(override.ReleaseID),
			"tenant_id":      types.StringValue(override.TenantID),
		})
		overrideList = append(overrideList, obj)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: getAutoDeployReleaseOverrideAttrTypes()}, overrideList)
}

func getAutoDeployReleaseOverrideAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"environment_id": types.StringType,
		"release_id":     types.StringType,
		"tenant_id":      types.StringType,
	}
}

func flattenReleaseCreationStrategy(strategy *projects.ReleaseCreationStrategy) types.List {
	if strategy == nil {
		return types.ListValueMust(types.ObjectType{AttrTypes: getReleaseCreationStrategyAttrTypes()}, []attr.Value{})
	}

	obj := types.ObjectValueMust(getReleaseCreationStrategyAttrTypes(), map[string]attr.Value{
		"channel_id":                       types.StringValue(strategy.ChannelID),
		"release_creation_package_step_id": types.StringValue(strategy.ReleaseCreationPackageStepID),
		"release_creation_package":         flattenDeploymentActionPackage(strategy.ReleaseCreationPackage),
	})

	return types.ListValueMust(types.ObjectType{AttrTypes: getReleaseCreationStrategyAttrTypes()}, []attr.Value{obj})
}

func convertMapStringToMapAttrValue(m map[string]string) map[string]attr.Value {
	result := make(map[string]attr.Value, len(m))
	for k, v := range m {
		result[k] = types.StringValue(v)
	}
	return result
}

func flattenDeploymentActionPackage(pkg *packages.DeploymentActionPackage) types.List {
	if pkg == nil {
		return types.ListNull(types.ObjectType{AttrTypes: getDonorPackageAttrTypes()})
	}

	obj := types.ObjectValueMust(
		getDonorPackageAttrTypes(),
		map[string]attr.Value{
			"deployment_action": types.StringValue(pkg.DeploymentAction),
			"package_reference": types.StringValue(pkg.PackageReference),
		},
	)

	return types.ListValueMust(types.ObjectType{AttrTypes: getDonorPackageAttrTypes()}, []attr.Value{obj})
}

func getVersioningStrategyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"donor_package":         types.ListType{ElemType: types.ObjectType{AttrTypes: getDonorPackageAttrTypes()}},
		"donor_package_step_id": types.StringType,
		"template":              types.StringType,
	}
}

func getDonorPackageAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"deployment_action": types.StringType,
		"package_reference": types.StringType,
	}
}

func getConnectivityPolicyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"allow_deployments_to_no_targets": types.BoolType,
		"exclude_unhealthy_targets":       types.BoolType,
		"skip_machine_behavior":           types.StringType,
		"target_roles":                    types.ListType{ElemType: types.StringType},
	}
}

func getReleaseCreationStrategyAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"channel_id":                       types.StringType,
		"release_creation_package_step_id": types.StringType,
		"release_creation_package": types.ListType{ElemType: types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"deployment_action": types.StringType,
				"package_reference": types.StringType,
			},
		}},
	}
}

func getGitLibraryPersistenceSettingsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"git_credential_id":  types.StringType,
		"url":                types.StringType,
		"base_path":          types.StringType,
		"default_branch":     types.StringType,
		"protected_branches": types.SetType{ElemType: types.StringType},
	}
}

func getGitAnonymousPersistenceSettingsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"url":                types.StringType,
		"base_path":          types.StringType,
		"default_branch":     types.StringType,
		"protected_branches": types.SetType{ElemType: types.StringType},
	}
}

func getGitUsernamePasswordPersistenceSettingsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"url":                types.StringType,
		"base_path":          types.StringType,
		"default_branch":     types.StringType,
		"protected_branches": types.SetType{ElemType: types.StringType},
		"username":           types.StringType,
		"password":           types.StringType,
	}
}

func getJSMExtensionSettingsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"connection_id":             types.StringType,
		"is_enabled":                types.BoolType,
		"service_desk_project_name": types.StringType,
	}
}
func getTemplateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.StringType,
		"name":             types.StringType,
		"label":            types.StringType,
		"help_text":        types.StringType,
		"default_value":    types.StringType,
		"display_settings": types.MapType{ElemType: types.StringType},
	}
}

func getServiceNowExtensionSettingsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"connection_id":                       types.StringType,
		"is_enabled":                          types.BoolType,
		"is_state_automatically_transitioned": types.BoolType,
		"standard_change_template_name":       types.StringType,
	}
}
