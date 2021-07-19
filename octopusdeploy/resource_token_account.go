package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTokenAccount(ctx, d, createdAccount.(*octopusdeploy.TokenAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] token account created (%s)", d.Id())
	return nil
}

func resourceTokenAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting token account (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Accounts.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] token account deleted")
	return nil
}

func resourceTokenAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading token account (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] token account (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if err := setTokenAccount(ctx, d, accountResource.(*octopusdeploy.TokenAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] token account read: (%s)", d.Id())
	return nil
}

func resourceTokenAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandTokenAccount(d)

	log.Printf("[INFO] updating token account: %#v", account)

	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTokenAccount(ctx, d, updatedAccount.(*octopusdeploy.TokenAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] token account updated (%s)", d.Id())
	return nil
}
