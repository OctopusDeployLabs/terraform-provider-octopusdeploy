package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		DeleteContext: resourceUserDelete,
		Importer:      getImporter(),
		ReadContext:   resourceUserRead,
		Schema:        getUserSchema(),
		UpdateContext: resourceUserUpdate,
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	user := expandUser(d)

	client := m.(*octopusdeploy.Client)
	createdUser, err := client.Users.Add(user)
	if createdUser != nil && err == nil {
		d.SetId(createdUser.GetID())
		return nil
	}

	return diag.FromErr(err)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Users.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	user, err := client.Users.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	setUser(ctx, d, user)
	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	user := expandUser(d)

	client := m.(*octopusdeploy.Client)
	updatedUser, err := client.Users.Update(user)
	if updatedUser != nil && err == nil {
		d.SetId(updatedUser.GetID())
		return nil
	}

	return diag.FromErr(err)
}
