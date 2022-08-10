package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectGroupCreate,
		DeleteContext: resourceProjectGroupDelete,
		Description:   "This resource manages project groups in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceProjectGroupRead,
		Schema:        getProjectGroupSchema(),
		UpdateContext: resourceProjectGroupUpdate,
	}
}

func resourceProjectGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectGroup := expandProjectGroup(d)

	log.Printf("[INFO] creating project group: %#v", projectGroup)

	client := m.(*client.Client)
	createdProjectGroup, err := client.ProjectGroups.Add(projectGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setProjectGroup(ctx, d, createdProjectGroup); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdProjectGroup.GetID())

	log.Printf("[INFO] project group created (%s)", d.Id())
	return nil
}

func resourceProjectGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting project group (%s)", d.Id())

	client := m.(*client.Client)
	if err := client.ProjectGroups.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] project group deleted (%s)", d.Id())
	d.SetId("")
	return nil
}

func resourceProjectGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading project group (%s)", d.Id())

	client := m.(*client.Client)
	projectGroup, err := client.ProjectGroups.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "project group")
	}

	if err := setProjectGroup(ctx, d, projectGroup); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] project group read (%s)", d.Id())
	return nil
}

func resourceProjectGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating project group (%s)", d.Id())

	projectGroup := expandProjectGroup(d)
	client := m.(*client.Client)
	updatedProjectGroup, err := client.ProjectGroups.Update(*projectGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setProjectGroup(ctx, d, updatedProjectGroup); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] project group updated (%s)", d.Id())
	return nil
}
