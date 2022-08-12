package octopusdeploy

import (
	"context"
	"log"
	"net/http"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceScopedUserRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScopedUserRoleCreate,
		DeleteContext: resourceScopedUserRoleDelete,
		Description:   "This resource manages scoped user roles in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceScopedUserRoleRead,
		Schema:        getScopedUserRoleSchema(),
		UpdateContext: resourceScopedUserRoleUpdate,
	}
}

func resourceScopedUserRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	scopedUserRole := expandScopedUserRole(d)

	log.Printf("[INFO] creating scoped user role: %#v", scopedUserRole)

	client := m.(*client.Client)
	createdScopedUserRole, err := client.ScopedUserRoles.Add(scopedUserRole)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setScopedUserRole(ctx, d, createdScopedUserRole); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdScopedUserRole.GetID())

	log.Printf("[INFO] scoped user role created (%s)", d.Id())
	return resourceScopedUserRoleRead(ctx, d, m)
}

func resourceScopedUserRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting scoped user role (%s)", d.Id())

	client := m.(*client.Client)
	if err := client.ScopedUserRoles.DeleteByID(d.Id()); err != nil {
		apiError := err.(*core.APIError)
		if apiError.StatusCode != http.StatusNotFound {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	log.Printf("[INFO] scoped user role deleted")
	return nil
}

func resourceScopedUserRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading scoped user role (%s)", d.Id())

	client := m.(*client.Client)
	scopedUserRole, err := client.ScopedUserRoles.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "scoped user role")
	}

	if err := setScopedUserRole(ctx, d, scopedUserRole); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] scoped user role read (%s)", d.Id())
	return nil
}

func resourceScopedUserRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating scoped user role (%s)", d.Id())

	scopedUserRole := expandScopedUserRole(d)
	client := m.(*client.Client)
	updatedScopedUserRole, err := client.ScopedUserRoles.Update(scopedUserRole)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setScopedUserRole(ctx, d, updatedScopedUserRole); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] scoped user role updated (%s)", d.Id())
	return resourceScopedUserRoleRead(ctx, d, m)
}
