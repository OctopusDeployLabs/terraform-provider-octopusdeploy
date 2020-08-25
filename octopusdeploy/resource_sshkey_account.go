package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSSHKeyCreate,
		Read:   resourceSSHKeyRead,
		Update: resourceSSHKeyUpdate,
		Delete: resourceSSHKeyDelete,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"passphrase": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSSHKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountID := d.Id()
	account, err := client.Account.Get(accountID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading SSH Key Pair %s: %s", accountID, err.Error())
	}

	d.Set("name", account.Name)
	d.Set("passphrase", account.Password)

	return nil
}

func buildSSHKeyResource(d *schema.ResourceData) *octopusdeploy.Account {
	var account = octopusdeploy.NewAccount(d.Get("name").(string), octopusdeploy.SshKeyPair)

	account.Name = d.Get("username").(string)
	account.Password = octopusdeploy.SensitiveValue{NewValue: d.Get("password").(string)}

	return account
}

func resourceSSHKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newAccount := buildSSHKeyResource(d)
	account, err := client.Account.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error reading SSH Key Pair %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceSSHKeyUpdate(d *schema.ResourceData, m interface{}) error {
	account := buildSSHKeyResource(d)
	account.ID = d.Id()

	client := m.(*octopusdeploy.Client)

	updatedAccount, err := client.Account.Update(account)

	if err != nil {
		return fmt.Errorf("error reading SSH Key Pair %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceSSHKeyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountID := d.Id()

	err := client.Account.Delete(accountID)

	if err != nil {
		return fmt.Errorf("error reading SSH Key Pair id %s: %s", accountID, err.Error())
	}

	d.SetId("")
	return nil
}
