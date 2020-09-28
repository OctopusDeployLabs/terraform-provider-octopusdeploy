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
	id := d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.GetByID(id)
	if err != nil {
		return nil, createResourceOperationError(errorReadingAccount, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil, fmt.Errorf(errorAccountNotFound, id)
	}

	d.Set(constName, resource.Name)
	d.Set(constDescription, resource.Description)
	d.Set(constEnvironments, resource.EnvironmentIDs)
	d.Set(constTenantedDeploymentParticipation, resource.TenantedDeploymentParticipation)
	d.Set(constTenants, resource.TenantIDs)
	d.Set(constTenantTags, resource.EnvironmentIDs)

	return resource, nil
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
	apiClient := m.(*client.Client)
	account, err := apiClient.Accounts.Add(account)
	if err != nil {
		return createResourceOperationError(errorCreatingAccount, account.Name, err)
	}

	d.SetId(account.ID)

	return nil
}

func resourceAccountUpdateCommon(d *schema.ResourceData, m interface{}, account *model.Account) error {
	account.ID = d.Id()

	apiClient := m.(*client.Client)
	updatedAccount, err := apiClient.Accounts.Update(*account)
	if err != nil {
		return createResourceOperationError(errorUpdatingAccount, d.Id(), err)
	}

	d.SetId(updatedAccount.ID)

	return nil
}

func resourceAccountDeleteCommon(d *schema.ResourceData, m interface{}) error {
	accountID := d.Id()

	apiClient := m.(*client.Client)
	err := apiClient.Accounts.DeleteByID(accountID)
	if err != nil {
		return createResourceOperationError(errorDeletingAccount, accountID, err)
	}

	d.SetId(constEmptyString)

	return nil
}
