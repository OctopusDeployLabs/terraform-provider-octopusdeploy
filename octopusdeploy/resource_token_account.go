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
		DeleteContext: resourceAccountDeleteCommon,
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

	createdTokenAccount := createdAccount.(*octopusdeploy.TokenAccount)
	createdTokenAccount.Token = account.Token

	flattenTokenAccount(ctx, d, createdTokenAccount)
	return nil
}

func resourceTokenAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	tokenAccount := accountResource.(*octopusdeploy.TokenAccount)

	flattenTokenAccount(ctx, d, tokenAccount)
	return nil
}

func resourceTokenAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandTokenAccount(d)

	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccount(accountResource.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	updatedTokenAccount := accountResource.(*octopusdeploy.TokenAccount)
	updatedTokenAccount.Token = account.Token

	flattenTokenAccount(ctx, d, updatedTokenAccount)
	return nil
}
