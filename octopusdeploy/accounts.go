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
		constName: {
			Type:     schema.TypeString,
			Required: true,
		},
		constDescription: {
			Type:     schema.TypeString,
			Optional: true,
		},
		constEnvironments: {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
		},
		constTenantedDeploymentParticipation: getTenantedDeploymentSchema(),
		constTenants: {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
		},
		constTenantTags: {
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
	account, err := octopusClient.Accounts.GetByID(accountID)

	if err != nil {
		return nil, fmt.Errorf(errorReadingAccount, accountID, err.Error())
	}

	if account == nil {
		d.SetId(constEmptyString)
		return nil, fmt.Errorf(errorAccountNotFound, accountID)
	}

	d.Set(constName, account.Name)
	d.Set(constDescription, account.Description)
	d.Set(constEnvironments, account.EnvironmentIDs)
	d.Set(constTenantedDeploymentParticipation, account.TenantedDeploymentParticipation)
	d.Set(constTenants, account.TenantIDs)
	d.Set(constTenantTags, account.EnvironmentIDs)

	return account, nil
}

func buildAccountResourceCommon(d *schema.ResourceData, accountType enum.AccountType) *model.Account {
	var account, _ = model.NewAccount(d.Get(constName).(string), accountType)

	if account == nil {
		log.Println(nameIsNil("buildAccountResourceCommon"))
	}

	if v, ok := d.GetOk(constTenantTags); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constTenantedDeploymentParticipation); ok {
		account.TenantedDeploymentParticipation, _ = enum.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk(constTenants); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constTenantTags); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return account
}

func resourceAccountCreateCommon(d *schema.ResourceData, m interface{}, account *model.Account) error {
	octopusClient := m.(*client.Client)

	account, err := octopusClient.Accounts.Add(account)

	if err != nil {
		return createResourceOperationError(errorCreatingAccount, account.Name, err)
	}

	d.SetId(account.ID)

	return nil
}

func resourceAccountUpdateCommon(d *schema.ResourceData, m interface{}, account *model.Account) error {
	account.ID = d.Id()

	octopusClient := m.(*client.Client)

	updatedAccount, err := octopusClient.Accounts.Update(*account)

	if err != nil {
		return createResourceOperationError(errorUpdatingAccount, d.Id(), err)
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceAccountDeleteCommon(d *schema.ResourceData, m interface{}) error {
	octopusClient := m.(*client.Client)

	accountID := d.Id()

	err := octopusClient.Accounts.DeleteByID(accountID)

	if err != nil {
		return createResourceOperationError(errorDeletingAccount, accountID, err)
	}

	d.SetId(constEmptyString)
	return nil
}
