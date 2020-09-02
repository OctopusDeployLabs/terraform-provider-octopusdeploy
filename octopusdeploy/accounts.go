package model

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func getCommonAccountsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
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
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
		},
		"tenant_tags": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
		},
	}
}

func fetchAndReadAccount(d *schema.ResourceData, m interface{}) (*octopusdeploy.Account, error) {
	client := m.(*octopusdeploy.Client)

	accountID := d.Id()
	account, err := client.Account.Get(accountID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil, fmt.Errorf("account %s not found: %s ", accountID, err.Error())
	}

	if err != nil {
		return nil, fmt.Errorf("error readingaccount %s: %s", accountID, err.Error())
	}

	d.Set("name", account.Name)
	d.Set("description", account.Description)
	d.Set("environments", account.EnvironmentIDs)
	d.Set("tenanted_deployment_participation", account.TenantedDeploymentParticipation)
	d.Set("tenants", account.TenantID)
	d.Set("tenant_tags", account.EnvironmentIDs)

	return account, nil
}

func buildAccountResourceCommon(d *schema.ResourceData, accountType octopusdeploy.AccountType) *octopusdeploy.Account {
	var account = octopusdeploy.NewAccount(d.Get("name").(string), accountType)

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentParticipation, _ = octopusdeploy.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenants"); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return account
}

func resourceAccountCreateCommon(d *schema.ResourceData, m interface{}, account *octopusdeploy.Account) error {
	client := m.(*octopusdeploy.Client)

	account, err := client.Account.Add(account)

	if err != nil {
		return fmt.Errorf("error creating account %s: %s", account.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceAccountUpdateCommon(d *schema.ResourceData, m interface{}, account *octopusdeploy.Account) error {
	account.ID = d.Id()

	client := m.(*octopusdeploy.Client)

	updatedAccount, err := client.Account.Update(account)

	if err != nil {
		return fmt.Errorf("error updating username password account id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceAccountDeleteCommon(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountID := d.Id()

	err := client.Account.Delete(accountID)

	if err != nil {
		return fmt.Errorf("error deleting username password account id %s: %s", accountID, err.Error())
	}

	d.SetId("")
	return nil
}
