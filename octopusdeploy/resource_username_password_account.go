package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUsernamePasswordAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUsernamePasswordAccountCreate,
		DeleteContext: resourceUsernamePasswordAccountDelete,
		Importer:      getImporter(),
		ReadContext:   resourceUsernamePasswordAccountRead,
		Schema:        getUsernamePasswordAccountSchema(),
		UpdateContext: resourceUsernamePasswordAccountUpdate,
	}
}

func resourceUsernamePasswordAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandUsernamePasswordAccount(d)

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		diag.FromErr(err)
	}

	createdUsernamePasswordAccount := createdAccount.(*octopusdeploy.UsernamePasswordAccount)

	setUsernamePasswordAccount(ctx, d, createdUsernamePasswordAccount)
	return nil
}

func resourceUsernamePasswordAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Accounts.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceUsernamePasswordAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if accountResource == nil {
		d.SetId("")
		return nil
	}

	accountResource, err = octopusdeploy.ToAccount(accountResource.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	usernamePasswordAccount := accountResource.(*octopusdeploy.UsernamePasswordAccount)

	setUsernamePasswordAccount(ctx, d, usernamePasswordAccount)
	return nil
}

func resourceUsernamePasswordAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandUsernamePasswordAccount(d)

	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.Update(account)
	if err != nil {
		diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccount(accountResource.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	updatedUsernamePasswordAccount := accountResource.(*octopusdeploy.UsernamePasswordAccount)

	setUsernamePasswordAccount(ctx, d, updatedUsernamePasswordAccount)
	return nil
}
