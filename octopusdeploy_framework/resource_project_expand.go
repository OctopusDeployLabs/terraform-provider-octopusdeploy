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

	if !model.GitLibraryPersistenceSettings.IsNull() {
		var gitLibrarySettingsList []gitLibraryPersistenceSettingsModel
		diags := model.GitLibraryPersistenceSettings.ElementsAs(ctx, &gitLibrarySettingsList, false)
		if diags.HasError() {
			fmt.Printf("Error converting Git library persistence settings: %v\n", diags)
		} else {
			fmt.Printf("Number of Git library persistence settings: %d\n", len(gitLibrarySettingsList))
			if len(gitLibrarySettingsList) > 0 {
				project.PersistenceSettings = expandGitLibraryPersistenceSettings(ctx, gitLibrarySettingsList[0])
			}
		}
	} else if !model.GitUsernamePasswordPersistenceSettings.IsNull() {
		var gitUsernamePasswordSettingsList []gitUsernamePasswordPersistenceSettingsModel
		diags := model.GitUsernamePasswordPersistenceSettings.ElementsAs(ctx, &gitUsernamePasswordSettingsList, false)
		if diags.HasError() {
			fmt.Printf("Error converting Git username/password persistence settings: %v\n", diags)
		} else {
			fmt.Printf("Number of Git username/password persistence settings: %d\n", len(gitUsernamePasswordSettingsList))
			if len(gitUsernamePasswordSettingsList) > 0 {
				project.PersistenceSettings = expandGitUsernamePasswordPersistenceSettings(ctx, gitUsernamePasswordSettingsList[0])
			}
		}
	} else if !model.GitAnonymousPersistenceSettings.IsNull() {
		var gitAnonymousSettingsList []gitAnonymousPersistenceSettingsModel
		diags := model.GitAnonymousPersistenceSettings.ElementsAs(ctx, &gitAnonymousSettingsList, false)
		if diags.HasError() {
			fmt.Printf("Error converting Git anonymous persistence settings: %v\n", diags)
		} else {
			fmt.Printf("Number of Git anonymous persistence settings: %d\n", len(gitAnonymousSettingsList))
			if len(gitAnonymousSettingsList) > 0 {
				project.PersistenceSettings = expandGitAnonymousPersistenceSettings(ctx, gitAnonymousSettingsList[0])
			}
		}
	}

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
		diags := model.Template.ElementsAs(ctx, &templates, false)
		if diags.HasError() {
			fmt.Printf("Error converting templates: %v\n", diags)
		} else {
			fmt.Printf("Number of templates: %d\n", len(templates))
			project.Templates = expandTemplates(templates)
		}
	} else {
		fmt.Println("Template is null")
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
	url, _ := url.Parse(model.URL.ValueString())
	var protectedBranches []string
	model.ProtectedBranches.ElementsAs(ctx, &protectedBranches, false)

	return projects.NewGitPersistenceSettings(
		model.BasePath.ValueString(),
		&credentials.Reference{
			ID: model.GitCredentialID.ValueString(),
		},
		model.DefaultBranch.ValueString(),
		protectedBranches,
		url,
	)
}

func expandGitUsernamePasswordPersistenceSettings(ctx context.Context, model gitUsernamePasswordPersistenceSettingsModel) projects.GitPersistenceSettings {
	url, _ := url.Parse(model.URL.ValueString())
	var protectedBranches []string
	model.ProtectedBranches.ElementsAs(ctx, &protectedBranches, false)

	return projects.NewGitPersistenceSettings(
		model.BasePath.ValueString(),
		&credentials.UsernamePassword{
			Username: model.Username.ValueString(),
			Password: core.NewSensitiveValue(model.Password.ValueString()),
		},
		model.DefaultBranch.ValueString(),
		protectedBranches,
		url,
	)
}

func expandGitAnonymousPersistenceSettings(ctx context.Context, model gitAnonymousPersistenceSettingsModel) projects.GitPersistenceSettings {
	url, _ := url.Parse(model.URL.ValueString())
	var protectedBranches []string
	model.ProtectedBranches.ElementsAs(ctx, &protectedBranches, false)

	return projects.NewGitPersistenceSettings(
		model.BasePath.ValueString(),
		&credentials.Anonymous{},
		model.DefaultBranch.ValueString(),
		protectedBranches,
		url,
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
