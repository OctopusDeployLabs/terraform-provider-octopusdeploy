package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureSubscriptionAccount(ctx, d, createdAccount.(*octopusdeploy.AzureSubscriptionAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] Azure subscription account created (%s)", d.Id())
	return nil
}

func resourceAzureSubscriptionAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Azure subscription account (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Accounts.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] Azure subscription account deleted")
	return nil
}

func resourceAzureSubscriptionAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading Azure subscription account (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			log.Printf("[INFO] Azure subscription account (%s) not found; deleting from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccount(accountResource.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	azureSubscriptionAccount := accountResource.(*octopusdeploy.AzureSubscriptionAccount)
	if err := setAzureSubscriptionAccount(ctx, d, azureSubscriptionAccount); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure subscription account read (%s)", d.Id())
	return nil
}

func resourceAzureSubscriptionAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAzureSubscriptionAccount(d)

	log.Printf("[INFO] updating Azure subscription account %#v", account)

	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err := octopusdeploy.ToAccount(updatedAccount.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureSubscriptionAccount(ctx, d, accountResource.(*octopusdeploy.AzureSubscriptionAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure subscription account updated (%s)", d.Id())
	return nil
}
