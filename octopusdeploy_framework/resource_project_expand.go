package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
		var connectivityPolicy connectivityPolicyModel
		model.ConnectivityPolicy.ElementsAs(ctx, &connectivityPolicy, false)
		project.ConnectivityPolicy = expandConnectivityPolicy(connectivityPolicy)
	}

	// TODO: git_library_persistence_settings
	//if v, ok := d.GetOk("git_library_persistence_settings"); ok {
	//	project.PersistenceSettings = expandGitPersistenceSettings(ctx, v, expandLibraryGitCredential)
	//}
	//if v, ok := d.GetOk("git_username_password_persistence_settings"); ok {
	//	project.PersistenceSettings = expandGitPersistenceSettings(ctx, v, expandUsernamePasswordGitCredential)
	//}
	//if v, ok := d.GetOk("git_anonymous_persistence_settings"); ok {
	//	project.PersistenceSettings = expandGitPersistenceSettings(ctx, v, expandAnonymousGitCredential)
	//}
	//
	//if project.PersistenceSettings != nil {
	//	tflog.Info(ctx, fmt.Sprintf("expanded persistence settings {%v}", project.PersistenceSettings))
	//}

	if !model.JiraServiceManagementExtensionSettings.IsNull() {
		var settings jiraServiceManagementExtensionSettingsModel
		model.JiraServiceManagementExtensionSettings.ElementsAs(ctx, &settings, false)
		project.ExtensionSettings = append(project.ExtensionSettings, expandJiraServiceManagementExtensionSettings(settings))
	}

	if !model.ServiceNowExtensionSettings.IsNull() {
		var settings servicenowExtensionSettingsModel
		model.ServiceNowExtensionSettings.ElementsAs(ctx, &settings, false)
		project.ExtensionSettings = append(project.ExtensionSettings, expandServiceNowExtensionSettings(settings))
	}

	if !model.VersioningStrategy.IsNull() {
		var strategy versioningStrategyModel
		model.VersioningStrategy.ElementsAs(ctx, &strategy, false)
		project.VersioningStrategy = expandVersioningStrategy(strategy)
	}

	if !model.ReleaseCreationStrategy.IsNull() {
		var strategy releaseCreationStrategyModel
		model.ReleaseCreationStrategy.ElementsAs(ctx, &strategy, false)
		project.ReleaseCreationStrategy = expandReleaseCreationStrategy(strategy)
	}

	if !model.Template.IsNull() {
		var templates []templateModel
		model.Template.ElementsAs(ctx, &templates, false)
		project.Templates = expandTemplates(templates)
	}

	if !model.AutoDeployReleaseOverrides.IsNull() {
		var overrideModels []autoDeployReleaseOverrideModel
		diags := model.AutoDeployReleaseOverrides.ElementsAs(ctx, &overrideModels, false)
		if !diags.HasError() {
			project.AutoDeployReleaseOverrides = expandAutoDeployReleaseOverrides(ctx, overrideModels)
		}
	}

	return project
}

func expandAutoDeployReleaseOverrides(ctx context.Context, models []autoDeployReleaseOverrideModel) []projects.AutoDeployReleaseOverride {
	result := make([]projects.AutoDeployReleaseOverride, 0, len(models))

	for _, model := range models {
		override := projects.AutoDeployReleaseOverride{
			EnvironmentID: model.EnvironmentID.ValueString(),
		}

		// TenantID is optional, so we only set it if it's not null
		if !model.TenantID.IsNull() {
			override.TenantID = model.TenantID.ValueString()
		}

		result = append(result, override)
	}

	return result
}

func expandGitPersistenceSettings(model gitPersistenceSettingsModel) projects.PersistenceSettings {
	gitUrl, _ := url.Parse(model.URL.ValueString())

	basePath := model.BasePath.ValueString()
	defaultBranch := model.DefaultBranch.ValueString()

	var protectedBranches []string
	model.ProtectedBranches.ElementsAs(context.Background(), &protectedBranches, false)

	var gitCredential credentials.GitCredential

	if !model.GitCredentialID.IsNull() {
		// Library Git Credential
		gitCredential = credentials.NewReference(model.GitCredentialID.ValueString())
	} else if !model.Username.IsNull() && !model.Password.IsNull() {
		// Username and Password Git Credential
		gitCredential = credentials.NewUsernamePassword(model.Username.ValueString(), core.NewSensitiveValue(model.Password.ValueString()))
	} else {
		// Anonymous Git Credential
		gitCredential = credentials.NewAnonymous()
	}

	return projects.NewGitPersistenceSettings(
		basePath,
		gitCredential,
		defaultBranch,
		protectedBranches,
		gitUrl,
	)
}

func expandConnectivityPolicy(model connectivityPolicyModel) *core.ConnectivityPolicy {
	var targetRoles []string
	if !model.TargetRoles.IsNull() && !model.TargetRoles.IsUnknown() {
		for _, v := range model.TargetRoles.Elements() {
			if strVal, ok := v.(types.String); ok {
				targetRoles = append(targetRoles, strVal.ValueString())
			}
		}
	}

	return &core.ConnectivityPolicy{
		AllowDeploymentsToNoTargets: model.AllowDeploymentsToNoTargets.ValueBool(),
		ExcludeUnhealthyTargets:     model.ExcludeUnhealthyTargets.ValueBool(),
		SkipMachineBehavior:         core.SkipMachineBehavior(model.SkipMachineBehavior.ValueString()),
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

func expandVersioningStrategy(model versioningStrategyModel) *projects.VersioningStrategy {
	strategy := &projects.VersioningStrategy{
		Template: model.Template.ValueString(),
	}
	if !model.DonorPackageStepID.IsNull() {
		donorPackageStepID := model.DonorPackageStepID.ValueString()
		strategy.DonorPackageStepID = &donorPackageStepID
	}
	if !model.DonorPackage.IsNull() {
		var donorPackage deploymentActionPackageModel
		model.DonorPackage.As(context.Background(), &donorPackage, basetypes.ObjectAsOptions{})
		strategy.DonorPackage = expandDeploymentActionPackage(donorPackage)
	}
	return strategy
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

func expandTemplates(models []templateModel) []actiontemplates.ActionTemplateParameter {
	templates := make([]actiontemplates.ActionTemplateParameter, len(models))
	for i, model := range models {
		defaultValue := core.NewPropertyValue(model.DefaultValue.ValueString(), false)

		displaySettings := make(map[string]string)
		if !model.DisplaySettings.IsNull() && !model.DisplaySettings.IsUnknown() {
			for k, v := range model.DisplaySettings.Elements() {
				if strVal, ok := v.(types.String); ok {
					displaySettings[k] = strVal.ValueString()
				}
			}
		}

		templates[i] = actiontemplates.ActionTemplateParameter{
			Name:            model.Name.ValueString(),
			Label:           model.Label.ValueString(),
			HelpText:        model.HelpText.ValueString(),
			DefaultValue:    &defaultValue,
			DisplaySettings: displaySettings,
		}
	}
	return templates
}
