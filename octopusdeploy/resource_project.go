package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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

	log.Printf("[INFO] creating project: %#v", project)

	client := m.(*octopusdeploy.Client)
	createdProject, err := client.Projects.Add(project)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setProject(ctx, d, createdProject); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdProject.GetID())

	log.Printf("[INFO] project created (%s)", d.Id())
	return nil
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting project (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Projects.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] project deleted")
	return nil
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading project (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	project, err := client.Projects.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			log.Printf("[INFO] project (%s) not found; deleting from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := setProject(ctx, d, project); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] project read (%s)", d.Id())
	return nil
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating project (%s)", d.Id())

	project := expandProject(d)
	client := m.(*octopusdeploy.Client)
	updatedProject, err := client.Projects.Update(project)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setProject(ctx, d, updatedProject); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] project updated (%s)", d.Id())
	return nil
}
