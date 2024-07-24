package octopusdeploy_framework

import (
	"context"
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

func flattenProject(ctx context.Context, project *projects.Project) (*projectResourceModel, diag.Diagnostics) {
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
		IsDiscreteChannelRelease:        types.BoolValue(project.IsDiscreteChannelRelease),
		IsVersionControlled:             types.BoolValue(project.IsVersionControlled),
		TenantedDeploymentParticipation: types.StringValue(string(project.TenantedDeploymentMode)),
		VariableSetID:                   types.StringValue(project.VariableSetID),
		ReleaseNotesTemplate:            types.StringValue(project.ReleaseNotesTemplate),
		Slug:                            types.StringValue(project.Slug),
	}

	var diags diag.Diagnostics

	model.IncludedLibraryVariableSets, diags = types.ListValueFrom(ctx, types.StringType, project.IncludedLibraryVariableSets)
	if diags.HasError() {
		return nil, diags
	}

	model.ConnectivityPolicy, diags = flattenConnectivityPolicy(ctx, project.ConnectivityPolicy)
	if diags.HasError() {
		return nil, diags
	}

	if project.PersistenceSettings != nil {
		switch project.PersistenceSettings.Type() {
		case projects.PersistenceSettingsTypeVersionControlled:
			gitSettings := project.PersistenceSettings.(projects.GitPersistenceSettings)
			switch gitSettings.Credential().Type() {
			case credentials.GitCredentialTypeAnonymous:
				model.GitAnonymousPersistenceSettings, diags = flattenGitAnonymousPersistenceSettings(ctx, gitSettings)
			case credentials.GitCredentialTypeReference:
				model.GitLibraryPersistenceSettings, diags = flattenGitLibraryPersistenceSettings(ctx, gitSettings)
			case credentials.GitCredentialTypeUsernamePassword:
				model.GitUsernamePasswordPersistenceSettings, diags = flattenGitUsernamePasswordPersistenceSettings(ctx, gitSettings)
			}
			if diags.HasError() {
				return nil, diags
			}
		}
	}

	for _, extensionSetting := range project.ExtensionSettings {
		switch extensionSetting.ExtensionID() {
		case extensions.JiraServiceManagementExtensionID:
			if jsmSettings, ok := extensionSetting.(*projects.JiraServiceManagementExtensionSettings); ok {
				model.JiraServiceManagementExtensionSettings, diags = flattenJiraServiceManagementExtensionSettings(ctx, jsmSettings)
				if diags.HasError() {
					return nil, diags
				}
			}
		case extensions.ServiceNowExtensionID:
			if snowSettings, ok := extensionSetting.(*projects.ServiceNowExtensionSettings); ok {
				model.ServicenowExtensionSettings, diags = flattenServiceNowExtensionSettings(ctx, snowSettings)
				if diags.HasError() {
					return nil, diags
				}
			}
		}
	}

	if project.VersioningStrategy != nil {
		model.VersioningStrategy, diags = flattenVersioningStrategy(ctx, project.VersioningStrategy)
		if diags.HasError() {
			return nil, diags
		}
	}

	if project.ReleaseCreationStrategy != nil {
		model.ReleaseCreationStrategy, diags = flattenReleaseCreationStrategy(ctx, project.ReleaseCreationStrategy)
		if diags.HasError() {
			return nil, diags
		}
	}

	if len(project.Templates) > 0 {
		model.Template, diags = flattenTemplates(ctx, project.Templates)
		if diags.HasError() {
			return nil, diags
		}
	}

	return model, nil
}

func flattenConnectivityPolicy(ctx context.Context, policy *core.ConnectivityPolicy) (types.Object, diag.Diagnostics) {
	if policy == nil {
		return types.ObjectNull(map[string]attr.Type{
			"allow_deployments_to_no_targets": types.BoolType,
			"exclude_unhealthy_targets":       types.BoolType,
			"skip_machine_behavior":           types.StringType,
			"target_roles":                    types.ListType{ElemType: types.StringType},
		}), nil
	}

	return types.ObjectValueFrom(ctx, map[string]attr.Type{
		"allow_deployments_to_no_targets": types.BoolType,
		"exclude_unhealthy_targets":       types.BoolType,
		"skip_machine_behavior":           types.StringType,
		"target_roles":                    types.ListType{ElemType: types.StringType},
	}, map[string]attr.Value{
		"allow_deployments_to_no_targets": types.BoolValue(policy.AllowDeploymentsToNoTargets),
		"exclude_unhealthy_targets":       types.BoolValue(policy.ExcludeUnhealthyTargets),
		"skip_machine_behavior":           types.StringValue(string(policy.SkipMachineBehavior)),
		"target_roles":                    util.FlattenStringList(policy.TargetRoles),
	})
}

func flattenGitAnonymousPersistenceSettings(ctx context.Context, settings projects.GitPersistenceSettings) (types.Object, diag.Diagnostics) {
	protectedBranches, diags := util.TerraformSetFromStringArray(ctx, settings.ProtectedBranchNamePatterns())
	if diags.HasError() {
		return types.ObjectNull(nil), diags
	}

	return types.ObjectValueFrom(ctx, map[string]attr.Type{
		"url":                types.StringType,
		"base_path":          types.StringType,
		"default_branch":     types.StringType,
		"protected_branches": types.SetType{ElemType: types.StringType},
	}, map[string]attr.Value{
		"url":                types.StringValue(settings.URL().String()),
		"base_path":          types.StringValue(settings.BasePath()),
		"default_branch":     types.StringValue(settings.DefaultBranch()),
		"protected_branches": protectedBranches,
	})
}

func flattenGitLibraryPersistenceSettings(ctx context.Context, settings projects.GitPersistenceSettings) (types.Object, diag.Diagnostics) {
	credential := settings.Credential().(*credentials.Reference)
	protectedBranches, diags := util.TerraformSetFromStringArray(ctx, settings.ProtectedBranchNamePatterns())
	if diags.HasError() {
		return types.ObjectNull(nil), diags
	}
	return types.ObjectValueFrom(ctx, map[string]attr.Type{
		"url":                types.StringType,
		"base_path":          types.StringType,
		"default_branch":     types.StringType,
		"protected_branches": types.SetType{ElemType: types.StringType},
		"git_credential_id":  types.StringType,
	}, map[string]attr.Value{
		"url":                types.StringValue(settings.URL().String()),
		"base_path":          types.StringValue(settings.BasePath()),
		"default_branch":     types.StringValue(settings.DefaultBranch()),
		"protected_branches": protectedBranches,
		"git_credential_id":  types.StringValue(credential.ID),
	})
}

func flattenGitUsernamePasswordPersistenceSettings(ctx context.Context, settings projects.GitPersistenceSettings) (types.Object, diag.Diagnostics) {
	credential := settings.Credential().(*credentials.UsernamePassword)
	protectedBranches, diags := util.TerraformSetFromStringArray(ctx, settings.ProtectedBranchNamePatterns())
	if diags.HasError() {
		return types.ObjectNull(nil), diags
	}

	var passwordValue string
	if credential.Password != nil && credential.Password.NewValue != nil {
		passwordValue = *credential.Password.NewValue
	}

	return types.ObjectValueFrom(ctx, map[string]attr.Type{
		"url":                types.StringType,
		"base_path":          types.StringType,
		"default_branch":     types.StringType,
		"protected_branches": types.SetType{ElemType: types.StringType},
		"username":           types.StringType,
		"password":           types.StringType,
	}, map[string]attr.Value{
		"url":                types.StringValue(settings.URL().String()),
		"base_path":          types.StringValue(settings.BasePath()),
		"default_branch":     types.StringValue(settings.DefaultBranch()),
		"protected_branches": protectedBranches,
		"username":           types.StringValue(credential.Username),
		"password":           types.StringValue(passwordValue),
	})
}

func flattenJiraServiceManagementExtensionSettings(ctx context.Context, settings *projects.JiraServiceManagementExtensionSettings) (types.Object, diag.Diagnostics) {
	return types.ObjectValueFrom(ctx, map[string]attr.Type{
		"connection_id":             types.StringType,
		"is_enabled":                types.BoolType,
		"service_desk_project_name": types.StringType,
	}, map[string]attr.Value{
		"connection_id":             types.StringValue(settings.ConnectionID()),
		"is_enabled":                types.BoolValue(settings.IsChangeControlled()),
		"service_desk_project_name": types.StringValue(settings.ServiceDeskProjectName),
	})
}

func flattenServiceNowExtensionSettings(ctx context.Context, settings *projects.ServiceNowExtensionSettings) (types.Object, diag.Diagnostics) {
	return types.ObjectValueFrom(ctx, map[string]attr.Type{
		"connection_id":                       types.StringType,
		"is_enabled":                          types.BoolType,
		"is_state_automatically_transitioned": types.BoolType,
		"standard_change_template_name":       types.StringType,
	}, map[string]attr.Value{
		"connection_id":                       types.StringValue(settings.ConnectionID()),
		"is_enabled":                          types.BoolValue(settings.IsChangeControlled()),
		"is_state_automatically_transitioned": types.BoolValue(settings.IsStateAutomaticallyTransitioned),
		"standard_change_template_name":       types.StringValue(settings.StandardChangeTemplateName),
	})
}

func flattenDeploymentActionPackage(ctx context.Context, deploymentActionPackage *packages.DeploymentActionPackage) (types.Object, diag.Diagnostics) {
	if deploymentActionPackage == nil {
		return types.ObjectNull(map[string]attr.Type{
			"deployment_action": types.StringType,
			"package_reference": types.StringType,
		}), nil
	}

	return types.ObjectValueFrom(ctx, map[string]attr.Type{
		"deployment_action": types.StringType,
		"package_reference": types.StringType,
	}, map[string]attr.Value{
		"deployment_action": types.StringValue(deploymentActionPackage.DeploymentAction),
		"package_reference": types.StringValue(deploymentActionPackage.PackageReference),
	})
}

func flattenVersioningStrategy(ctx context.Context, versioningStrategy *projects.VersioningStrategy) (types.Object, diag.Diagnostics) {
	if versioningStrategy == nil {
		return types.ObjectNull(map[string]attr.Type{
			"donor_package": types.ObjectType{AttrTypes: map[string]attr.Type{
				"deployment_action": types.StringType,
				"package_reference": types.StringType,
			}},
			"donor_package_step_id": types.StringType,
			"template":              types.StringType,
		}), nil
	}

	donorPackage, diags := flattenDeploymentActionPackage(ctx, versioningStrategy.DonorPackage)
	if diags.HasError() {
		return types.ObjectNull(nil), diags
	}

	return types.ObjectValueFrom(ctx, map[string]attr.Type{
		"donor_package": types.ObjectType{AttrTypes: map[string]attr.Type{
			"deployment_action": types.StringType,
			"package_reference": types.StringType,
		}},
		"donor_package_step_id": types.StringType,
		"template":              types.StringType,
	}, map[string]attr.Value{
		"donor_package":         donorPackage,
		"donor_package_step_id": types.StringPointerValue(versioningStrategy.DonorPackageStepID),
		"template":              types.StringValue(versioningStrategy.Template),
	})
}
func flattenReleaseCreationStrategy(ctx context.Context, strategy *projects.ReleaseCreationStrategy) (types.Object, diag.Diagnostics) {
	if strategy == nil {
		return types.ObjectNull(map[string]attr.Type{
			"channel_id":                       types.StringType,
			"release_creation_package_step_id": types.StringType,
			"release_creation_package": types.ObjectType{AttrTypes: map[string]attr.Type{
				"deployment_action": types.StringType,
				"package_reference": types.StringType,
			}},
		}), nil
	}

	releaseCreationPackage, diags := flattenDeploymentActionPackage(ctx, strategy.ReleaseCreationPackage)
	if diags.HasError() {
		return types.ObjectNull(nil), diags
	}

	return types.ObjectValueFrom(ctx, map[string]attr.Type{
		"channel_id":                       types.StringType,
		"release_creation_package_step_id": types.StringType,
		"release_creation_package": types.ObjectType{AttrTypes: map[string]attr.Type{
			"deployment_action": types.StringType,
			"package_reference": types.StringType,
		}},
	}, map[string]attr.Value{
		"channel_id":                       types.StringValue(strategy.ChannelID),
		"release_creation_package_step_id": types.StringValue(strategy.ReleaseCreationPackageStepID),
		"release_creation_package":         releaseCreationPackage,
	})
}

func convertDisplaySettings(m map[string]string) map[string]attr.Value {
	result := make(map[string]attr.Value, len(m))
	for k, v := range m {
		result[k] = types.StringValue(v)
	}
	return result
}

func flattenTemplates(ctx context.Context, templates []actiontemplates.ActionTemplateParameter) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(templates) == 0 {
		return types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{
			"id":               types.StringType,
			"name":             types.StringType,
			"label":            types.StringType,
			"help_text":        types.StringType,
			"default_value":    types.StringType,
			"display_settings": types.MapType{ElemType: types.StringType},
		}}), diags
	}

	templateList := make([]attr.Value, 0, len(templates))
	for _, template := range templates {
		displaySettingsValue, mapDiags := types.MapValue(types.StringType, convertDisplaySettings(template.DisplaySettings))
		diags.Append(mapDiags...)
		if diags.HasError() {
			return types.ListNull(nil), diags
		}

		flattenedTemplate, objectDiags := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"id":               types.StringType,
			"name":             types.StringType,
			"label":            types.StringType,
			"help_text":        types.StringType,
			"default_value":    types.StringType,
			"display_settings": types.MapType{ElemType: types.StringType},
		}, map[string]attr.Value{
			"id":               types.StringValue(template.ID),
			"name":             types.StringValue(template.Name),
			"label":            types.StringValue(template.Label),
			"help_text":        types.StringValue(template.HelpText),
			"default_value":    types.StringValue(template.DefaultValue.Value),
			"display_settings": displaySettingsValue,
		})
		diags.Append(objectDiags...)
		if diags.HasError() {
			return types.ListNull(nil), diags
		}
		templateList = append(templateList, flattenedTemplate)
	}

	listValue, listDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":               types.StringType,
		"name":             types.StringType,
		"label":            types.StringType,
		"help_text":        types.StringType,
		"default_value":    types.StringType,
		"display_settings": types.MapType{ElemType: types.StringType},
	}}, templateList)
	diags.Append(listDiags...)

	return listValue, diags
}
