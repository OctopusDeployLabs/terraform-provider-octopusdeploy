package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzureServicePrincipalAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureServicePrincipalAccountCreate,
		DeleteContext: resourceAzureServicePrincipalAccountDelete,
		Description:   "This resource manages Azure service principal accounts in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceAzureServicePrincipalAccountRead,
		Schema:        getAzureServicePrincipalAccountSchema(),
		UpdateContext: resourceAzureServicePrincipalAccountUpdate,
	}
}

func resourceAzureServicePrincipalAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAzureServicePrincipalAccount(d)

	log.Printf("[INFO] creating Azure service principal account: %#v", account)

	client := m.(*client.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureServicePrincipalAccount(ctx, d, createdAccount.(*accounts.AzureServicePrincipalAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] Azure service principal account created (%s)", d.Id())
	return nil
}

func resourceAzureServicePrincipalAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Azure service principal account (%s)", d.Id())

	client := m.(*client.Client)
	if err := client.Accounts.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] Azure service principal account deleted")
	return nil
}

func resourceAzureServicePrincipalAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading Azure service principal account (%s)", d.Id())

	client := m.(*client.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Azure service principal account")
	}

	azureServicePrincipalAccount := accountResource.(*accounts.AzureServicePrincipalAccount)
	if err := setAzureServicePrincipalAccount(ctx, d, azureServicePrincipalAccount); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure service principal account read (%s)", d.Id())
	return nil
}

func resourceAzureServicePrincipalAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAzureServicePrincipalAccount(d)

	log.Printf("[INFO] updating Azure service principal account %#v", account)

	client := m.(*client.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureServicePrincipalAccount(ctx, d, updatedAccount.(*accounts.AzureServicePrincipalAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure service principal account updated (%s)", d.Id())
	return nil
}
