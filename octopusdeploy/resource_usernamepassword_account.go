package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUsernamePassword() *schema.Resource {
	return &schema.Resource{
		Create: resourceUsernamePasswordCreate,
		Read:   resourceUsernamePasswordRead,
		Update: resourceUsernamePasswordUpdate,
		Delete: resourceUsernamePasswordDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenanted_deployment_participation": getTenantedDeploymentSchema(),
			"environments": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"account_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UsernamePassword",
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceUsernamePasswordRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountID := d.Id()
	account, err := client.Account.Get(accountID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading username password account %s: %s", accountID, err.Error())
	}

	d.Set("name", account.Name)
	d.Set("description", account.Description)
	d.Set("environments", account.EnvironmentIDs)
	d.Set("password", account.Password)

	return nil
}

func buildUsernamePasswordResource(d *schema.ResourceData) *octopusdeploy.Account {
	var account = octopusdeploy.NewAccount(d.Get("name").(string), octopusdeploy.UsernamePassword)

	account.Name = d.Get("name").(string)
	account.Password = octopusdeploy.SensitiveValue{NewValue: d.Get("password").(string)}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentParticipation, _ = octopusdeploy.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return account
}

func resourceUsernamePasswordCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newAccount := buildUsernamePasswordResource(d)
	account, err := client.Account.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error creating username password account %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceUsernamePasswordUpdate(d *schema.ResourceData, m interface{}) error {
	account := buildUsernamePasswordResource(d)
	account.ID = d.Id()

	client := m.(*octopusdeploy.Client)

	updatedAccount, err := client.Account.Update(account)

	if err != nil {
		return fmt.Errorf("error updating username password account id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceUsernamePasswordDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountID := d.Id()

	err := client.Account.Delete(accountID)

	if err != nil {
		return fmt.Errorf("error deleting username password account id %s: %s", accountID, err.Error())
	}

	d.SetId("")
	return nil
}
