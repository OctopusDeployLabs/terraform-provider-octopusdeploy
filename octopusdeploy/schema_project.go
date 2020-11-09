package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandProject(d *schema.ResourceData) *octopusdeploy.Project {
	name := d.Get("name").(string)
	lifecycleID := d.Get("lifecycle_id").(string)
	projectGroupID := d.Get("project_group_id").(string)

	project := octopusdeploy.NewProject(name, lifecycleID, projectGroupID)
	project.ID = d.Id()

	if v, ok := d.GetOk("default_guided_failure_mode"); ok {
		project.DefaultGuidedFailureMode = v.(string)
	}

	if v, ok := d.GetOk("default_to_skip_if_already_installed"); ok {
		project.DefaultToSkipIfAlreadyInstalled = v.(bool)
	}

	if v, ok := d.GetOk("description"); ok {
		project.Description = v.(string)
	}

	if v, ok := d.GetOk("is_discrete_channel_release"); ok {
		project.IsDiscreteChannelRelease = v.(bool)
	}

	if v, ok := d.GetOk("included_library_variable_sets"); ok {
		project.IncludedLibraryVariableSets = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		project.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	return project
}

func flattenProject(ctx context.Context, d *schema.ResourceData, project *octopusdeploy.Project) {
	d.Set("auto_create_release", project.AutoCreateRelease)
	d.Set("auto_deploy_release_overrides", project.AutoDeployReleaseOverrides)
	d.Set("cloned_from_project_id", project.ClonedFromProjectID)
	d.Set("default_guided_failure_mode", project.DefaultGuidedFailureMode)
	d.Set("default_to_skip_if_already_installed", project.DefaultToSkipIfAlreadyInstalled)
	d.Set("deployment_changes_template", project.DeploymentChangesTemplate)
	d.Set("deployment_process_id", project.DeploymentProcessID)
	d.Set("description", project.Description)
	d.Set("is_discrete_channel_release", project.IsDiscreteChannelRelease)
	d.Set("extension_settings", project.ExtensionSettings)
	d.Set("included_library_variable_sets", project.IncludedLibraryVariableSets)
	d.Set("is_disabled", project.IsDisabled)
	d.Set("is_discrete_channel_release", project.IsDiscreteChannelRelease)
	d.Set("is_version_controlled", project.IsVersionControlled)
	d.Set("lifecycle_id", project.LifecycleID)
	d.Set("name", project.Name)
	d.Set("project_connectivity_policy", flattenProjectConnectivityPolicy(project.ProjectConnectivityPolicy))
	d.Set("project_group_id", project.ProjectGroupID)
	d.Set("release_creation_strategy", flattenReleaseCreationStrategy(project.ReleaseCreationStrategy))
	d.Set("release_notes_template", project.ReleaseNotesTemplate)
	d.Set("slug", project.Slug)
	d.Set("space_id", project.SpaceID)
	d.Set("templates", project.Templates)
	d.Set("tenanted_deployment_participation", project.TenantedDeploymentMode)
	d.Set("variable_set_id", project.VariableSetID)
	d.Set("version_control_settings", flattenVersionControlSettings(project.VersionControlSettings))
	d.Set("versioning_strategy", flattenVersioningStrategy(project.VersioningStrategy))

	d.SetId(project.GetID())
}

func flattenProjectConnectivityPolicy(projectConnectivityPolicy *octopusdeploy.ProjectConnectivityPolicy) []interface{} {
	if projectConnectivityPolicy == nil {
		return nil
	}

	flattenedProjectConnectivityPolicy := make(map[string]interface{})
	flattenedProjectConnectivityPolicy["allow_deployments_to_no_targets"] = projectConnectivityPolicy.AllowDeploymentsToNoTargets
	flattenedProjectConnectivityPolicy["skip_machine_behavior"] = projectConnectivityPolicy.SkipMachineBehavior
	flattenedProjectConnectivityPolicy["target_roles"] = projectConnectivityPolicy.TargetRoles
	return []interface{}{flattenedProjectConnectivityPolicy}
}

func flattenReleaseCreationStrategy(releaseCreationStrategy *octopusdeploy.ReleaseCreationStrategy) []interface{} {
	if releaseCreationStrategy == nil {
		return nil
	}

	flattenedReleaseCreationStrategy := make(map[string]interface{})
	flattenedReleaseCreationStrategy["channel_id"] = releaseCreationStrategy.ChannelID
	flattenedReleaseCreationStrategy["release_creation_package"] = releaseCreationStrategy.ReleaseCreationPackage
	flattenedReleaseCreationStrategy["release_creation_package_step_id"] = releaseCreationStrategy.ReleaseCreationPackageStepID
	return []interface{}{flattenedReleaseCreationStrategy}
}

func flattenVersionControlSettings(versionControlSettings *octopusdeploy.VersionControlSettings) []interface{} {
	if versionControlSettings == nil {
		return nil
	}

	flattenedVersionControlSettings := make(map[string]interface{})
	flattenedVersionControlSettings["default_branch"] = versionControlSettings.DefaultBranch
	flattenedVersionControlSettings["password"] = versionControlSettings.Password
	flattenedVersionControlSettings["url"] = versionControlSettings.URL
	flattenedVersionControlSettings["username"] = versionControlSettings.Username
	return []interface{}{flattenedVersionControlSettings}
}

func flattenVersioningStrategy(versioningStrategy octopusdeploy.VersioningStrategy) []interface{} {
	flattenedVersioningStrategy := make(map[string]interface{})
	flattenedVersioningStrategy["donor_package"] = versioningStrategy.DonorPackage
	flattenedVersioningStrategy["donor_package_step_id"] = versioningStrategy.DonorPackageStepID
	flattenedVersioningStrategy["template"] = versioningStrategy.Template
	return []interface{}{flattenedVersioningStrategy}
}

func getProjectDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"auto_create_release": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"auto_deploy_release_overrides": {
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Type: schema.TypeList,
		},
		"cloned_from_project_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"default_guided_failure_mode": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"default_to_skip_if_already_installed": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"deployment_changes_template": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"deployment_process_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"description": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"extension_settings": {
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Type: schema.TypeList,
		},
		"included_library_variable_sets": {
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Type: schema.TypeList,
		},
		"is_disabled": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"is_discrete_channel_release": {
			Computed:    true,
			Description: "Treats releases of different channels to the same environment as a separate deployment dimension",
			Type:        schema.TypeBool,
		},
		"is_version_controlled": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"lifecycle_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"project_connectivity_policy": {
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"allow_deployments_to_no_targets": {
						Computed: true,
						Type:     schema.TypeBool,
					},
					"skip_machine_behavior": {
						Computed: true,
						Type:     schema.TypeString,
					},
					"target_roles": {
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Computed: true,
						Type:     schema.TypeList,
					},
				},
			},
			Type: schema.TypeSet,
		},
		"project_group_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"release_creation_strategy": {
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Type: schema.TypeList,
		},
		"release_notes_template": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"slug": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"space_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"templates": {
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Type: schema.TypeList,
		},
		"tenanted_deployment_participation": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"variable_set_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"version_control_settings": {
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"default_branch": {
						Computed: true,
						Type:     schema.TypeString,
					},
					"password": {
						Computed:  true,
						Sensitive: true,
						Type:      schema.TypeString,
					},
					"url": {
						Computed: true,
						Type:     schema.TypeString,
					},
					"username": {
						Computed:  true,
						Sensitive: true,
						Type:      schema.TypeString,
					},
				},
			},
			Type: schema.TypeSet,
		},
		"versioning_strategy": {
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"donor_package": {
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Type: schema.TypeList,
					},
					"donor_package_step_id": {
						Computed: true,
						Type:     schema.TypeString,
					},
					"template": {
						Computed: true,
						Type:     schema.TypeString,
					},
				},
			},
			Type: schema.TypeSet,
		},
	}
}

func getProjectSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"auto_create_release": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"auto_deploy_release_overrides": {
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Type: schema.TypeList,
		},
		"cloned_from_project_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"default_guided_failure_mode": {
			Optional: true,
			Type:     schema.TypeString,
			Default:  "EnvironmentDefault",
			ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
				"EnvironmentDefault",
				"Off",
				"On",
			}, false)),
		},
		"default_to_skip_if_already_installed": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"deployment_changes_template": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"deployment_process_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"discrete_channel_release": {
			Description: "Treats releases of different channels to the same environment as a separate deployment dimension",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"extension_settings": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		"included_library_variable_sets": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		"is_disabled": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"is_discrete_channel_release": {
			Optional:    true,
			Description: "Treats releases of different channels to the same environment as a separate deployment dimension",
			Type:        schema.TypeBool,
		},
		"is_version_controlled": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"lifecycle_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"project_connectivity_policy": {
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"allow_deployments_to_no_targets": {
						Optional: true,
						Type:     schema.TypeBool,
					},
					"skip_machine_behavior": {
						Default:  "None",
						Optional: true,
						Type:     schema.TypeString,
						ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
							"SkipUnavailableMachines",
							"None",
						}, false)),
					},
					"target_roles": {
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Optional: true,
						Type:     schema.TypeList,
					},
				},
			},
			Type: schema.TypeSet,
		},
		"project_group_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"release_creation_strategy": {
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"channel_id": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"release_creation_package": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"release_creation_package_step_id": {
						Optional: true,
						Type:     schema.TypeString,
					},
				},
			},
			Type: schema.TypeSet,
		},
		"release_notes_template": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"slug": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"space_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"templates": {
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Type: schema.TypeList,
		},
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"variable_set_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"version_control_settings": {
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"default_branch": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"password": {
						Optional:  true,
						Sensitive: true,
						Type:      schema.TypeString,
					},
					"url": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"username": {
						Optional:  true,
						Sensitive: true,
						Type:      schema.TypeString,
					},
				},
			},
			Type: schema.TypeSet,
		},
		"versioning_strategy": {
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"donor_package": {
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Type: schema.TypeList,
					},
					"donor_package_step_id": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"template": {
						Optional: true,
						Type:     schema.TypeString,
					},
				},
			},
			Type: schema.TypeSet,
		},
	}
}
