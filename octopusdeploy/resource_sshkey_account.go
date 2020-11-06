package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSSHKey() *schema.Resource {
	schemaMap := getCommonAccountsSchema()
	schemaMap[constUsername] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap[constPassphrase] = &schema.Schema{
		Optional:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	return &schema.Resource{
		CreateContext: resourceSSHKeyAccountCreate,
		DeleteContext: resourceAccountDeleteCommon,
		ReadContext:   resourceSSHKeyAccountRead,
		Schema:        schemaMap,
		UpdateContext: resourceSSHKeyAccountUpdate,
	}
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

	account := accountResource.(*octopusdeploy.SSHKeyAccount)

	d.Set(constDescription, account.Description)
	d.Set(constEnvironments, account.EnvironmentIDs)
	d.Set(constName, account.GetName())
	d.Set(constTenantedDeploymentParticipation, account.TenantedDeploymentMode)
	d.Set(constTenants, account.TenantIDs)
	d.Set(constTenantTags, account.TenantTags)

	// TODO: determine what to do here...
	// d.Set(constPassphrase, account.PrivateKeyPassphrase)

	d.SetId(account.GetID())

	return nil
}

func buildSSHKeyAccount(d *schema.ResourceData) (*octopusdeploy.SSHKeyAccount, error) {
	var name string
	if v, ok := d.GetOk(constName); ok {
		name = v.(string)
	}

	var username string
	if v, ok := d.GetOk(constUsername); ok {
		username = v.(string)
	}

	var passphrase string
	if v, ok := d.GetOk(constPassphrase); ok {
		passphrase = v.(string)
	}

	account, err := octopusdeploy.NewSSHKeyAccount(name, username, octopusdeploy.NewSensitiveValue(passphrase))
	if err != nil {
		return nil, err
	}

	if v, ok := d.GetOk(constTenantedDeploymentParticipation); ok {
		account.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk(constTenantTags); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constTenants); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	return account, nil
}

func resourceSSHKeyAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, err := buildSSHKeyAccount(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := m.(*octopusdeploy.Client)
	sshKeyAccount, err := client.Accounts.Add(account)
	if err != nil {
		diag.FromErr(err)
	}

	d.SetId(sshKeyAccount.GetID())

	return nil
}

func resourceSSHKeyAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, err := buildSSHKeyAccount(d)
	if err != nil {
		return diag.FromErr(err)
	}
	account.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(updatedAccount.GetID())

	return nil
}
