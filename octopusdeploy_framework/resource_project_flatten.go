package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
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
		ID:                              types.StringValue(project.GetID()),
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
		ReleaseNotesTemplate:            types.StringValue(project.ReleaseNotesTemplate),
		Slug:                            types.StringValue(project.Slug),
		ClonedFromProjectID:             util.StringOrNull(project.ClonedFromProjectID),
	}

	var diags diag.Diagnostics

	model.IncludedLibraryVariableSets = util.FlattenStringList(project.IncludedLibraryVariableSets)
	model.AutoDeployReleaseOverrides = flattenAutoDeployReleaseOverrides(ctx, project.AutoDeployReleaseOverrides)

	if state.ConnectivityPolicy.IsNull() {
		model.ConnectivityPolicy = types.ListNull(types.ObjectType{AttrTypes: getConnectivityPolicyAttrTypes()})
	} else {
		model.ConnectivityPolicy = flattenConnectivityPolicy(ctx, project.ConnectivityPolicy)
	}

	if state.ReleaseCreationStrategy.IsNull() {
		model.ReleaseCreationStrategy = types.ListNull(types.ObjectType{AttrTypes: getReleaseCreationStrategyAttrTypes()})
	} else {
		model.ReleaseCreationStrategy = flattenReleaseCreationStrategy(ctx, project.ReleaseCreationStrategy)
	}

	if state.VersioningStrategy.IsNull() {
		model.VersioningStrategy = types.ListNull(types.ObjectType{AttrTypes: getVersioningStrategyAttrTypes()})
	} else {
		model.VersioningStrategy = flattenVersioningStrategy(project.VersioningStrategy)
	}

	// Template
	templateList, d := flattenTemplates(ctx, project.Templates)
	diags.Append(d...)
	model.Template = templateList

	// TODO GitPersistenceSetting
	model.GitLibraryPersistenceSettings, d = flattenGitLibraryPersistenceSettings(ctx, project.PersistenceSettings)
	model.GitAnonymousPersistenceSettings, d = flattenGitAnonymousPersistenceSettings(ctx, project.PersistenceSettings)
	model.GitUsernamePasswordPersistenceSettings, d = flattenGitUsernamePasswordPersistenceSettings(ctx, project.PersistenceSettings)

	// Extension Settings
	model.JiraServiceManagementExtensionSettings = flattenJiraServiceManagementExtensionSettings(ctx, nil)
	model.ServiceNowExtensionSettings = flattenServiceNowExtensionSettings(ctx, nil)

	for _, extensionSetting := range project.ExtensionSettings {
		switch extensionSetting.ExtensionID() {
		case extensions.JiraServiceManagementExtensionID:
			if jsmSettings, ok := extensionSetting.(*projects.JiraServiceManagementExtensionSettings); ok {
				model.JiraServiceManagementExtensionSettings = flattenJiraServiceManagementExtensionSettings(ctx, jsmSettings)
			}
		case extensions.ServiceNowExtensionID:
			if snowSettings, ok := extensionSetting.(*projects.ServiceNowExtensionSettings); ok {
				model.ServiceNowExtensionSettings = flattenServiceNowExtensionSettings(ctx, snowSettings)
			}
		}
	}

	return model, diags
}

func flattenConnectivityPolicy(ctx context.Context, policy *core.ConnectivityPolicy) types.List {
	if policy == nil {
		return types.ListValueMust(types.ObjectType{AttrTypes: getConnectivityPolicyAttrTypes()}, []attr.Value{})
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
		return types.ListValueMust(types.ObjectType{AttrTypes: getVersioningStrategyAttrTypes()}, []attr.Value{})
	}

	obj := types.ObjectValueMust(getVersioningStrategyAttrTypes(), map[string]attr.Value{
		"donor_package":         flattenDeploymentActionPackage(strategy.DonorPackage),
		"donor_package_step_id": types.StringPointerValue(strategy.DonorPackageStepID),
		"template":              types.StringValue(strategy.Template),
	})

	return types.ListValueMust(types.ObjectType{AttrTypes: getVersioningStrategyAttrTypes()}, []attr.Value{obj})
}

func flattenGitLibraryPersistenceSettings(ctx context.Context, settings projects.PersistenceSettings) (types.List, diag.Diagnostics) {
	return types.ListNull(types.ObjectType{AttrTypes: getGitLibraryPersistenceSettingsAttrTypes()}), nil
}

func flattenGitAnonymousPersistenceSettings(ctx context.Context, settings projects.PersistenceSettings) (types.List, diag.Diagnostics) {
	return types.ListNull(types.ObjectType{AttrTypes: getGitAnonymousPersistenceSettingsAttrTypes()}), nil
}

func flattenGitUsernamePasswordPersistenceSettings(ctx context.Context, settings projects.PersistenceSettings) (types.List, diag.Diagnostics) {
	return types.ListNull(types.ObjectType{AttrTypes: getGitUsernamePasswordPersistenceSettingsAttrTypes()}), nil
}

func flattenJiraServiceManagementExtensionSettings(ctx context.Context, settings *projects.JiraServiceManagementExtensionSettings) types.List {
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

func flattenServiceNowExtensionSettings(ctx context.Context, settings *projects.ServiceNowExtensionSettings) types.List {
	if settings == nil {
		return types.ListValueMust(types.ObjectType{AttrTypes: getServiceNowExtensionSettingsAttrTypes()}, []attr.Value{})
	}

	obj := types.ObjectValueMust(getServiceNowExtensionSettingsAttrTypes(), map[string]attr.Value{
		"connection_id":                       types.StringValue(settings.ConnectionID()),
		"is_enabled":                          types.BoolValue(settings.IsChangeControlled()),
		"is_state_automatically_transitioned": types.BoolValue(settings.IsStateAutomaticallyTransitioned),
		"standard_change_template_name":       types.StringValue(settings.StandardChangeTemplateName),
	})

	return types.ListValueMust(types.ObjectType{AttrTypes: getServiceNowExtensionSettingsAttrTypes()}, []attr.Value{obj})
}

func flattenTemplates(ctx context.Context, templates []actiontemplates.ActionTemplateParameter) (types.List, diag.Diagnostics) {
	if len(templates) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: getTemplateAttrTypes()}), nil
	}

	templateList := make([]attr.Value, 0, len(templates))
	for _, template := range templates {
		obj, diags := types.ObjectValueFrom(ctx, getTemplateAttrTypes(), map[string]attr.Value{
			"id":            types.StringValue(template.ID),
			"name":          types.StringValue(template.Name),
			"label":         types.StringValue(template.Label),
			"help_text":     types.StringValue(template.HelpText),
			"default_value": types.StringValue(template.DefaultValue.Value),
			"display_settings": types.MapValueMust(
				types.StringType,
				convertMapStringToMapAttrValue(template.DisplaySettings),
			),
		})
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: getTemplateAttrTypes()}), diags
		}
		templateList = append(templateList, obj)
	}

	return types.ListValueMust(types.ObjectType{AttrTypes: getTemplateAttrTypes()}, templateList), nil
}

func flattenAutoDeployReleaseOverrides(ctx context.Context, overrides []projects.AutoDeployReleaseOverride) types.List {
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

func flattenReleaseCreationStrategy(ctx context.Context, strategy *projects.ReleaseCreationStrategy) types.List {
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
