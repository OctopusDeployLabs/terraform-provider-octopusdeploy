package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzureServicePrincipalAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureServicePrincipalAccountCreate,
		DeleteContext: resourceAzureServicePrincipalAccountDelete,
		Importer:      getImporter(),
		ReadContext:   resourceAzureServicePrincipalAccountRead,
		Schema:        getAzureServicePrincipalAccountSchema(),
		UpdateContext: resourceAzureServicePrincipalAccountUpdate,
	}
}

func resourceAzureServicePrincipalAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAzureServicePrincipalAccount(d)

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())
	return resourceAzureServicePrincipalAccountRead(ctx, d, m)
}

func resourceAzureServicePrincipalAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Accounts.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceAzureServicePrincipalAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	azureSubscriptionAccount := accountResource.(*octopusdeploy.AzureServicePrincipalAccount)

	setAzureServicePrincipalAccount(ctx, d, azureSubscriptionAccount)
	return nil
}

func resourceAzureServicePrincipalAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAzureServicePrincipalAccount(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAzureServicePrincipalAccountRead(ctx, d, m)
}
