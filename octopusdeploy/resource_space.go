package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSpace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSpaceCreate,
		DeleteContext: resourceSpaceDelete,
		Description:   "This resource manages spaces in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceSpaceRead,
		Schema:        getSpaceSchema(),
		UpdateContext: resourceSpaceUpdate,
	}
}

func resourceSpaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	space := expandSpace(d)

	client := m.(*octopusdeploy.Client)
	createdSpace, err := client.Spaces.Add(space)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdSpace.GetID())
	return resourceSpaceRead(ctx, d, m)
}

func resourceSpaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	space := expandSpace(d)
	space.TaskQueueStopped = true

	client := m.(*octopusdeploy.Client)
	updatedSpace, err := client.Spaces.Update(space)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Spaces.DeleteByID(updatedSpace.GetID())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceSpaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	space, err := client.Spaces.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	setSpace(ctx, d, space)
	return nil
}

func resourceSpaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	space := expandSpace(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Spaces.Update(space)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceChannelRead(ctx, d, m)
}
