package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSSHKeyAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSHKeyAccountCreate,
		DeleteContext: resourceSSHKeyAccountDelete,
		Importer:      getImporter(),
		ReadContext:   resourceSSHKeyAccountRead,
		Schema:        getSSHKeyAccountSchema(),
		UpdateContext: resourceSSHKeyAccountUpdate,
	}
}

func resourceSSHKeyAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandSSHKeyAccount(d)

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		diag.FromErr(err)
	}

	createdSSHKeyAccount := createdAccount.(*octopusdeploy.SSHKeyAccount)

	setSSHKeyAccount(ctx, d, createdSSHKeyAccount)
	return nil
}

func resourceSSHKeyAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Accounts.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceSSHKeyAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if accountResource == nil {
		d.SetId("")
		return nil
	}

	accountResource, err = octopusdeploy.ToAccount(accountResource.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	sshKeyAccount := accountResource.(*octopusdeploy.SSHKeyAccount)

	setSSHKeyAccount(ctx, d, sshKeyAccount)
	return nil
}

func resourceSSHKeyAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandSSHKeyAccount(d)

	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccount(accountResource.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	updatedSSHKeyAccount := accountResource.(*octopusdeploy.SSHKeyAccount)

	setSSHKeyAccount(ctx, d, updatedSSHKeyAccount)
	return nil
}
