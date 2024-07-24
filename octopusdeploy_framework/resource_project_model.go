package octopusdeploy_framework

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type projectResourceModel struct {
	ID                                     types.String `tfsdk:"id"`
	SpaceID                                types.String `tfsdk:"space_id"`
	Name                                   types.String `tfsdk:"name"`
	Description                            types.String `tfsdk:"description"`
	LifecycleID                            types.String `tfsdk:"lifecycle_id"`
	ProjectGroupID                         types.String `tfsdk:"project_group_id"`
	IsDisabled                             types.Bool   `tfsdk:"is_disabled"`
	AutoCreateRelease                      types.Bool   `tfsdk:"auto_create_release"`
	DefaultGuidedFailureMode               types.String `tfsdk:"default_guided_failure_mode"`
	DefaultToSkipIfAlreadyInstalled        types.Bool   `tfsdk:"default_to_skip_if_already_installed"`
	DeploymentChangesTemplate              types.String `tfsdk:"deployment_changes_template"`
	DeploymentProcessID                    types.String `tfsdk:"deployment_process_id"`
	DiscreteChannelRelease                 types.Bool   `tfsdk:"discrete_channel_release"`
	IsDiscreteChannelRelease               types.Bool   `tfsdk:"is_discrete_channel_release"`
	IsVersionControlled                    types.Bool   `tfsdk:"is_version_controlled"`
	TenantedDeploymentParticipation        types.String `tfsdk:"tenanted_deployment_participation"`
	VariableSetID                          types.String `tfsdk:"variable_set_id"`
	ReleaseNotesTemplate                   types.String `tfsdk:"release_notes_template"`
	Slug                                   types.String `tfsdk:"slug"`
	ClonedFromProjectID                    types.String `tfsdk:"cloned_from_project_id"`
	VersioningStrategy                     types.List   `tfsdk:"versioning_strategy"`
	ConnectivityPolicy                     types.List   `tfsdk:"connectivity_policy"`
	ReleaseCreationStrategy                types.List   `tfsdk:"release_creation_strategy"`
	Template                               types.List   `tfsdk:"template"`
	GitAnonymousPersistenceSettings        types.List   `tfsdk:"git_anonymous_persistence_settings"`
	GitLibraryPersistenceSettings          types.List   `tfsdk:"git_library_persistence_settings"`
	GitUsernamePasswordPersistenceSettings types.List   `tfsdk:"git_username_password_persistence_settings"`
	JiraServiceManagementExtensionSettings types.List   `tfsdk:"jira_service_management_extension_settings"`
	ServiceNowExtensionSettings            types.List   `tfsdk:"servicenow_extension_settings"`
	IncludedLibraryVariableSets            types.List   `tfsdk:"included_library_variable_sets"`
	AutoDeployReleaseOverrides             types.List   `tfsdk:"auto_deploy_release_overrides"`
}

type connectivityPolicyModel struct {
	AllowDeploymentsToNoTargets types.Bool   `tfsdk:"allow_deployments_to_no_targets"`
	ExcludeUnhealthyTargets     types.Bool   `tfsdk:"exclude_unhealthy_targets"`
	SkipMachineBehavior         types.String `tfsdk:"skip_machine_behavior"`
	TargetRoles                 types.List   `tfsdk:"target_roles"`
}
type autoDeployReleaseOverrideModel struct {
	EnvironmentID types.String `tfsdk:"environment_id"`
	TenantID      types.String `tfsdk:"tenant_id"`
}

type gitPersistenceSettingsModel struct {
	URL               types.String `tfsdk:"url"`
	BasePath          types.String `tfsdk:"base_path"`
	DefaultBranch     types.String `tfsdk:"default_branch"`
	ProtectedBranches types.Set    `tfsdk:"protected_branches"`
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	GitCredentialID   types.String `tfsdk:"git_credential_id"`
}

type jiraServiceManagementExtensionSettingsModel struct {
	ConnectionID           types.String `tfsdk:"connection_id"`
	IsEnabled              types.Bool   `tfsdk:"is_enabled"`
	ServiceDeskProjectName types.String `tfsdk:"service_desk_project_name"`
}

type servicenowExtensionSettingsModel struct {
	ConnectionID                     types.String `tfsdk:"connection_id"`
	IsEnabled                        types.Bool   `tfsdk:"is_enabled"`
	IsStateAutomaticallyTransitioned types.Bool   `tfsdk:"is_state_automatically_transitioned"`
	StandardChangeTemplateName       types.String `tfsdk:"standard_change_template_name"`
}

type versioningStrategyModel struct {
	DonorPackageStepID types.String `tfsdk:"donor_package_step_id"`
	Template           types.String `tfsdk:"template"`
	DonorPackage       types.Object `tfsdk:"donor_package"`
}

type releaseCreationStrategyModel struct {
	ChannelID                    types.String `tfsdk:"channel_id"`
	ReleaseCreationPackageStepID types.String `tfsdk:"release_creation_package_step_id"`
	ReleaseCreationPackage       types.Object `tfsdk:"release_creation_package"`
}

type deploymentActionPackageModel struct {
	DeploymentAction types.String `tfsdk:"deployment_action"`
	PackageReference types.String `tfsdk:"package_reference"`
}

type templateModel struct {
	Name            types.String `tfsdk:"name"`
	Label           types.String `tfsdk:"label"`
	HelpText        types.String `tfsdk:"help_text"`
	DefaultValue    types.String `tfsdk:"default_value"`
	DisplaySettings types.Map    `tfsdk:"display_settings"`
}
