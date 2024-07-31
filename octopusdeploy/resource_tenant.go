package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
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
	mutex.Lock()
	defer mutex.Unlock()

	tenant := expandTenant(d)

	log.Printf("[INFO] creating tenant: %#v", tenant)

	client := m.(*client.Client)
	createdTenant, err := tenants.Add(client, tenant)
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
	mutex.Lock()
	defer mutex.Unlock()

	log.Printf("[INFO] deleting tenant (%s)", d.Id())

	client := m.(*client.Client)
	if err := tenants.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] tenant deleted (%s)", d.Id())
	d.SetId("")
	return nil
}

func resourceTenantRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading tenant (%s)", d.Id())

	client := m.(*client.Client)
	tenant, err := tenants.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "tenant")
	}

	if err := setTenant(ctx, d, tenant); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] tenant read (%s)", d.Id())
	return nil
}

func resourceTenantUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	log.Printf("[INFO] updating tenant (%s)", d.Id())

	client := m.(*client.Client)
	tenantFromApi, err := tenants.GetByID(client, d.Get("space_id").(string), d.Id())

	tenant := expandTenant(d)

	// the project environments are not managed here, so we need to maintain the collection when updating
	tenant.ProjectEnvironments = tenantFromApi.ProjectEnvironments
	updatedTenant, err := tenants.Update(client, tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTenant(ctx, d, updatedTenant); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] tenant updated (%s)", d.Id())
	return nil
}
