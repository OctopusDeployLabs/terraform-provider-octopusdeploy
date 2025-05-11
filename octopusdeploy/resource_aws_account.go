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
	account, err := expandAmazonWebServicesAccount(d)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] creating AWS account")

	client := m.(*client.Client)
	createdAccount, err := accounts.Add(client, account)

	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAmazonWebServicesAccount(ctx, d, createdAccount.(*accounts.AmazonWebServicesAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] AWS account created (%s)", d.Id())
	return nil
}

func resourceAmazonWebServicesAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting AWS account (%s)", d.Id())

	client := m.(*client.Client)
	if err := client.Accounts.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] AWS account deleted")
	return nil
}

func resourceAmazonWebServicesAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading AWS account (%s)", d.Id())

	client := m.(*client.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "AWS account")
	}

	amazonWebServicesAccount := accountResource.(*accounts.AmazonWebServicesAccount)
	if err := setAmazonWebServicesAccount(ctx, d, amazonWebServicesAccount); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] AWS account read: %#v", amazonWebServicesAccount)
	return nil
}

func resourceAmazonWebServicesAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, err := expandAmazonWebServicesAccount(d)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] updating AWS account: %#v", account)

	client := m.(*client.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAmazonWebServicesAccount(ctx, d, updatedAccount.(*accounts.AmazonWebServicesAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] AWS account updated (%s)", d.Id())
	return nil
}
