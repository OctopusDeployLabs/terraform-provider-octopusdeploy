package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLifecycle() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLifecycleCreate,
		DeleteContext: resourceLifecycleDelete,
		Description:   "This resource manages lifecycles in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceLifecycleRead,
		Schema:        getLifecycleSchema(),
		UpdateContext: resourceLifecycleUpdate,
	}
}

func resourceLifecycleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lifecycle := expandLifecycle(d)

	log.Printf("[INFO] creating lifecycle: %#v", lifecycle)

	client := m.(*client.Client)
	createdLifecycle, err := lifecycles.Add(client, lifecycle)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setLifecycle(ctx, d, createdLifecycle); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdLifecycle.GetID())

	log.Printf("[INFO] lifecycle created (%s)", d.Id())
	return nil
}

func resourceLifecycleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting lifecycle (%s)", d.Id())

	client := m.(*client.Client)
	if err := lifecycles.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] lifecycle deleted (%s)", d.Id())
	d.SetId("")
	return nil
}

func resourceLifecycleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading lifecycle (%s)", d.Id())

	client := m.(*client.Client)
	lifecycle, err := lifecycles.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "lifecycle")
	}

	if err := setLifecycle(ctx, d, lifecycle); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] lifecycle read (%s)", d.Id())
	return nil
}

func resourceLifecycleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating lifecycle (%s)", d.Id())

	lifecycle := expandLifecycle(d)

	client := m.(*client.Client)
	updatedLifecycle, err := lifecycles.Update(client, lifecycle)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setLifecycle(ctx, d, updatedLifecycle); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] lifecycle updated (%s)", d.Id())
	return nil
}
