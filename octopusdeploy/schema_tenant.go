package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandTenant(d *schema.ResourceData) *tenants.Tenant {
	name := d.Get("name").(string)

	tenant := tenants.NewTenant(name)
	tenant.ID = d.Id()

	if v, ok := d.GetOk("cloned_from_tenant_id"); ok {
		tenant.ClonedFromTenantID = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		tenant.Description = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		tenant.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		tenant.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return tenant
}

func flattenTenant(tenant *tenants.Tenant) map[string]interface{} {
	if tenant == nil {
		return nil
	}

	return map[string]interface{}{
		"cloned_from_tenant_id": tenant.ClonedFromTenantID,
		"description":           tenant.Description,
		"id":                    tenant.GetID(),
		"name":                  tenant.Name,
		"space_id":              tenant.SpaceID,
		"tenant_tags":           tenant.TenantTags,
	}
}

func getTenantDataSchema() map[string]*schema.Schema {
	dataSchema := getTenantSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"cloned_from_tenant_id": getQueryClonedFromTenantID(),
		"id":                    getDataSchemaID(),
		"ids":                   getQueryIDs(),
		"is_clone":              getQueryIsClone(),
		"name":                  getQueryName(),
		"partial_name":          getQueryPartialName(),
		"project_id":            getQueryProjectID(),
		"skip":                  getQuerySkip(),
		"tags":                  getQueryTags(),
		"space_id":              getQuerySpaceID(),
		"tenants": {
			Computed:    true,
			Description: "A list of tenants that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"take": getQueryTake(),
	}
}

func getTenantSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cloned_from_tenant_id": {
			Description: "The ID of the tenant from which this tenant was cloned.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"description": getDescriptionSchema("tenant"),
		"id":          getIDSchema(),
		"name":        getNameSchema(true),
		"space_id":    getSpaceIDSchema(),
		"tenant_tags": getTenantTagsSchema(),
	}
}

func setTenant(ctx context.Context, d *schema.ResourceData, tenant *tenants.Tenant) error {
	d.Set("cloned_from_tenant_id", tenant.ClonedFromTenantID)
	d.Set("description", tenant.Description)
	d.Set("id", tenant.GetID())
	d.Set("name", tenant.Name)
	d.Set("space_id", tenant.SpaceID)

	if err := d.Set("tenant_tags", tenant.TenantTags); err != nil {
		return fmt.Errorf("error setting tenant_tags: %s", err)
	}

	return nil
}
