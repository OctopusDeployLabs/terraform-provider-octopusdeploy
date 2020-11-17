package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzureSubscriptionAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureSubscriptionAccountCreate,
		DeleteContext: resourceAzureSubscriptionAccountDelete,
		Importer:      getImporter(),
		ReadContext:   resourceAzureSubscriptionAccountRead,
		Schema:        getAzureSubscriptionAccountSchema(),
		UpdateContext: resourceAzureSubscriptionAccountUpdate,
	}
}

func resourceAzureSubscriptionAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAzureSubscriptionAccount(d)

	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	createdAzureSubscriptionAccount := accountResource.(*octopusdeploy.AzureSubscriptionAccount)

	setAzureSubscriptionAccount(ctx, d, createdAzureSubscriptionAccount)
	return nil
}

func resourceAzureSubscriptionAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Accounts.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceAzureSubscriptionAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	azureSubscriptionAccount := accountResource.(*octopusdeploy.AzureSubscriptionAccount)

	setAzureSubscriptionAccount(ctx, d, azureSubscriptionAccount)
	return nil
}

func resourceAzureSubscriptionAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAzureSubscriptionAccount(d)

	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccount(accountResource.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	updatedAzureSubscriptionAccount := accountResource.(*octopusdeploy.AzureSubscriptionAccount)

	setAzureSubscriptionAccount(ctx, d, updatedAzureSubscriptionAccount)
	return nil
}
