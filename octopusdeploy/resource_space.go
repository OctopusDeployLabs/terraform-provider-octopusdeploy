package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
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

	log.Printf("[INFO] creating space: %#v", space)

	client := m.(*client.Client)
	createdSpace, err := client.Spaces.Add(space)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setSpace(ctx, d, createdSpace); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdSpace.GetID())

	log.Printf("[INFO] space created (%s)", d.Id())
	return nil
}

func resourceSpaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting space (%s)", d.Id())

	space := expandSpace(d)
	space.TaskQueueStopped = true

	client := m.(*client.Client)
	updatedSpace, err := spaces.Update(client, space)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.Spaces.DeleteByID(updatedSpace.GetID()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] space deleted")
	return nil
}

func resourceSpaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading space (%s)", d.Id())

	client := m.(*client.Client)
	space, err := spaces.GetByID(client, d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "space")
	}

	if err := setSpace(ctx, d, space); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] space read (%s)", d.Id())
	return nil
}

func resourceSpaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating space (%s)", d.Id())

	space := expandSpace(d)
	client := m.(*client.Client)
	updatedSpace, err := spaces.Update(client, space)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setSpace(ctx, d, updatedSpace); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] space updated (%s)", d.Id())
	return nil
}
