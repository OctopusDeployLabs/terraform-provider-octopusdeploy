package octopusdeploy

import (
	"context"
	"log"

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

	log.Printf("[INFO] creating account: %#v", accountResource)

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(accountResource)
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccountResource(createdAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAccountResource(ctx, d, accountResource); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(accountResource.GetID())

	log.Printf("[INFO] account created (%s)", d.Id())
	return nil
}

func resourceAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting account (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Accounts.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] account deleted")
	return nil
}

func resourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading account (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			log.Printf("[INFO] account (%s) not found; deleting from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	accountResource := account.(*octopusdeploy.AccountResource)

	if err := setAccountResource(ctx, d, accountResource); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] account read (%s)", d.Id())
	return nil
}

func resourceAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating account (%s)", d.Id())

	accountResource := expandAccountResource(d)
	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(accountResource)
	if err != nil {
		return diag.FromErr(err)
	}

	updatedAccountResource := updatedAccount.(*octopusdeploy.AccountResource)

	if err := setAccountResource(ctx, d, updatedAccountResource); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] account updated (%s)", d.Id())
	return nil
}
