package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTokenAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTokenAccountCreate,
		DeleteContext: resourceTokenAccountDelete,
		Importer:      getImporter(),
		ReadContext:   resourceTokenAccountRead,
		Schema:        getTokenAccountSchema(),
		UpdateContext: resourceTokenAccountUpdate,
	}
}

func resourceTokenAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandTokenAccount(d)

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())
	return resourceTokenAccountRead(ctx, d, m)
}

func resourceTokenAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Accounts.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceTokenAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccount(accountResource.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	tokenAccount := accountResource.(*octopusdeploy.TokenAccount)

	setTokenAccount(ctx, d, tokenAccount)
	return nil
}

func resourceTokenAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandTokenAccount(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTokenAccountRead(ctx, d, m)
}
