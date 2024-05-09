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

func resourceAmazonWebServicesOpenIDConnectAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAmazonWebServicesOpenIDConnectAccountCreate,
		DeleteContext: resourceAmazonWebServicesOpenIDConnectAccountDelete,
		Description:   "This resource manages AWS OIDC accounts in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceAmazonWebServicesOpenIDConnectAccountRead,
		Schema:        getAmazonWebServicesOpenIDConnectAccountSchema(),
		UpdateContext: resourceAmazonWebServicesOpenIDConnectAccountUpdate,
	}
}

func resourceAmazonWebServicesOpenIDConnectAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAmazonWebServicesOpenIDConnectAccount(d)

	log.Printf("[INFO] creating AWS OIDC account")

	client := m.(*client.Client)
	createdAccount, err := accounts.Add(client, account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAmazonWebServicesOpenIDConnectAccount(ctx, d, createdAccount.(*accounts.AwsOIDCAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] AWS OIDC account created (%s)", d.Id())
	return nil
}

func resourceAmazonWebServicesOpenIDConnectAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting AWS OIDC account (%s)", d.Id())

	client := m.(*client.Client)
	if err := client.Accounts.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] AWS OIDC account deleted")
	return nil
}

func resourceAmazonWebServicesOpenIDConnectAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading AWS OIDC account (%s)", d.Id())

	client := m.(*client.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "AWS OIDC account")
	}

	awsOIDCAccount := accountResource.(*accounts.AwsOIDCAccount)
	if err := setAmazonWebServicesOpenIDConnectAccount(ctx, d, awsOIDCAccount); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] AWS OIDC account read: %#v", awsOIDCAccount)
	return nil
}

func resourceAmazonWebServicesOpenIDConnectAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAmazonWebServicesOpenIDConnectAccount(d)

	log.Printf("[INFO] updating AWS OIDC account: %#v", account)

	client := m.(*client.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAmazonWebServicesOpenIDConnectAccount(ctx, d, updatedAccount.(*accounts.AwsOIDCAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] AWS OIDC account updated (%s)", d.Id())
	return nil
}
