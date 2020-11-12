package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectsRead,
		Schema:      getProjectDataSchema(),
	}
}

func dataSourceProjectsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.ProjectsQuery{
		ClonedFromProjectID: d.Get("cloned_from_project_id").(string),
		IDs:                 expandArray(d.Get("ids").([]interface{})),
		IsClone:             d.Get("is_clone").(bool),
		Name:                d.Get("name").(string),
		PartialName:         d.Get("partial_name").(string),
		Skip:                d.Get("skip").(int),
		Take:                d.Get("take").(int),
	}

	client := m.(*octopusdeploy.Client)
	projects, err := client.Projects.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedProjects := []interface{}{}
	for _, project := range projects.Items {
		flattenedProject := map[string]interface{}{
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
			"connectivity_policy":                  flattenProjectConnectivityPolicy(project.ConnectivityPolicy),
			"project_group_id":                     project.ProjectGroupID,
			"release_creation_strategy":            flattenReleaseCreationStrategy(project.ReleaseCreationStrategy),
			"release_notes_template":               project.ReleaseNotesTemplate,
			"slug":                                 project.Slug,
			"space_id":                             project.SpaceID,
			"templates":                            project.Templates,
			"tenanted_deployment_participation":    project.TenantedDeploymentMode,
			"variable_set_id":                      project.VariableSetID,
			"version_control_settings":             flattenVersionControlSettings(project.VersionControlSettings),
			"versioning_strategy":                  flattenVersioningStrategy(project.VersioningStrategy),
		}
		flattenedProjects = append(flattenedProjects, flattenedProject)
	}

	d.Set("projects", flattenedProjects)
	d.SetId("Projects " + time.Now().UTC().String())

	return nil
}
