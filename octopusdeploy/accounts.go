package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func fetchAndReadAccount(ctx context.Context, d *schema.ResourceData, m interface{}) (octopusdeploy.IAccount, diag.Diagnostics) {
	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return nil, diag.FromErr(err)
	}

	accountResource := account.(*octopusdeploy.AccountResource)

	setAccountResource(ctx, d, accountResource)
	return accountResource, nil
}

func resourceAccountCreateCommon(ctx context.Context, d *schema.ResourceData, m interface{}, account octopusdeploy.IAccount) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	createdAccountResource := createdAccount.(*octopusdeploy.AccountResource)

	setAccountResource(ctx, d, createdAccountResource)
	return nil
}

func resourceAccountUpdateCommon(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountResource := expandAccountResource(d)

	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(accountResource)
	if err != nil {
		return diag.FromErr(err)
	}

	updatedAccountResource := updatedAccount.(*octopusdeploy.AccountResource)

	setAccountResource(ctx, d, updatedAccountResource)
	return nil
}

func resourceAccountDeleteCommon(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Accounts.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
