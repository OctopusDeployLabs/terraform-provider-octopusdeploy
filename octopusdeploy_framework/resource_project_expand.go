package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/credentials"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
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

	if !model.IncludedLibraryVariableSets.IsNull() {
		var includedSets []string
		model.IncludedLibraryVariableSets.ElementsAs(ctx, &includedSets, false)
		project.IncludedLibraryVariableSets = includedSets
	}

	if !model.ConnectivityPolicy.IsNull() {
		var connectivityPolicy connectivityPolicyModel
		model.ConnectivityPolicy.As(ctx, &connectivityPolicy, basetypes.ObjectAsOptions{})
		project.ConnectivityPolicy = expandConnectivityPolicy(connectivityPolicy)
	}

	if !model.GitAnonymousPersistenceSettings.IsNull() {
		var settings gitAnonymousPersistenceSettingsModel
		model.GitAnonymousPersistenceSettings.As(ctx, &settings, basetypes.ObjectAsOptions{})
		project.PersistenceSettings = expandGitAnonymousPersistenceSettings(settings)
	} else if !model.GitLibraryPersistenceSettings.IsNull() {
		var settings gitLibraryPersistenceSettingsModel
		model.GitLibraryPersistenceSettings.As(ctx, &settings, basetypes.ObjectAsOptions{})
		project.PersistenceSettings = expandGitLibraryPersistenceSettings(settings)
	} else if !model.GitUsernamePasswordPersistenceSettings.IsNull() {
		var settings gitUsernamePasswordPersistenceSettingsModel
		model.GitUsernamePasswordPersistenceSettings.As(ctx, &settings, basetypes.ObjectAsOptions{})
		project.PersistenceSettings = expandGitUsernamePasswordPersistenceSettings(settings)
	}

	if !model.JiraServiceManagementExtensionSettings.IsNull() {
		var settings jiraServiceManagementExtensionSettingsModel
		model.JiraServiceManagementExtensionSettings.As(ctx, &settings, basetypes.ObjectAsOptions{})
		project.ExtensionSettings = append(project.ExtensionSettings, expandJiraServiceManagementExtensionSettings(settings))
	}
	if !model.ServicenowExtensionSettings.IsNull() {
		var settings servicenowExtensionSettingsModel
		model.ServicenowExtensionSettings.As(ctx, &settings, basetypes.ObjectAsOptions{})
		project.ExtensionSettings = append(project.ExtensionSettings, expandServiceNowExtensionSettings(settings))
	}

	if !model.VersioningStrategy.IsNull() {
		var strategy versioningStrategyModel
		model.VersioningStrategy.As(ctx, &strategy, basetypes.ObjectAsOptions{})
		project.VersioningStrategy = expandVersioningStrategy(strategy)
	}

	if !model.ReleaseCreationStrategy.IsNull() {
		var strategy releaseCreationStrategyModel
		model.ReleaseCreationStrategy.As(ctx, &strategy, basetypes.ObjectAsOptions{})
		project.ReleaseCreationStrategy = expandReleaseCreationStrategy(strategy)
	}

	if !model.Template.IsNull() {
		var templates []templateModel
		model.Template.ElementsAs(ctx, &templates, false)
		project.Templates = expandTemplates(templates)
	}

	return project
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

func expandGitAnonymousPersistenceSettings(model gitAnonymousPersistenceSettingsModel) projects.GitPersistenceSettings {
	url, _ := url.Parse(model.URL.ValueString())
	return projects.NewGitPersistenceSettings(
		model.BasePath.ValueString(),
		credentials.NewAnonymous(),
		model.DefaultBranch.ValueString(),
		util.ExpandStringArray(model.ProtectedBranches),
		url,
	)
}

func expandGitLibraryPersistenceSettings(model gitLibraryPersistenceSettingsModel) projects.GitPersistenceSettings {
	url, _ := url.Parse(model.URL.ValueString())
	return projects.NewGitPersistenceSettings(
		model.BasePath.ValueString(),
		credentials.NewReference(model.GitCredentialID.ValueString()),
		model.DefaultBranch.ValueString(),
		util.ExpandStringArray(model.ProtectedBranches),
		url,
	)
}

func expandGitUsernamePasswordPersistenceSettings(model gitUsernamePasswordPersistenceSettingsModel) projects.GitPersistenceSettings {
	url, _ := url.Parse(model.URL.ValueString())
	passwordSensitiveValue := core.NewSensitiveValue(model.Password.ValueString())
	return projects.NewGitPersistenceSettings(
		model.BasePath.ValueString(),
		credentials.NewUsernamePassword(model.Username.ValueString(), passwordSensitiveValue),
		model.DefaultBranch.ValueString(),
		util.ExpandStringArray(model.ProtectedBranches),
		url,
	)
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
