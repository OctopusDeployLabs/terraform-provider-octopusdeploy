package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		DeleteContext: resourceProjectDelete,
		Description:   "This resource manages projects in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceProjectRead,
		Schema:        getProjectSchema(),
		UpdateContext: resourceProjectUpdate,
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	project := expandProject(ctx, d)

	// DANGER: the go provider is about to nil the persistence settings, to stop the API from exploding. Take a copy
	// so we can make decisions.
	persistenceSettings := project.PersistenceSettings

	tflog.Info(ctx, fmt.Sprintf("creating project (%s)", project.Name))

	var spaceID string
	if v, ok := d.GetOk("space_id"); ok {
		spaceID = v.(string)
	}

	client := m.(*client.Client)
	createdProject, err := projects.Add(client, project)
	if err != nil {
		return diag.FromErr(err)
	}

	if persistenceSettings != nil && persistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		tflog.Info(ctx, "converting project to use VCS")

		vcsProject, err := projects.ConvertToVCS(client, createdProject, "converting project to use VCS", "", persistenceSettings.(projects.GitPersistenceSettings))
		if err != nil {
			projects.DeleteByID(client, spaceID, createdProject.GetID())
			return diag.FromErr(err)
		}
		createdProject.PersistenceSettings = vcsProject.PersistenceSettings
	}

	createdProject, err = projects.GetByID(client, spaceID, createdProject.GetID())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setProject(ctx, d, createdProject); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdProject.GetID())

	tflog.Info(ctx, fmt.Sprintf("project created (%s)", d.Id()))
	return nil
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("deleting project (%s)", d.Id()))

	var spaceID string
	if v, ok := d.GetOk("space_id"); ok {
		spaceID = v.(string)
	}

	client := m.(*client.Client)
	if err := projects.DeleteByID(client, spaceID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("project deleted (%s)", d.Id()))
	d.SetId("")
	return nil
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("reading project (%s)", d.Id()))

	var spaceID string
	if v, ok := d.GetOk("space_id"); ok {
		spaceID = v.(string)
	}

	client := m.(*client.Client)
	project, err := projects.GetByID(client, spaceID, d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "project")
	}

	if err := setProject(ctx, d, project); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("project read (%s)", d.Id()))
	return nil
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("updating project (%s)", d.Id()))

	client := m.(*client.Client)
	project := expandProject(ctx, d)
	var updatedProject *projects.Project
	var err error

	projectLinks, err := projects.GetByID(client, project.SpaceID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if project.PersistenceSettings != nil && project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		convertToVcsLink := projectLinks.Links["ConvertToVcs"]

		if len(convertToVcsLink) != 0 {
			versionControlSettings := expandVersionControlSettingsForProjectConversion(ctx, d)

			tflog.Info(ctx, fmt.Sprintf("converting project to use VCS (%s)", d.Id()))

			project.Links["ConvertToVcs"] = convertToVcsLink
			vcsProject, err := projects.ConvertToVCS(client, project, "converting project to use VCS", "", versionControlSettings)
			if err != nil {
				return diag.FromErr(err)
			}
			project.PersistenceSettings = vcsProject.PersistenceSettings
		}
	}

	project.Links = projectLinks.Links

	updatedProject, err = projects.Update(client, project)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setProject(ctx, d, updatedProject); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("project updated (%s)", d.Id()))
	return nil
}
