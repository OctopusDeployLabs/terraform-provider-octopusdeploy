package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAmazonWebServicesAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAmazonWebServicesAccountCreate,
		DeleteContext: resourceAmazonWebServicesAccountDelete,
		Description:   "This resource manages AWS accounts in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceAmazonWebServicesAccountRead,
		Schema:        getAmazonWebServicesAccountSchema(),
		UpdateContext: resourceAmazonWebServicesAccountUpdate,
	}
}

func resourceAmazonWebServicesAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAmazonWebServicesAccount(d)

	log.Printf("[INFO] creating AWS account: %#v", account)

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAmazonWebServicesAccount(ctx, d, createdAccount.(*octopusdeploy.AmazonWebServicesAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] AWS account created (%s)", d.Id())
	return nil
}

func resourceAmazonWebServicesAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting AWS account (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Accounts.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] AWS account deleted")
	return nil
}

func resourceAmazonWebServicesAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading AWS account (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			log.Printf("[INFO] AWS account (%s) not found; deleting from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccount(accountResource.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	amazonWebServicesAccount := accountResource.(*octopusdeploy.AmazonWebServicesAccount)

	if err := setAmazonWebServicesAccount(ctx, d, amazonWebServicesAccount); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] AWS account read: %#v", amazonWebServicesAccount)
	return nil
}

func resourceAmazonWebServicesAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAmazonWebServicesAccount(d)

	log.Printf("[INFO] updating AWS account: %#v", account)

	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err := octopusdeploy.ToAccount(updatedAccount.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAmazonWebServicesAccount(ctx, d, accountResource.(*octopusdeploy.AmazonWebServicesAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] AWS account updated (%s)", d.Id())
	return nil
}
