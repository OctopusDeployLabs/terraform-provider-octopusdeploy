package octopusdeploy

import (
	"context"
	"fmt"

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

	if v, ok := d.GetOk("auto_create_release"); ok {
		project.AutoCreateRelease = v.(bool)
	}

	if v, ok := d.GetOk("auto_deploy_release_overrides"); ok {
		project.AutoDeployReleaseOverrides = v.([]octopusdeploy.AutoDeployReleaseOverride)
	}

	if v, ok := d.GetOk("cloned_from_project_id"); ok {
		project.ClonedFromProjectID = v.(string)
	}

	if v, ok := d.GetOk("connectivity_policy"); ok {
		project.ConnectivityPolicy = expandConnectivityPolicy(v.([]interface{}))
	}

	if v, ok := d.GetOk("default_guided_failure_mode"); ok {
		project.DefaultGuidedFailureMode = v.(string)
	}

	if v, ok := d.GetOk("default_to_skip_if_already_installed"); ok {
		project.DefaultToSkipIfAlreadyInstalled = v.(bool)
	}

	if v, ok := d.GetOk("deployment_changes_template"); ok {
		project.DeploymentChangesTemplate = v.(string)
	}

	if v, ok := d.GetOk("deployment_process_id"); ok {
		project.DeploymentProcessID = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		project.Description = v.(string)
	}

	if v, ok := d.GetOk("extension_settings"); ok {
		project.ExtensionSettings = expandExtensionSettingsValues(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("git_persistence_settings"); ok {
		project.PersistenceSettings = expandGitPersistenceSettings(v)
	}

	if v, ok := d.GetOk("included_library_variable_sets"); ok {
		project.IncludedLibraryVariableSets = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("is_disabled"); ok {
		project.IsDisabled = v.(bool)
	}

	if v, ok := d.GetOk("is_discrete_channel_release"); ok {
		project.IsDiscreteChannelRelease = v.(bool)
	}

	if v, ok := d.GetOk("is_version_controlled"); ok {
		project.IsVersionControlled = v.(bool)
	}

	if v, ok := d.GetOk("release_creation_strategy"); ok {
		project.ReleaseCreationStrategy = expandReleaseCreationStrategy(v.([]interface{}))
	}

	if v, ok := d.GetOk("release_notes_template"); ok {
		project.ReleaseNotesTemplate = v.(string)
	}

	if v, ok := d.GetOk("slug"); ok {
		project.Slug = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		project.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("template"); ok {
		project.Templates = expandActionTemplateParameters(v.([]interface{}))
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		project.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("versioning_strategy"); ok {
		project.VersioningStrategy = expandVersioningStrategy(v)
	}

	return project
}

func flattenProject(project *octopusdeploy.Project) map[string]interface{} {
	if project == nil {
		return nil
	}

	projectMap := map[string]interface{}{
		"auto_create_release":                  project.AutoCreateRelease,
		"auto_deploy_release_overrides":        project.AutoDeployReleaseOverrides,
		"cloned_from_project_id":               project.ClonedFromProjectID,
		"connectivity_policy":                  flattenConnectivityPolicy(project.ConnectivityPolicy),
		"default_guided_failure_mode":          project.DefaultGuidedFailureMode,
		"default_to_skip_if_already_installed": project.DefaultToSkipIfAlreadyInstalled,
		"deployment_changes_template":          project.DeploymentChangesTemplate,
		"deployment_process_id":                project.DeploymentProcessID,
		"description":                          project.Description,
		"extension_settings":                   project.ExtensionSettings,
		"id":                                   project.GetID(),
		"included_library_variable_sets":       project.IncludedLibraryVariableSets,
		"is_disabled":                          project.IsDisabled,
		"is_discrete_channel_release":          project.IsDiscreteChannelRelease,
		"is_version_controlled":                project.IsVersionControlled,
		"lifecycle_id":                         project.LifecycleID,
		"name":                                 project.Name,
		"project_group_id":                     project.ProjectGroupID,
		"release_creation_strategy":            flattenReleaseCreationStrategy(project.ReleaseCreationStrategy),
		"release_notes_template":               project.ReleaseNotesTemplate,
		"slug":                                 project.Slug,
		"space_id":                             project.SpaceID,
		"template":                             flattenActionTemplateParameters(project.Templates),
		"tenanted_deployment_participation":    project.TenantedDeploymentMode,
		"variable_set_id":                      project.VariableSetID,
		"versioning_strategy":                  flattenVersioningStrategy(project.VersioningStrategy),
	}

	if project.PersistenceSettings != nil {
		if project.PersistenceSettings.GetType() == "VersionControlled" {
			projectMap["git_persistence_settings"] = flattenGitPersistenceSettings(project.PersistenceSettings, projectMap["git_persistence_settings.0.credentials.0.password"].(string))
		}
	}

	return projectMap
}

func getProjectDataSchema() map[string]*schema.Schema {
	dataSchema := getProjectSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"cloned_from_project_id": getQueryClonedFromProjectID(),
		"id":                     getDataSchemaID(),
		"ids":                    getQueryIDs(),
		"is_clone":               getQueryIsClone(),
		"name":                   getQueryName(),
		"partial_name":           getQueryPartialName(),
		"projects": {
			Computed:    true,
			Description: "A list of projects that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"skip": getQuerySkip(),
		"take": getQueryTake(),
	}
}

func getProjectSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_deployments_to_no_targets": {
			Deprecated: "This value is only valid for an associated connectivity policy and should not be specified here.",
			Optional:   true,
			Type:       schema.TypeBool,
		},
		"auto_create_release": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"auto_deploy_release_overrides": {
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"cloned_from_project_id": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"connectivity_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getConnectivityPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"default_guided_failure_mode": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"EnvironmentDefault",
				"Off",
				"On",
			}, false)),
		},
		"default_to_skip_if_already_installed": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"deployment_changes_template": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"deployment_process_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"description": {
			Computed:    true,
			Description: "The description of this project.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"discrete_channel_release": {
			Description: "Treats releases of different channels to the same environment as a separate deployment dimension",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"extension_settings": {
			Optional: true,
			Elem:     &schema.Resource{Schema: getExtensionSettingsSchema()},
			Type:     schema.TypeSet,
		},
		"git_persistence_settings": {
			Description: "Provides Git-related persistence settings for a version-controlled project.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"base_path": {
						Default:     ".octopus",
						Description: "The base path associated with these version control settings.",
						Optional:    true,
						Type:        schema.TypeString,
					},
					"credentials": {
						Description: "The credentials associated with these version control settings.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"password": {
									Description:      "The password for the Git credential.",
									Required:         true,
									Sensitive:        true,
									Type:             schema.TypeString,
									ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
								},
								"username": {
									Description:      "The username for the Git credential.",
									Required:         true,
									Type:             schema.TypeString,
									ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
								},
							},
						},
						MaxItems: 1,
						Optional: true,
						Type:     schema.TypeList,
					},
					"default_branch": {
						Default:     "main",
						Description: "The default branch associated with these version control settings.",
						Optional:    true,
						Type:        schema.TypeString,
					},
					"url": {
						Description:      "The URL associated with these version control settings.",
						Required:         true,
						Type:             schema.TypeString,
						ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
					},
				},
			},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"id": getIDSchema(),
		"included_library_variable_sets": {
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"is_disabled": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"is_discrete_channel_release": {
			Computed:    true,
			Description: "Treats releases of different channels to the same environment as a separate deployment dimension",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"is_version_controlled": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"lifecycle_id": {
			Description:      "The lifecycle ID associated with this project.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"name": {
			Description:      "The name of the project in Octopus Deploy. This name must be unique.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"project_group_id": {
			Description:      "The project group ID associated with this project.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"release_creation_strategy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getReleaseCreationStrategySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"release_notes_template": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"slug": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"space_id": {
			Computed:         true,
			Description:      "The space ID associated with this project.",
			Optional:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"template": {
			Elem:     &schema.Resource{Schema: getActionTemplateParameterSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"variable_set_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"versioning_strategy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getVersionStrategySchema()},
			Optional: true,
			Type:     schema.TypeSet,
		},
	}
}

func setProject(ctx context.Context, d *schema.ResourceData, project *octopusdeploy.Project) error {
	d.Set("auto_create_release", project.AutoCreateRelease)

	if err := d.Set("auto_deploy_release_overrides", project.AutoDeployReleaseOverrides); err != nil {
		return fmt.Errorf("error setting auto_deploy_release_overrides: %s", err)
	}

	d.Set("cloned_from_project_id", project.ClonedFromProjectID)

	if err := d.Set("connectivity_policy", flattenConnectivityPolicy(project.ConnectivityPolicy)); err != nil {
		return fmt.Errorf("error setting connectivity_policy: %s", err)
	}

	d.Set("default_guided_failure_mode", project.DefaultGuidedFailureMode)
	d.Set("default_to_skip_if_already_installed", project.DefaultToSkipIfAlreadyInstalled)
	d.Set("deployment_changes_template", project.DeploymentChangesTemplate)
	d.Set("deployment_process_id", project.DeploymentProcessID)
	d.Set("description", project.Description)

	if err := d.Set("extension_settings", project.ExtensionSettings); err != nil {
		return fmt.Errorf("error setting extension_settings: %s", err)
	}

	if err := d.Set("included_library_variable_sets", project.IncludedLibraryVariableSets); err != nil {
		return fmt.Errorf("error setting included_library_variable_sets: %s", err)
	}

	d.Set("is_disabled", project.IsDisabled)
	d.Set("is_discrete_channel_release", project.IsDiscreteChannelRelease)
	d.Set("is_version_controlled", project.IsVersionControlled)
	d.Set("lifecycle_id", project.LifecycleID)
	d.Set("name", project.Name)

	if project.PersistenceSettings != nil {
		if project.PersistenceSettings.GetType() == "VersionControlled" {
			if err := d.Set("git_persistence_settings", flattenGitPersistenceSettings(project.PersistenceSettings, d.Get("git_persistence_settings.0.credentials.0.password").(string))); err != nil {
				return fmt.Errorf("error setting git_persistence_settings: %s", err)
			}
		}
	}

	d.Set("project_group_id", project.ProjectGroupID)

	if err := d.Set("release_creation_strategy", flattenReleaseCreationStrategy(project.ReleaseCreationStrategy)); err != nil {
		return fmt.Errorf("error setting release_creation_strategy: %s", err)
	}

	d.Set("release_notes_template", project.ReleaseNotesTemplate)
	d.Set("slug", project.Slug)
	d.Set("space_id", project.SpaceID)

	if err := d.Set("template", flattenActionTemplateParameters(project.Templates)); err != nil {
		return fmt.Errorf("error setting templates: %s", err)
	}

	d.Set("tenanted_deployment_participation", project.TenantedDeploymentMode)
	d.Set("variable_set_id", project.VariableSetID)

	if err := d.Set("versioning_strategy", flattenVersioningStrategy(project.VersioningStrategy)); err != nil {
		return fmt.Errorf("error setting versioning_strategy: %s", err)
	}

	d.Set("id", project.GetID())

	return nil
}
