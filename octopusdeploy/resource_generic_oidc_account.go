package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceGenericOpenIDConnectAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGenericOpenIDConnectAccountCreate,
		DeleteContext: resourceGenericOpenIDConnectAccountDelete,
		Description:   "This resource manages Generic OpenID Connect accounts in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceGenericOpenIDConnectAccountRead,
		Schema:        getGenericOpenIdConnectAccountSchema(),
		UpdateContext: resourceGenericOpenIDConnectAccountUpdate,
	}
}

func resourceGenericOpenIDConnectAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandGenericOpenIDConnectAccount(d)

	log.Printf("[INFO] creating Generic OpenID Connect account: %#v", account)

	client := m.(*client.Client)
	createdAccount, err := accounts.Add(client, account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setGenericOpenIDConnectAccount(ctx, d, createdAccount.(*accounts.GenericOIDCAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] Generic OpenID Connect account created (%s)", d.Id())
	return nil
}

func resourceGenericOpenIDConnectAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Generic OpenID Connect account (%s)", d.Id())

	client := m.(*client.Client)
	if err := accounts.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] Generic OpenID Connect account deleted")
	return nil
}

func resourceGenericOpenIDConnectAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading Generic OpenID Connect account (%s)", d.Id())

	client := m.(*client.Client)
	accountResource, err := accounts.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Generic OpenID Connect account")
	}

	genericOIDCAccount := accountResource.(*accounts.GenericOIDCAccount)
	if err := setGenericOpenIDConnectAccount(ctx, d, genericOIDCAccount); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Generic OpenID Connect account read (%s)", d.Id())
	return nil
}

func resourceGenericOpenIDConnectAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandGenericOpenIDConnectAccount(d)

	log.Printf("[INFO] updating Generic OpenID Connect account %#v", account)

	client := m.(*client.Client)
	updatedAccount, err := accounts.Update(client, account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setGenericOpenIDConnectAccount(ctx, d, updatedAccount.(*accounts.GenericOIDCAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Generic OpenID Connect account updated (%s)", d.Id())
	return nil
}
