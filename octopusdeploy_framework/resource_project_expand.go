package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/url"
)

func expandProject(ctx context.Context, model projectResourceModel) *projects.Project {
	project := projects.NewProject(
		model.Name.ValueString(),
		model.LifecycleID.ValueString(),
		model.ProjectGroupID.ValueString(),
	)

	project.ID = model.ID.ValueString()
	project.SpaceID = model.SpaceID.ValueString()
	project.Description = model.Description.ValueString()
	project.IsDisabled = model.IsDisabled.ValueBool()
	project.AutoCreateRelease = model.AutoCreateRelease.ValueBool()
	project.DefaultGuidedFailureMode = model.DefaultGuidedFailureMode.ValueString()
	project.DefaultToSkipIfAlreadyInstalled = model.DefaultToSkipIfAlreadyInstalled.ValueBool()
	project.DeploymentChangesTemplate = model.DeploymentChangesTemplate.ValueString()
	project.DeploymentProcessID = model.DeploymentProcessID.ValueString()
	project.IsDiscreteChannelRelease = model.IsDiscreteChannelRelease.ValueBool()
	project.IsVersionControlled = model.IsVersionControlled.ValueBool()
	project.TenantedDeploymentMode = core.TenantedDeploymentMode(model.TenantedDeploymentParticipation.ValueString())
	project.ReleaseNotesTemplate = model.ReleaseNotesTemplate.ValueString()
	project.Slug = model.Slug.ValueString()
	project.ClonedFromProjectID = model.ClonedFromProjectID.ValueString()

	if !model.IncludedLibraryVariableSets.IsNull() {
		var includedSets []string
		model.IncludedLibraryVariableSets.ElementsAs(ctx, &includedSets, false)
		project.IncludedLibraryVariableSets = includedSets
	}

	if !model.ConnectivityPolicy.IsNull() {
		project.ConnectivityPolicy = expandConnectivityPolicy(ctx, model.ConnectivityPolicy)
	}

	if !model.GitLibraryPersistenceSettings.IsNull() {
		var gitLibrarySettingsList []gitLibraryPersistenceSettingsModel
		diags := model.GitLibraryPersistenceSettings.ElementsAs(ctx, &gitLibrarySettingsList, false)
		if diags.HasError() {
			tflog.Error(ctx, fmt.Sprintf("Error converting Git library persistence settings: %v\n", diags))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Number of Git library persistence settings: %d\n", len(gitLibrarySettingsList)))
			if len(gitLibrarySettingsList) > 0 {
				project.PersistenceSettings = expandGitLibraryPersistenceSettings(ctx, gitLibrarySettingsList[0])
				project.IsVersionControlled = true
			}
		}
	} else if !model.GitUsernamePasswordPersistenceSettings.IsNull() {
		var gitUsernamePasswordSettingsList []gitUsernamePasswordPersistenceSettingsModel
		diags := model.GitUsernamePasswordPersistenceSettings.ElementsAs(ctx, &gitUsernamePasswordSettingsList, false)
		if diags.HasError() {
			tflog.Error(ctx, fmt.Sprintf("Error converting Git username/password persistence settings: %v\n", diags))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Number of Git username/password persistence settings: %d\n", len(gitUsernamePasswordSettingsList)))
			if len(gitUsernamePasswordSettingsList) > 0 {
				project.PersistenceSettings = expandGitUsernamePasswordPersistenceSettings(ctx, gitUsernamePasswordSettingsList[0])
				project.IsVersionControlled = true
			}
		}
	} else if !model.GitAnonymousPersistenceSettings.IsNull() {
		var gitAnonymousSettingsList []gitAnonymousPersistenceSettingsModel
		diags := model.GitAnonymousPersistenceSettings.ElementsAs(ctx, &gitAnonymousSettingsList, false)
		if diags.HasError() {
			tflog.Error(ctx, fmt.Sprintf("Error converting Git anonymous persistence settings: %v\n", diags))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Number of Git anonymous persistence settings: %d\n", len(gitAnonymousSettingsList)))
			if len(gitAnonymousSettingsList) > 0 {
				project.PersistenceSettings = expandGitAnonymousPersistenceSettings(ctx, gitAnonymousSettingsList[0])
				project.IsVersionControlled = true
			}
		}
	}

	if !model.JiraServiceManagementExtensionSettings.IsNull() {
		var settingsList []jiraServiceManagementExtensionSettingsModel
		diags := model.JiraServiceManagementExtensionSettings.ElementsAs(ctx, &settingsList, false)
		if !diags.HasError() && len(settingsList) > 0 {
			settings := settingsList[0]
			project.ExtensionSettings = append(project.ExtensionSettings, expandJiraServiceManagementExtensionSettings(settings))
		}
	}

	if !model.ServiceNowExtensionSettings.IsNull() {
		var settingsList []servicenowExtensionSettingsModel
		diags := model.ServiceNowExtensionSettings.ElementsAs(ctx, &settingsList, false)
		if !diags.HasError() && len(settingsList) > 0 {
			settings := settingsList[0]
			project.ExtensionSettings = append(project.ExtensionSettings, expandServiceNowExtensionSettings(settings))
		}
	}

	if !model.VersioningStrategy.IsNull() {
		project.VersioningStrategy = expandVersioningStrategy(ctx, model.VersioningStrategy)
	}

	if !model.ReleaseCreationStrategy.IsNull() {
		var strategy releaseCreationStrategyModel
		model.ReleaseCreationStrategy.ElementsAs(ctx, &strategy, false)
		project.ReleaseCreationStrategy = expandReleaseCreationStrategy(strategy)
	}

	if !model.Template.IsNull() {
		var templates []templateModel
		diags := model.Template.ElementsAs(ctx, &templates, false)
		if diags.HasError() {
			tflog.Error(ctx, fmt.Sprintf("Error converting templates: %v\n", diags))
		} else {
			tflog.Info(ctx, fmt.Sprintf("Number of templates: %d\n", len(templates)))
			project.Templates = expandTemplates(templates)
		}
	} else {
		tflog.Debug(ctx, "Template is null")
		project.Templates = []actiontemplates.ActionTemplateParameter{}
	}

	if !model.AutoDeployReleaseOverrides.IsNull() {
		var overrideModels []autoDeployReleaseOverrideModel
		diags := model.AutoDeployReleaseOverrides.ElementsAs(ctx, &overrideModels, false)
		if !diags.HasError() {
			project.AutoDeployReleaseOverrides = expandAutoDeployReleaseOverrides(overrideModels)
		}
	}

	return project
}

func expandGitLibraryPersistenceSettings(ctx context.Context, model gitLibraryPersistenceSettingsModel) projects.GitPersistenceSettings {
	gitUrl, _ := url.Parse(model.URL.ValueString())
	var protectedBranches []string
	model.ProtectedBranches.ElementsAs(ctx, &protectedBranches, false)

	return projects.NewGitPersistenceSettings(
		model.BasePath.ValueString(),
		&credentials.Reference{
			ID: model.GitCredentialID.ValueString(),
		},
		model.DefaultBranch.ValueString(),
		protectedBranches,
		gitUrl,
	)
}

func expandGitUsernamePasswordPersistenceSettings(ctx context.Context, model gitUsernamePasswordPersistenceSettingsModel) projects.GitPersistenceSettings {
	gitUrl, _ := url.Parse(model.URL.ValueString())
	var protectedBranches []string
	model.ProtectedBranches.ElementsAs(ctx, &protectedBranches, false)

	usernamePasswordCredential := credentials.NewUsernamePassword(
		model.Username.ValueString(),
		core.NewSensitiveValue(model.Password.ValueString()),
	)

	return projects.NewGitPersistenceSettings(
		model.BasePath.ValueString(),
		usernamePasswordCredential,
		model.DefaultBranch.ValueString(),
		protectedBranches,
		gitUrl,
	)
}

func expandGitAnonymousPersistenceSettings(ctx context.Context, model gitAnonymousPersistenceSettingsModel) projects.GitPersistenceSettings {
	gitUrl, _ := url.Parse(model.URL.ValueString())
	var protectedBranches []string
	model.ProtectedBranches.ElementsAs(ctx, &protectedBranches, false)

	return projects.NewGitPersistenceSettings(
		model.BasePath.ValueString(),
		&credentials.Anonymous{},
		model.DefaultBranch.ValueString(),
		protectedBranches,
		gitUrl,
	)
}

func expandAutoDeployReleaseOverrides(models []autoDeployReleaseOverrideModel) []projects.AutoDeployReleaseOverride {
	result := make([]projects.AutoDeployReleaseOverride, 0, len(models))

	for _, model := range models {
		override := projects.AutoDeployReleaseOverride{
			EnvironmentID: model.EnvironmentID.ValueString(),
		}

		if !model.TenantID.IsNull() {
			override.TenantID = model.TenantID.ValueString()
		}

		result = append(result, override)
	}

	return result
}

func expandConnectivityPolicy(ctx context.Context, connectivityPolicyList types.List) *core.ConnectivityPolicy {
	if connectivityPolicyList.IsNull() || connectivityPolicyList.IsUnknown() {
		return nil
	}

	var policyList []connectivityPolicyModel
	diags := connectivityPolicyList.ElementsAs(ctx, &policyList, false)
	if diags.HasError() {
		return nil
	}

	if len(policyList) == 0 {
		return nil
	}
	policy := policyList[0]

	var targetRoles []string
	if !policy.TargetRoles.IsNull() && !policy.TargetRoles.IsUnknown() {
		policy.TargetRoles.ElementsAs(ctx, &targetRoles, false)
	}

	skipMachineBehavior := core.SkipMachineBehavior(policy.SkipMachineBehavior.ValueString())

	return &core.ConnectivityPolicy{
		AllowDeploymentsToNoTargets: policy.AllowDeploymentsToNoTargets.ValueBool(),
		ExcludeUnhealthyTargets:     policy.ExcludeUnhealthyTargets.ValueBool(),
		SkipMachineBehavior:         skipMachineBehavior,
		TargetRoles:                 targetRoles,
	}
}

func expandJiraServiceManagementExtensionSettings(model jiraServiceManagementExtensionSettingsModel) *projects.JiraServiceManagementExtensionSettings {
	return projects.NewJiraServiceManagementExtensionSettings(
		model.ConnectionID.ValueString(),
		model.IsEnabled.ValueBool(),
		model.ServiceDeskProjectName.ValueString(),
	)
}

func expandServiceNowExtensionSettings(model servicenowExtensionSettingsModel) *projects.ServiceNowExtensionSettings {
	return projects.NewServiceNowExtensionSettings(
		model.ConnectionID.ValueString(),
		model.IsEnabled.ValueBool(),
		model.StandardChangeTemplateName.ValueString(),
		model.IsStateAutomaticallyTransitioned.ValueBool(),
	)
}

func expandVersioningStrategy(ctx context.Context, versioningStrategyList types.List) *projects.VersioningStrategy {
	if versioningStrategyList.IsNull() || versioningStrategyList.IsUnknown() {
		return nil
	}

	var strategyList []versioningStrategyModel
	diags := versioningStrategyList.ElementsAs(ctx, &strategyList, false)
	if diags.HasError() {
		return nil
	}

	if len(strategyList) == 0 {
		return nil
	}
	strategy := strategyList[0]

	versioningStrategy := &projects.VersioningStrategy{
		Template: strategy.Template.ValueString(),
	}

	if !strategy.DonorPackageStepID.IsNull() {
		donorPackageStepID := strategy.DonorPackageStepID.ValueString()
		versioningStrategy.DonorPackageStepID = &donorPackageStepID
	}

	if !strategy.DonorPackage.IsNull() {
		var donorPackageList []deploymentActionPackageModel
		diags := strategy.DonorPackage.ElementsAs(ctx, &donorPackageList, false)
		if !diags.HasError() && len(donorPackageList) > 0 {
			donorPackage := donorPackageList[0]
			versioningStrategy.DonorPackage = &packages.DeploymentActionPackage{
				DeploymentAction: donorPackage.DeploymentAction.ValueString(),
				PackageReference: donorPackage.PackageReference.ValueString(),
			}
		}
	}
	return versioningStrategy
}

func expandReleaseCreationStrategy(model releaseCreationStrategyModel) *projects.ReleaseCreationStrategy {
	strategy := &projects.ReleaseCreationStrategy{
		ChannelID:                    model.ChannelID.ValueString(),
		ReleaseCreationPackageStepID: model.ReleaseCreationPackageStepID.ValueString(),
	}
	if !model.ReleaseCreationPackage.IsNull() {
		var releaseCreationPackage deploymentActionPackageModel
		model.ReleaseCreationPackage.As(context.Background(), &releaseCreationPackage, basetypes.ObjectAsOptions{})
		strategy.ReleaseCreationPackage = expandDeploymentActionPackage(releaseCreationPackage)
	}
	return strategy
}

func expandDeploymentActionPackage(model deploymentActionPackageModel) *packages.DeploymentActionPackage {
	return &packages.DeploymentActionPackage{
		DeploymentAction: model.DeploymentAction.ValueString(),
		PackageReference: model.PackageReference.ValueString(),
	}
}
func expandTemplates(templates []templateModel) []actiontemplates.ActionTemplateParameter {
	result := make([]actiontemplates.ActionTemplateParameter, len(templates))
	for i, template := range templates {
		defaultValue := core.NewPropertyValue("", false)
		if !template.DefaultValue.IsNull() {
			defaultValue = core.NewPropertyValue(template.DefaultValue.ValueString(), false)
		}

		displaySettings := make(map[string]string)
		if !template.DisplaySettings.IsNull() && !template.DisplaySettings.IsUnknown() {
			template.DisplaySettings.ElementsAs(context.Background(), &displaySettings, false)
		}

		result[i] = actiontemplates.ActionTemplateParameter{
			DefaultValue:    &defaultValue,
			DisplaySettings: displaySettings,
			HelpText:        template.HelpText.ValueString(),
			Label:           template.Label.ValueString(),
			Name:            template.Name.ValueString(),
		}

		if !template.ID.IsNull() {
			result[i].Resource.ID = template.ID.ValueString()
		}
	}
	return result
}
