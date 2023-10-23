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

	log.Printf("[INFO] creating SSH key account: %#v", account)

	client := m.(*client.Client)
	createdAccount, err := accounts.Add(client, account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setSSHKeyAccount(ctx, d, createdAccount.(*accounts.SSHKeyAccount)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	log.Printf("[INFO] SSH key account created (%s)", d.Id())
	return nil
}

func resourceSSHKeyAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting SSH key account (%s)", d.Id())

	client := m.(*client.Client)
	if err := accounts.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] SSH key account deleted")
	return nil
}

func resourceSSHKeyAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading SSH key account (%s)", d.Id())

	client := m.(*client.Client)
	accountResource, err := accounts.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "SSH key account")
	}

	sshKeyAccount := accountResource.(*accounts.SSHKeyAccount)
	if err := setSSHKeyAccount(ctx, d, sshKeyAccount); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] SSH key account read (%s)", d.Id())
	return nil
}

func resourceSSHKeyAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating SSH key account (%s)", d.Id())

	account := expandSSHKeyAccount(d)
	client := m.(*client.Client)
	updatedAccount, err := accounts.Update(client, account)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setSSHKeyAccount(ctx, d, updatedAccount.(*accounts.SSHKeyAccount)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] SSH key account updated (%s)", d.Id())
	return nil
}
