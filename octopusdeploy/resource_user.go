package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	resourceUserImporter := &schema.ResourceImporter{
		StateContext: schema.ImportStatePassthroughContext,
	}

	return &schema.Resource{
		CreateContext: resourceUserCreate,
		DeleteContext: resourceUserDelete,
		Importer:      resourceUserImporter,
		ReadContext:   resourceUserRead,
		Schema:        getUserSchema(),
		UpdateContext: resourceUserUpdate,
	}
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	user, err := client.Users.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	flattenUser(ctx, d, user)
	return nil
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	user := expandUser(d)

	client := m.(*octopusdeploy.Client)
	createdUser, err := client.Users.Add(user)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenUser(ctx, d, createdUser)
	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	user := expandUser(d)

	client := m.(*octopusdeploy.Client)
	updatedUser, err := client.Users.Update(user)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenUser(ctx, d, updatedUser)
	return nil
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
