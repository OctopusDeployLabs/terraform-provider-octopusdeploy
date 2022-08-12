package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		DeleteContext: resourceUserDelete,
		Description:   "This resource manages users in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceUserRead,
		Schema:        getUserSchema(),
		UpdateContext: resourceUserUpdate,
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	user := expandUser(d)

	log.Printf("[DEBUG] creating user")

	client := m.(*client.Client)
	createdUser, err := client.Users.Add(user)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setUser(ctx, d, createdUser); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdUser.GetID())

	log.Printf("[DEBUG] user created (%s)", d.Id())
	return nil
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting user (%s)", d.Id())

	client := m.(*client.Client)
	if err := client.Users.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] user deleted")
	return nil
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading user (%s)", d.Id())

	client := m.(*client.Client)
	user, err := client.Users.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "user")
	}

	if err := setUser(ctx, d, user); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] user read (%s)", d.Id())
	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating user (%s)", d.Id())

	user := expandUser(d)
	client := m.(*client.Client)
	updatedUser, err := client.Users.Update(user)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setUser(ctx, d, updatedUser); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] user updated (%s)", d.Id())
	return nil
}
