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

func resourceAzureSubscriptionAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureSubscriptionAccountCreate,
		DeleteContext: resourceAzureSubscriptionAccountDelete,
		Description:   "This resource manages Azure subscription accounts in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceAzureSubscriptionAccountRead,
		Schema:        getAzureSubscriptionAccountSchema(),
		UpdateContext: resourceAzureSubscriptionAccountUpdate,
	}
}

func resourceAzureSubscriptionAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAzureSubscriptionAccount(d)

	log.Printf("[INFO] creating Azure subscription account: %#v", account)

	client := m.(*client.Client)
	createdAccount, err := accounts.Add(client, account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureSubscriptionAccount(ctx, d, createdAccount.(*accounts.AzureSubscriptionAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] Azure subscription account created (%s)", d.Id())
	return nil
}

func resourceAzureSubscriptionAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Azure subscription account (%s)", d.Id())

	client := m.(*client.Client)
	if err := accounts.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] Azure subscription account deleted")
	return nil
}

func resourceAzureSubscriptionAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading Azure subscription account (%s)", d.Id())

	client := m.(*client.Client)
	accountResource, err := accounts.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Azure subscription account")
	}

	azureSubscriptionAccount := accountResource.(*accounts.AzureSubscriptionAccount)
	if err := setAzureSubscriptionAccount(ctx, d, azureSubscriptionAccount); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure subscription account read (%s)", d.Id())
	return nil
}

func resourceAzureSubscriptionAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAzureSubscriptionAccount(d)

	log.Printf("[INFO] updating Azure subscription account %#v", account)

	client := m.(*client.Client)
	updatedAccount, err := accounts.Update(client, account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureSubscriptionAccount(ctx, d, updatedAccount.(*accounts.AzureSubscriptionAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure subscription account updated (%s)", d.Id())
	return nil
}
