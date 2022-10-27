package octopusdeploy

import (
	"context"
	"fmt"
	"log"

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
	project := expandProject(d)

	tflog.Info(ctx, fmt.Sprintf("creating project (%s)", project.Name))

	client := m.(*client.Client)
	createdProject, err := client.Projects.Add(project)
	if err != nil {
		return diag.FromErr(err)
	}

	if v, ok := d.GetOk("git_persistence_settings"); ok {
		versionControlSettings := expandVersionControlSettings(v)
		if versionControlSettings.Type == "VersionControlled" {
			log.Printf("[INFO] converting project to use VCS (%s)", d.Id())
			vcsProject, err := client.Projects.ConvertToVcs(createdProject, "converting project to use VCS", versionControlSettings)
			if err != nil {
				client.Projects.DeleteByID(createdProject.GetID())
				return diag.FromErr(err)
			}
			createdProject.PersistenceSettings = vcsProject.PersistenceSettings
		}
	}

	createdProject, err = client.Projects.GetByID(createdProject.GetID())
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
	log.Printf("[INFO] deleting project (%s)", d.Id())

	client := m.(*client.Client)
	if err := client.Projects.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] project deleted (%s)", d.Id())
	d.SetId("")
	return nil
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("reading project (%s)", d.Id()))

	client := m.(*client.Client)
	project, err := client.Projects.GetByID(d.Id())
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
	log.Printf("[INFO] updating project (%s)", d.Id())

	client := m.(*client.Client)
	project := expandProject(d)
	var updatedProject *projects.Project
	var err error

	projectLinks, err := client.Projects.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if project.PersistenceSettings != nil && project.PersistenceSettings.GetType() == "VersionControlled" {
		if v, ok := d.GetOk("git_persistence_settings"); ok {
			convertToVcsLink := projectLinks.Links["ConvertToVcs"]

			if len(convertToVcsLink) != 0 {
				versionControlSettings := expandVersionControlSettings(v)
				project.Links["ConvertToVcs"] = convertToVcsLink
				log.Printf("[INFO] converting project to use VCS (%s)", d.Id())
				project, err = client.Projects.ConvertToVcs(project, "converting project to use VCS", versionControlSettings)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	project.Links = projectLinks.Links

	updatedProject, err = client.Projects.Update(project)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setProject(ctx, d, updatedProject); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] project updated (%s)", d.Id())
	return nil
}
