package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTokenAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTokenAccountCreate,
		Read:   resourceTokenAccountRead,
		Update: resourceTokenAccountUpdate,
		Delete: resourceTokenAccountDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Token",
			},
			"environments": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"tenant_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenanted_deployment_participation": getTenantedDeploymentSchema(),
			"token": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTokenAccountRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountID := d.Id()
	account, err := client.Account.Get(accountID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading token account %s: %s", accountID, err.Error())
	}

	d.Set("name", account.Name)
	d.Set("description", account.Description)
	d.Set("environments", account.EnvironmentIDs)
	d.Set("tenanted_deployment_participation", account.TenantedDeploymentParticipation.String())
	d.Set("tenant_tags", account.TenantTags)
	d.Set("token", account.Token)

	return nil
}

func buildTokenAccountResource(d *schema.ResourceData) *octopusdeploy.Account {
	var account = octopusdeploy.NewAccount(d.Get("name").(string), octopusdeploy.Token)

	account.Name = d.Get("name").(string)
	account.Description = d.Get("description").(string)

	account.Token = octopusdeploy.SensitiveValue{NewValue: d.Get("token").(string)}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentParticipation, _ = octopusdeploy.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("environments"); ok {
		account.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	return account
}

func resourceTokenAccountCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newAccount := buildTokenAccountResource(d)
	account, err := client.Account.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error creating Token account %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceTokenAccountUpdate(d *schema.ResourceData, m interface{}) error {
	account := buildTokenAccountResource(d)
	account.ID = d.Id()

	client := m.(*octopusdeploy.Client)

	updatedAccount, err := client.Account.Update(account)

	if err != nil {
		return fmt.Errorf("error updating Token acccount id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceTokenAccountDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountID := d.Id()

	err := client.Account.Delete(accountID)

	if err != nil {
		return fmt.Errorf("error deleting Token account id %s: %s", accountID, err.Error())
	}

	d.SetId("")
	return nil
}
