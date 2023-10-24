package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTenantConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTenantConnectionCreate,
		DeleteContext: resourceTenantConnectionDelete,
		Description:   "This resource manages tenant connections in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceTenantConnectionRead,
		Schema:        getTenantSchema(),
		UpdateContext: resourceTenantConnectionUpdate,
	}
}

func resourceTenantConnectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tenantConnection := expandTenantConnection(d)

	log.Printf("[INFO] creating tenant connection: %#v", tenantConnection)

	client := m.(*client.Client)
	tenant, err := client.Tenants.GetByID(tenantConnection.TenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	tenant.ProjectEnvironments[tenantConnection.ProjectID] = tenantConnection.EnvironmentIDs
	tenant, err = client.Tenants.Update(tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	tenantConnection, err = expandTenantConnectionFromTenant(tenant, tenantConnection.ProjectID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTenantConnection(ctx, d, tenantConnection); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tenantConnection.GetID())

	log.Printf("[INFO] tenant connection created (%s)", d.Id())
	return nil
}

func resourceTenantConnectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting tenant connection (%s)", d.Id())

	tenantConnection, err := expandTenantConnectionFromID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := m.(*client.Client)
	tenant, err := client.Tenants.GetByID(tenantConnection.TenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	if _, ok := tenant.ProjectEnvironments[tenantConnection.ProjectID]; ok {
		delete(tenant.ProjectEnvironments, tenantConnection.ProjectID)
		if _, err := client.Tenants.Update(tenant); err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[INFO] tenant connection deleted (%s)", d.Id())
	d.SetId("")
	return nil
}

func resourceTenantConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading tenant connection (%s)", d.Id())

	tenantConnection, err := expandTenantConnectionFromID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := m.(*client.Client)
	tenant, err := client.Tenants.GetByID(tenantConnection.TenantID)
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "tenant")
	}

	tenantConnection, err = expandTenantConnectionFromTenant(tenant, tenantConnection.ProjectID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTenantConnection(ctx, d, tenantConnection); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] tenant connection read (%s)", d.Id())
	return nil
}

func resourceTenantConnectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating tenant connection (%s)", d.Id())

	tenantConnection := expandTenantConnection(d)
	client := m.(*client.Client)
	tenant, err := client.Tenants.GetByID(tenantConnection.TenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	tenant.ProjectEnvironments[tenantConnection.ProjectID] = tenantConnection.EnvironmentIDs
	tenant, err = client.Tenants.Update(tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	tenantConnection, err = expandTenantConnectionFromTenant(tenant, tenantConnection.ProjectID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTenantConnection(ctx, d, tenantConnection); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] tenant connection updated (%s)", d.Id())
	return nil
}
