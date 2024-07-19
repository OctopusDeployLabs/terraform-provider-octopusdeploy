package octopusdeploy

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTenantProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTenantProjectCreate,
		DeleteContext: resourceTenantProjectDelete,
		Description:   "This resource represents the connection between tenants and projects.",
		Importer:      getImporter(),
		ReadContext:   resourceTenantProjectRead,
		UpdateContext: resourceTenantProjectUpdate,
		Schema:        getTenantProjectSchema(),
	}
}

func resourceTenantProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	client := m.(*client.Client)
	k := extractRelationship(d, client)

	log.Printf("[INFO] updating tenant (%#v) connection to project (%#v)", k.tenantID, k.projectID)

	tenant, err := tenants.GetByID(client, k.spaceID, k.tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	tenant.ProjectEnvironments[k.projectID] = k.environmentIDs

	_, err = tenants.Update(client, tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] updated tenant (%s) connection to project (%#v)", k.tenantID, k.projectID)
	return nil
}

func resourceTenantProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	client := m.(*client.Client)
	k := extractRelationship(d, client)

	log.Printf("[INFO] connecting tenant (%#v) to project (%#v)", k.tenantID, k.projectID)

	tenant, err := tenants.GetByID(client, k.spaceID, k.tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	tenant.ProjectEnvironments[k.projectID] = k.environmentIDs

	_, err = tenants.Update(client, tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	id := k.spaceID + ":" + k.tenantID + ":" + k.projectID
	d.SetId(id)

	log.Printf("[INFO] tenant (%s) connected to project (%#v)", k.tenantID, k.projectID)
	return nil
}

func resourceTenantProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	client := m.(*client.Client)
	k := extractRelationship(d, client)

	log.Printf("[INFO] removing tenant (%#v) from project (%#v)", k.tenantID, k.projectID)

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

	delete(tenant.ProjectEnvironments, k.projectID)

	_, err = tenants.Update(client, tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] tenant (%#v) disconnected from project (%#v)", k.tenantID, k.projectID)
	d.SetId("")
	return nil
}

func resourceTenantProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	bits := strings.Split(d.Id(), ":")
	spaceID := bits[0]
	tenantID := bits[1]
	projectID := bits[2]

	tenant, err := tenants.GetByID(client, spaceID, tenantID)
	if err != nil {
		apiError := err.(*core.APIError)
		if apiError.StatusCode != http.StatusNotFound {
			return diag.FromErr(err)
		}
	}

	d.Set("environment_ids", tenant.ProjectEnvironments[projectID])

	return nil
}

func extractRelationship(d *schema.ResourceData, client *client.Client) person {
	tenantID := d.Get("tenant_id").(string)
	projectID := d.Get("project_id").(string)

	environmentIDs := []string{}
	if attr, ok := d.GetOk("environment_ids"); ok {
		environmentIDs = getSliceFromTerraformTypeList(attr)
	}

	spaceID := client.GetSpaceID()
	if v, ok := d.GetOk("space_id"); ok {
		spaceID = v.(string)
	}

	n := person{tenantID: tenantID, projectID: projectID, environmentIDs: environmentIDs, spaceID: spaceID}
	return n
}

type person struct {
	tenantID       string
	projectID      string
	environmentIDs []string
	spaceID        string
}
