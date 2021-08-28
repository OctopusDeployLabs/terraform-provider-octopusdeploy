package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectTriggerCreate,
		DeleteContext: resourceProjectTriggerDelete,
		Description:   "This resource manages project triggers in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceProjectTriggerRead,
		Schema:        getProjectTriggerSchema(),
		UpdateContext: resourceProjectTriggerUpdate,
	}
}

func resourceProjectTriggerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectTrigger := expandProjectTrigger(d)

	log.Printf("[INFO] creating project trigger: %#v", projectTrigger)

	client := m.(*octopusdeploy.Client)
	createdProjectTrigger, err := client.ProjectTriggers.Add(projectTrigger)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setProjectTrigger(ctx, d, createdProjectTrigger); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdProjectTrigger.GetID())

	log.Printf("[INFO] project trigger created (%s)", d.Id())
	return nil
}

func resourceProjectTriggerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting project trigger (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.ProjectTriggers.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] project trigger deleted (%s)", d.Id())
	d.SetId("")
	return nil
}

func resourceProjectTriggerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading project trigger (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	projectTrigger, err := client.ProjectTriggers.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] project trigger (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if err := setProjectTrigger(ctx, d, projectTrigger); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] project trigger read (%s)", d.Id())
	return nil
}

func resourceProjectTriggerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating project trigger (%s)", d.Id())

	projectTrigger := expandProjectTrigger(d)
	client := m.(*octopusdeploy.Client)
	updatedProjectTrigger, err := client.ProjectTriggers.Update(*projectTrigger)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setProjectTrigger(ctx, d, updatedProjectTrigger); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] project trigger updated (%s)", d.Id())
	return nil
}
