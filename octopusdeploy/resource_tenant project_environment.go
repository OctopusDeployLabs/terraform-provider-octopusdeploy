package octopusdeploy

import (
	"context"
	"log"
	"slices"

	"net/http"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTenantProjectEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTenantProjectEnvironmentCreate,
		DeleteContext: resourceTenantProjectEnvironmentDelete,
		Description:   "This resource manages tenants in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceTenantProjectEnvironmentRead,
		Schema:        getTenantProjectEnvironmentSchema(),
	}
}

func resourceTenantProjectEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	client := m.(*client.Client)
	k := extractRelationship(d, client)

	log.Printf("[INFO] connecting tenant (%#v) to project (%#v) for environment (%#v)", k.tenantID, k.projectID, k.environmentID)

	tenant, err := tenants.GetByID(client, k.spaceID, k.tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Append relationship if not present
	if !slices.Contains(tenant.ProjectEnvironments[k.projectID], k.environmentID) {
		tenant.ProjectEnvironments[k.projectID] = append(tenant.ProjectEnvironments[k.projectID], k.environmentID)
	}

	_, err = tenants.Update(client, tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	id := k.spaceID + ":" + k.tenantID + ":" + k.projectID + ":" + k.environmentID
	d.SetId(id)

	log.Printf("[INFO] tenant (%s) connected to project (%#v) for environment (%#v)", k.tenantID, k.projectID, k.environmentID)
	return nil
}

func resourceTenantProjectEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	client := m.(*client.Client)
	k := extractRelationship(d, client)

	log.Printf("[INFO] removing tenant (%#v) from project (%#v) for environment (%#v)", k.tenantID, k.projectID, k.environmentID)

	tenant, err := tenants.GetByID(client, k.spaceID, k.tenantID)
	if err != nil {
		apiError := err.(*core.APIError)
		if apiError.StatusCode == http.StatusNotFound {
			log.Printf("[INFO] tenant (%#v) no longer exists", k.tenantID)
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	p := tenant.ProjectEnvironments[k.projectID]
	if slices.Contains(p, k.environmentID) {
		for i := 0; i < len(p); i++ {
			if p[i] == k.environmentID {
				tenant.ProjectEnvironments[k.projectID] = slices.Delete(p, i, i+1)
			}
		}
	}

	_, err = tenants.Update(client, tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] tenant (%#v) disconnected from project (%#v) for environment (%#v)", k.tenantID, k.projectID, k.environmentID)
	d.SetId("")
	return nil
}

func resourceTenantProjectEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	k := extractRelationship(d, client)
	_, err := tenants.GetByID(client, k.spaceID, k.tenantID)
	if err != nil {
		apiError := err.(*core.APIError)
		if apiError.StatusCode != http.StatusNotFound {
			return diag.FromErr(err)
		}
	}
	return nil
}

func extractRelationship(d *schema.ResourceData, client *client.Client) person {
	tenantID := d.Get("tenant_id").(string)
	projectID := d.Get("project_id").(string)
	environmentID := d.Get("environment_id").(string)

	spaceID := client.GetSpaceID()
	if v, ok := d.GetOk("space_id"); ok {
		spaceID = v.(string)
	}

	n := person{tenantID: tenantID, projectID: projectID, environmentID: environmentID, spaceID: spaceID}
	return n
}

type person struct {
	tenantID      string
	projectID     string
	environmentID string
	spaceID       string
}
