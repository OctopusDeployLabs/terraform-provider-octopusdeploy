package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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

	client := m.(*octopusdeploy.Client)
	createdLifecycle, err := client.Lifecycles.Add(lifecycle)
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

	client := m.(*octopusdeploy.Client)
	err := client.Lifecycles.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] lifecycle deleted")
	return nil
}

func resourceLifecycleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading lifecycle (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	lifecycle, err := client.Lifecycles.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			log.Printf("[INFO] lifecycle (%s) not found; deleting from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
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

	client := m.(*octopusdeploy.Client)
	updatedLifecycle, err := client.Lifecycles.Update(lifecycle)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setLifecycle(ctx, d, updatedLifecycle); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] lifecycle updated (%s)", d.Id())
	return nil
}
