package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTenant() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTenantCreate,
		DeleteContext: resourceTenantDelete,
		Description:   "This resource manages tenants in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceTenantRead,
		Schema:        getTenantSchema(),
		UpdateContext: resourceTenantUpdate,
	}
}

func resourceTenantCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenant := expandTenant(d)

	log.Printf("[INFO] creating tenant: %#v", tenant)

	client := m.(*octopusdeploy.Client)
	createdTenant, err := client.Tenants.Add(tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTenant(ctx, d, createdTenant); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdTenant.GetID())

	log.Printf("[INFO] tenant created (%s)", d.Id())
	return nil
}

func resourceTenantDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting tenant (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Tenants.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] tenant deleted (%s)", d.Id())
	d.SetId("")
	return nil
}

func resourceTenantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading tenant (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	tenant, err := client.Tenants.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] tenant (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if err := setTenant(ctx, d, tenant); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] tenant read (%s)", d.Id())
	return nil
}

func resourceTenantUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating tenant (%s)", d.Id())

	tenant := expandTenant(d)
	client := m.(*octopusdeploy.Client)
	updatedTenant, err := client.Tenants.Update(tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTenant(ctx, d, updatedTenant); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] tenant updated (%s)", d.Id())
	return nil
}
