package octopusdeploy

import (
	"context"
	"log"

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
		Create:        resourceSSHKeyCreate,
		ReadContext:   resourceSSHKeyRead,
		Update:        resourceSSHKeyUpdate,
		DeleteContext: resourceAccountDeleteCommon,
		Schema:        schemaMap,
	}
}

func resourceSSHKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func buildSSHKeyResource(d *schema.ResourceData) (*octopusdeploy.SSHKeyAccount, error) {
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

func resourceSSHKeyCreate(d *schema.ResourceData, m interface{}) error {
	account, err := buildSSHKeyResource(d)
	if err != nil {
		return err
	}

	client := m.(*octopusdeploy.Client)
	resource, err := client.Accounts.Add(account)
	if err != nil {
		return createResourceOperationError(errorCreatingSSHKeyPair, account.Name, err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func resourceSSHKeyUpdate(d *schema.ResourceData, m interface{}) error {
	account, err := buildSSHKeyResource(d)
	if err != nil {
		return err
	}
	account.ID = d.Id() // set ID so Octopus API knows which account to update

	client := m.(*octopusdeploy.Client)
	resource, err := client.Accounts.Update(account)
	if err != nil {
		return createResourceOperationError(errorUpdatingSSHKeyPair, d.Id(), err)
	}

	d.SetId(resource.GetID())

	return nil
}
