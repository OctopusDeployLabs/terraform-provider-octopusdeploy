package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/enum"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environments": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"passphrase": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSSHKeyRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	accountID := d.Id()
	account, err := apiClient.Accounts.Get(accountID)

	if err == client.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading SSH Key Pair %s: %s", accountID, err.Error())
	}

	d.Set("name", account.Name)
	d.Set("passphrase", account.Password)
	d.Set("tenants", account.TenantIDs)

	return nil
}

func buildSSHKeyResource(d *schema.ResourceData) *model.Account {
	account, err := model.NewAccount(d.Get("name").(string), enum.SshKeyPair)
	if err != nil {
		return nil
	}

	account.Name = d.Get("username").(string)
	password := d.Get("password").(string)
	account.Password = &model.SensitiveValue{NewValue: &password}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentParticipation, _ = enum.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenants"); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	return account
}

func resourceSSHKeyCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	newAccount := buildSSHKeyResource(d)
	account, err := apiClient.Accounts.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error reading SSH Key Pair %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceSSHKeyUpdate(d *schema.ResourceData, m interface{}) error {
	account := buildSSHKeyResource(d)
	account.ID = d.Id()

	apiClient := m.(*client.Client)

	updatedAccount, err := apiClient.Accounts.Update(account)

	if err != nil {
		return fmt.Errorf("error reading SSH Key Pair %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceSSHKeyDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	accountID := d.Id()

	err := apiClient.Accounts.Delete(accountID)

	if err != nil {
		return fmt.Errorf("error reading SSH Key Pair id %s: %s", accountID, err.Error())
	}

	d.SetId("")
	return nil
}
