package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSSHKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSHKeyAccountCreate,
		DeleteContext: resourceAccountDeleteCommon,
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

	flattenSSHKeyAccount(ctx, d, createdSSHKeyAccount)
	return nil
}

func resourceSSHKeyAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccount(accountResource.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	sshKeyAccount := accountResource.(*octopusdeploy.SSHKeyAccount)

	flattenSSHKeyAccount(ctx, d, sshKeyAccount)
	return nil
}

func resourceSSHKeyAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account := expandSSHKeyAccount(d)

	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	updatedSSHKeyAccount := updatedAccount.(*octopusdeploy.SSHKeyAccount)

	flattenSSHKeyAccount(ctx, d, updatedSSHKeyAccount)
	return nil
}
