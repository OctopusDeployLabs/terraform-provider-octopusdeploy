package octopusdeploy

import (
	"context"
	"fmt"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandTenantConnectionFromTenant(tenant *tenants.Tenant, projectID string) (*TenantConnection, error) {
	environmentIDs, ok := tenant.ProjectEnvironments[projectID]
	if !ok {
		return nil, nil
	}

	tc := TenantConnection{
		TenantID:       tenant.ID,
		ProjectID:      projectID,
		EnvironmentIDs: environmentIDs,
	}

	tc.ID = tc.GetID()

	return &tc, nil
}

func expandTenantConnectionFromID(id string) (*TenantConnection, error) {
	errInvalidID := fmt.Errorf(
		"expected id to be in format '%s', got: %s",
		"tenant_id:project_id[:environment_id_0[+environment_id_1[+environment_id_n]]]",
		id,
	)

	parts := strings.Split(id, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return nil, errInvalidID
	}

	tc := TenantConnection{
		TenantID:       parts[0],
		ProjectID:      parts[1],
		EnvironmentIDs: []string{},
	}

	if len(parts) == 3 {
		tc.EnvironmentIDs = strings.Split(parts[2], "+")
	}

	tc.ID = tc.GetID()

	return &tc, nil
}

type TenantConnection struct {
	ID             string   `json:"id"`
	TenantID       string   `json:"tenant_id"`
	ProjectID      string   `json:"project_id"`
	EnvironmentIDs []string `json:"environment_ids"`
}

func (tc *TenantConnection) GetID() string {
	id := fmt.Sprintf("%s:%s", tc.TenantID, tc.ProjectID)

	if len(tc.EnvironmentIDs) > 0 {
		id = fmt.Sprintf("%s:%s", id, strings.Join(tc.EnvironmentIDs, "+"))
	}

	return id
}

func expandTenantConnection(d *schema.ResourceData) *TenantConnection {
	tenantConnection := &TenantConnection{}
	tenantConnection.ID = d.Id()

	if v, ok := d.GetOk("tenant_id"); ok {
		tenantConnection.TenantID = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		tenantConnection.ProjectID = v.(string)
	}

	if v, ok := d.GetOk("environment_ids"); ok {
		list := v.(*schema.Set).List()
		tenantConnection.EnvironmentIDs = expandArray(list)
	}

	return tenantConnection
}

func flattenTenantConnection(tenantConnection *TenantConnection) map[string]interface{} {
	if tenantConnection == nil {
		return nil
	}

	return map[string]interface{}{
		"id":              tenantConnection.GetID(),
		"tenant_id":       tenantConnection.TenantID,
		"project_id":      tenantConnection.ProjectID,
		"environment_ids": tenantConnection.EnvironmentIDs,
	}
}

func getTenantConnectionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tenant_id": {
			Description: "The ID of the tenant that this connection belongs to",
			Required:    true,
			Type:        schema.TypeString,
		},
		"project_id": {
			Description: "The ID of project that is being connected to the tenant",
			Required:    true,
			Type:        schema.TypeString,
		},
		"environment_ids": {
			Description: "The list of environment IDs for which the project will deploy to",
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}

func setTenantConnection(ctx context.Context, d *schema.ResourceData, tenantConnection *TenantConnection) error {
	d.Set("id", tenantConnection.GetID())
	d.Set("tenant_id", tenantConnection.TenantID)
	d.Set("project_id", tenantConnection.ProjectID)

	if tenantConnection.EnvironmentIDs == nil {
		d.Set("environment_ids", []string{})
	} else {
		d.Set("environment_ids", tenantConnection.EnvironmentIDs)
	}

	return nil
}
