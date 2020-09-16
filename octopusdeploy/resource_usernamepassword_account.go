package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
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
			"tenant_tags": {
				Description: "The tags for the tenants that this step applies to",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	apiClient := m.(*client.Client)

	accountID := d.Id()
	account, err := apiClient.Accounts.Get(accountID)

	if err == client.ErrItemNotFound {
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

func buildUsernamePasswordResource(d *schema.ResourceData) *model.Account {
	var account, err = model.NewAccount(d.Get("name").(string), enum.UsernamePassword)

	if err != nil {
		return nil
	}

	account.Name = d.Get("name").(string)
	pass := d.Get("password").(string)
	account.Password = &model.SensitiveValue{NewValue: &pass}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentParticipation, _ = enum.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return account
}

func resourceUsernamePasswordCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	newAccount := buildUsernamePasswordResource(d)
	account, err := apiClient.Accounts.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error creating username password account %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceUsernamePasswordUpdate(d *schema.ResourceData, m interface{}) error {
	account := buildUsernamePasswordResource(d)
	account.ID = d.Id()

	apiClient := m.(*client.Client)

	updatedAccount, err := apiClient.Accounts.Update(*account)

	if err != nil {
		return fmt.Errorf("error updating username password account id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceUsernamePasswordDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	accountID := d.Id()

	err := apiClient.Accounts.Delete(accountID)

	if err != nil {
		return fmt.Errorf("error deleting username password account id %s: %s", accountID, err.Error())
	}

	d.SetId("")
	return nil
}
