package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func fetchAndReadAccount(d *schema.ResourceData, m interface{}) (*model.Account, error) {
	octopusClient := m.(*client.Client)

	accountID := d.Id()
	account, err := octopusClient.Accounts.Get(accountID)

	if err == client.ErrItemNotFound {
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
	d.Set("tenants", account.TenantIDs)
	d.Set("tenant_tags", account.EnvironmentIDs)

	return account, nil
}

func buildAccountResourceCommon(d *schema.ResourceData, accountType enum.AccountType) *model.Account {
	var account, _ = model.NewAccount(d.Get("name").(string), accountType)

	if account == nil {
		log.Println(nameIsNil("buildAccountResourceCommon"))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentParticipation, _ = enum.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenants"); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return account
}

func resourceAccountCreateCommon(d *schema.ResourceData, m interface{}, account *model.Account) error {
	octopusClient := m.(*client.Client)

	account, err := octopusClient.Accounts.Add(account)

	if err != nil {
		return fmt.Errorf("error creating account %s: %s", account.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceAccountUpdateCommon(d *schema.ResourceData, m interface{}, account *model.Account) error {
	account.ID = d.Id()

	octopusClient := m.(*client.Client)

	updatedAccount, err := octopusClient.Accounts.Update(*account)

	if err != nil {
		return fmt.Errorf("error updating username password account id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceAccountDeleteCommon(d *schema.ResourceData, m interface{}) error {
	octopusClient := m.(*client.Client)

	accountID := d.Id()

	err := octopusClient.Accounts.Delete(accountID)

	if err != nil {
		return fmt.Errorf("error deleting account ID %s: %s", accountID, err.Error())
	}

	d.SetId("")
	return nil
}
