package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSSHKey() *schema.Resource {
	schemaMap := getCommonAccountsSchema()

	schemaMap[constUsername] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap[constPassphrase] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	return &schema.Resource{
		Create:        resourceSSHKeyCreate,
		Read:          resourceSSHKeyRead,
		Update:        resourceSSHKeyUpdate,
		DeleteContext: resourceAccountDeleteCommon,
		Schema:        schemaMap,
	}
}

func resourceSSHKeyRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingSSHKeyPair, id, err)
	}
	if account == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constAccount, m)

	accountResource := account.(*octopusdeploy.AccountResource)

	d.Set(constName, accountResource.Name)
	d.Set(constPassphrase, accountResource.PrivateKeyPassphrase)
	d.Set(constTenants, accountResource.TenantIDs)

	return nil
}

func buildSSHKeyResource(d *schema.ResourceData) (*octopusdeploy.SSHKeyAccount, error) {
	name := d.Get(constName).(string)
	username := d.Get(constUsername).(string)

	password := d.Get(constPassword).(string)
	if isEmpty(password) {
		log.Println("Key is nil. Must add in a password")
	}

	secretKey := octopusdeploy.NewSensitiveValue(password)

	account, err := octopusdeploy.NewSSHKeyAccount(name, username, secretKey)
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
