package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSSHKeyAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSHKeyAccountCreate,
		DeleteContext: resourceSSHKeyAccountDelete,
		Description:   "This resource manages SSH key accounts in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceSSHKeyAccountRead,
		Schema:        getSSHKeyAccountSchema(),
		UpdateContext: resourceSSHKeyAccountUpdate,
	}
}

func resourceSSHKeyAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandSSHKeyAccount(d)

	log.Printf("[INFO] creating SSH key account")

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setSSHKeyAccount(ctx, d, createdAccount.(*octopusdeploy.SSHKeyAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] SSH key account created (%s)", d.Id())
	return nil
}

func resourceSSHKeyAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting SSH key account (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Accounts.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] SSH key account deleted")
	return nil
}

func resourceSSHKeyAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading SSH key account (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] SSH key account (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	sshKeyAccount := accountResource.(*octopusdeploy.SSHKeyAccount)
	if err := setSSHKeyAccount(ctx, d, sshKeyAccount); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] SSH key account read (%s)", d.Id())
	return nil
}

func resourceSSHKeyAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating SSH key account (%s)", d.Id())

	account := expandSSHKeyAccount(d)
	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setSSHKeyAccount(ctx, d, updatedAccount.(*octopusdeploy.SSHKeyAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] SSH key account updated (%s)", d.Id())
	return nil
}
