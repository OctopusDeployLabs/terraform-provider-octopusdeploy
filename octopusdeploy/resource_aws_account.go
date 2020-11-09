package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAmazonWebServicesAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAmazonWebServicesAccountCreate,
		DeleteContext: resourceAccountDeleteCommon,
		ReadContext:   resourceAmazonWebServicesAccountRead,
		Schema:        getAmazonWebServicesAccountSchema(),
		UpdateContext: resourceAmazonWebServicesAccountUpdate,
	}
}

func resourceAmazonWebServicesAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAmazonWebServicesAccount(d)

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	createdAmazonWebServicesAccount := createdAccount.(*octopusdeploy.AmazonWebServicesAccount)

	flattenAmazonWebServicesAccount(ctx, d, createdAmazonWebServicesAccount)
	return nil
}

func resourceAmazonWebServicesAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccount(accountResource.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	amazonWebServicesAccount := accountResource.(*octopusdeploy.AmazonWebServicesAccount)

	flattenAmazonWebServicesAccount(ctx, d, amazonWebServicesAccount)
	return nil
}

func resourceAmazonWebServicesAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAmazonWebServicesAccount(d)

	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	updatedAmazonWebServicesAccount := updatedAccount.(*octopusdeploy.AmazonWebServicesAccount)

	flattenAmazonWebServicesAccount(ctx, d, updatedAmazonWebServicesAccount)
	return nil
}
