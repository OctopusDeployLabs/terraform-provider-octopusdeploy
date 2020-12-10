package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext:      resourceAccountCreate,
		DeleteContext:      resourceAccountDelete,
		DeprecationMessage: "Use an account-specific resource instead (i.e. octopusdeploy_aws_account, octopusdeploy_azure_service_principal, octopusdeploy_azure_subscription_account, octopusdeploy_ssh_key_account, octopusdeploy_token_account, octopusdeploy_username_password_account).",
		Description:        "This resource manages accounts in Octopus Deploy.",
		Importer:           getImporter(),
		ReadContext:        resourceAccountRead,
		Schema:             getAccountResourceSchema(),
		UpdateContext:      resourceAccountUpdate,
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

	d.SetId(accountResource.GetID())
	return resourceAccountRead(ctx, d, m)
}

func resourceAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Accounts.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	accountResource := account.(*octopusdeploy.AccountResource)

	setAccountResource(ctx, d, accountResource)
	return nil
}

func resourceAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountResource := expandAccountResource(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Accounts.Update(accountResource)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAccountRead(ctx, d, m)
}
