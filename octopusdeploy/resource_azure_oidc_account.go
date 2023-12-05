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

func resourceAzureOpenIDConnectAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureOpenIDConnectAccountCreate,
		DeleteContext: resourceAzureOpenIDConnectAccountDelete,
		Description:   "This resource manages Azure OpenID Connect accounts in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceAzureOpenIDConnectAccountRead,
		Schema:        getAzureOpenIdConnectAccountSchema(),
		UpdateContext: resourceAzureOpenIDConnectAccountUpdate,
	}
}

func resourceAzureOpenIDConnectAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAzureOpenIDConnectAccount(d)

	log.Printf("[INFO] creating Azure OpenID Connect account: %#v", account)

	client := m.(*client.Client)
	createdAccount, err := accounts.Add(client, account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureOpenIDConnectAccount(ctx, d, createdAccount.(*accounts.AzureOIDCAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] Azure OpenID Connect account created (%s)", d.Id())
	return nil
}

func resourceAzureOpenIDConnectAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Azure OpenID Connect account (%s)", d.Id())

	client := m.(*client.Client)
	if err := accounts.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] Azure OpenID Connect account deleted")
	return nil
}

func resourceAzureOpenIDConnectAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading Azure OpenID Connect account (%s)", d.Id())

	client := m.(*client.Client)
	accountResource, err := accounts.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Azure OpenID Connect account")
	}

	azureOIDCAccount := accountResource.(*accounts.AzureOIDCAccount)
	if err := setAzureOpenIDConnectAccount(ctx, d, azureOIDCAccount); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure OpenID Connect account read (%s)", d.Id())
	return nil
}

func resourceAzureOpenIDConnectAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandAzureOpenIDConnectAccount(d)

	log.Printf("[INFO] updating Azure OpenID Connect account %#v", account)

	client := m.(*client.Client)
	updatedAccount, err := accounts.Update(client, account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureOpenIDConnectAccount(ctx, d, updatedAccount.(*accounts.AzureOIDCAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure OpenID Connect account updated (%s)", d.Id())
	return nil
}
