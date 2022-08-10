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

func resourceUsernamePasswordAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUsernamePasswordAccountCreate,
		DeleteContext: resourceUsernamePasswordAccountDelete,
		Description:   "This resource manages username-password accounts in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceUsernamePasswordAccountRead,
		Schema:        getUsernamePasswordAccountSchema(),
		UpdateContext: resourceUsernamePasswordAccountUpdate,
	}
}

func resourceUsernamePasswordAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandUsernamePasswordAccount(d)

	log.Printf("[INFO] creating username-password account: %#v", account)

	client := m.(*client.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setUsernamePasswordAccount(ctx, d, createdAccount.(*accounts.UsernamePasswordAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] username-password account created (%s)", d.Id())
	return nil
}

func resourceUsernamePasswordAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting username-password account (%s)", d.Id())

	client := m.(*client.Client)
	if err := client.Accounts.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] username-password account deleted")
	return nil
}

func resourceUsernamePasswordAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading username-password account (%s)", d.Id())

	client := m.(*client.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "username-password account")
	}

	if err := setUsernamePasswordAccount(ctx, d, accountResource.(*accounts.UsernamePasswordAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] username-password account read: (%s)", d.Id())
	return nil
}

func resourceUsernamePasswordAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandUsernamePasswordAccount(d)

	log.Printf("[INFO] updating username-password account: %#v", account)

	client := m.(*client.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setUsernamePasswordAccount(ctx, d, updatedAccount.(*accounts.UsernamePasswordAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] username-password account updated (%s)", d.Id())
	return nil
}
