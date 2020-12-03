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

	if v, ok := d.GetOk("auto_create_release"); ok {
		project.AutoCreateRelease = v.(bool)
	}

	if v, ok := d.GetOk("auto_deploy_release_overrides"); ok {
		project.AutoDeployReleaseOverrides = v.([]*octopusdeploy.AutoDeployReleaseOverride)
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

	if v, ok := d.GetOk("templates"); ok {
		project.Templates = expandActionTemplateParameters(v.([]interface{}))
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		project.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	return project
}

func flattenProject(project *octopusdeploy.Project) map[string]interface{} {
	if project == nil {
		return nil
	}

	return map[string]interface{}{
		"auto_create_release":                  project.AutoCreateRelease,
		"auto_deploy_release_overrides":        project.AutoDeployReleaseOverrides,
		"cloned_from_project_id":               project.ClonedFromProjectID,
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
		"connectivity_policy":                  flattenConnectivityPolicy(project.ConnectivityPolicy),
		"project_group_id":                     project.ProjectGroupID,
		"release_creation_strategy":            flattenReleaseCreationStrategy(project.ReleaseCreationStrategy),
		"release_notes_template":               project.ReleaseNotesTemplate,
		"slug":                                 project.Slug,
		"space_id":                             project.SpaceID,
		"templates":                            flattenActionTemplateParameters(project.Templates),
		"tenanted_deployment_participation":    project.TenantedDeploymentMode,
		"variable_set_id":                      project.VariableSetID,
		"version_control_settings":             flattenVersionControlSettings(project.VersionControlSettings),
		"versioning_strategy":                  flattenVersioningStrategy(project.VersioningStrategy),
	}
}

func getProjectDataSchema() map[string]*schema.Schema {
	dataSchema := getProjectSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"cloned_from_project_id": getClonedFromProjectIDQuery(),
		"id":                     getIDDataSchema(),
		"ids":                    getIDsQuery(),
		"is_clone":               getIsCloneQuery(),
		"name":                   getNameQuery(),
		"partial_name":           getPartialNameQuery(),
		"projects": {
			Computed:    true,
			Description: "A list of projects that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"skip": getSkipQuery(),
		"take": getTakeQuery(),
	}
}

func getProjectSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_deployments_to_no_targets": {
			Deprecated: "Change this please!!!",
			Optional:   true,
			Type:       schema.TypeBool,
		},
		"auto_create_release": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"auto_deploy_release_overrides": {
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Type:     schema.TypeList,
		},
		"cloned_from_project_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"connectivity_policy": {
			Computed: true,
			Optional: true,
			Elem:     &schema.Resource{Schema: getConnectivityPolicySchema()},
			Type:     schema.TypeList,
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
			Computed: true,
			Type:     schema.TypeString,
		},
		"description": getDescriptionSchema(),
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
		"id": getIDSchema(),
		"included_library_variable_sets": {
			Elem:     &schema.Schema{Type: schema.TypeString},
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
		"name": getNameSchema(true),
		"project_group_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"release_creation_strategy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getReleaseCreationStrategySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"release_notes_template": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"slug": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"space_id": getSpaceIDSchema(),
		"templates": {
			Elem:     &schema.Resource{Schema: getActionTemplateParameterSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"variable_set_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"version_control_settings": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getVersionControlSettingsSchema()},
			Optional: true,
			Type:     schema.TypeSet,
		},
		"versioning_strategy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getVersionStrategySchema()},
			Optional: true,
			Type:     schema.TypeSet,
		},
	}
}

func setProject(ctx context.Context, d *schema.ResourceData, project *octopusdeploy.Project) {
	d.Set("auto_create_release", project.AutoCreateRelease)
	d.Set("auto_deploy_release_overrides", project.AutoDeployReleaseOverrides)
	d.Set("cloned_from_project_id", project.ClonedFromProjectID)
	d.Set("connectivity_policy", flattenConnectivityPolicy(project.ConnectivityPolicy))
	d.Set("default_guided_failure_mode", project.DefaultGuidedFailureMode)
	d.Set("default_to_skip_if_already_installed", project.DefaultToSkipIfAlreadyInstalled)
	d.Set("deployment_changes_template", project.DeploymentChangesTemplate)
	d.Set("deployment_process_id", project.DeploymentProcessID)
	d.Set("description", project.Description)
	d.Set("extension_settings", project.ExtensionSettings)
	d.Set("id", project.GetID())
	d.Set("included_library_variable_sets", project.IncludedLibraryVariableSets)
	d.Set("is_disabled", project.IsDisabled)
	d.Set("is_discrete_channel_release", project.IsDiscreteChannelRelease)
	d.Set("is_version_controlled", project.IsVersionControlled)
	d.Set("lifecycle_id", project.LifecycleID)
	d.Set("name", project.Name)
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
}
