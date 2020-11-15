package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountCreate,
		DeleteContext: resourceAccountDeleteCommon,
		Importer:      getImporter(),
		ReadContext:   resourceAccountRead,
		Schema:        getAccountResourceSchema(),
		UpdateContext: resourceAccountUpdate,
	}
}

func resourceAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountResource := expandAccountResource(d)

	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.Add(accountResource)
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccountResource(account)
	if err != nil {
		return diag.FromErr(err)
	}

	setAccountResource(ctx, d, accountResource)
	return nil
}

func resourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if account == nil {
		d.SetId("")
		return nil
	}

	accountResource := account.(*octopusdeploy.AccountResource)

	setAccountResource(ctx, d, accountResource)
	return nil
}

func resourceAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountResource := expandAccountResource(d)

	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.Update(accountResource)
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccountResource(account)
	if err != nil {
		return diag.FromErr(err)
	}

	setAccountResource(ctx, d, accountResource)
	return nil
}
