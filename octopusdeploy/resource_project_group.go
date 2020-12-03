package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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

	client := m.(*octopusdeploy.Client)
	createdProjectGroup, err := client.ProjectGroups.Add(projectGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	setProjectGroup(ctx, d, createdProjectGroup)
	return nil
}

func resourceProjectGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.ProjectGroups.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceProjectGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	projectGroup, err := client.ProjectGroups.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	setProjectGroup(ctx, d, projectGroup)
	return nil
}

func resourceProjectGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectGroup := expandProjectGroup(d)

	client := m.(*octopusdeploy.Client)
	updatedProjectGroup, err := client.ProjectGroups.Update(*projectGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	setProjectGroup(ctx, d, updatedProjectGroup)
	return nil
}
