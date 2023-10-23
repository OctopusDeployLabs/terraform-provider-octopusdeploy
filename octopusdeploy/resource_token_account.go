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

func resourceTokenAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTokenAccountCreate,
		DeleteContext: resourceTokenAccountDelete,
		Importer:      getImporter(),
		ReadContext:   resourceTokenAccountRead,
		Schema:        getTokenAccountSchema(),
		UpdateContext: resourceTokenAccountUpdate,
	}
}

func resourceTokenAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandTokenAccount(d)

	log.Printf("[INFO] creating token account: %#v", account)

	client := m.(*client.Client)
	createdAccount, err := accounts.Add(client, account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTokenAccount(ctx, d, createdAccount.(*accounts.TokenAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] token account created (%s)", d.Id())
	return nil
}

func resourceTokenAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting token account (%s)", d.Id())

	client := m.(*client.Client)
	if err := accounts.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] token account deleted")
	return nil
}

func resourceTokenAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading token account (%s)", d.Id())

	client := m.(*client.Client)
	accountResource, err := accounts.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "token account")
	}

	if err := setTokenAccount(ctx, d, accountResource.(*accounts.TokenAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] token account read: (%s)", d.Id())
	return nil
}

func resourceTokenAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandTokenAccount(d)

	log.Printf("[INFO] updating token account: %#v", account)

	client := m.(*client.Client)
	updatedAccount, err := accounts.Update(client, account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTokenAccount(ctx, d, updatedAccount.(*accounts.TokenAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] token account updated (%s)", d.Id())
	return nil
}
