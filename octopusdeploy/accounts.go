package octopusdeploy

import (
	"context"
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getCommonAccountsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		constName: {
			Required: true,
			Type:     schema.TypeString,
		},
		constDescription: {
			Optional: true,
			Type:     schema.TypeString,
		},
		constEnvironments: {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		constTenantedDeploymentParticipation: getTenantedDeploymentSchema(),
		constTenants: {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		constTenantTags: {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}

func fetchAndReadAccount(d *schema.ResourceData, m interface{}) (octopusdeploy.IAccount, error) {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.GetByID(id)
	if err != nil {
		return nil, createResourceOperationError(errorReadingAccount, id, err)
	}
	if account == nil {
		d.SetId(constEmptyString)
		return nil, fmt.Errorf(errorAccountNotFound, id)
	}

	accountResource := account.(*octopusdeploy.AccountResource)

	d.Set(constName, accountResource.GetName())
	d.Set(constDescription, accountResource.Description)
	d.Set(constEnvironments, accountResource.EnvironmentIDs)
	d.Set(constTenantedDeploymentParticipation, accountResource.TenantedDeploymentMode)
	d.Set(constTenants, accountResource.TenantIDs)
	d.Set(constTenantTags, accountResource.EnvironmentIDs)

	return accountResource, nil
}

func buildAccountResourceCommon(d *schema.ResourceData, accountType string) octopusdeploy.IAccount {
	var account = octopusdeploy.NewAccountResource(d.Get(constName).(string), octopusdeploy.AccountType(accountType))

	if account == nil {
		log.Println(nameIsNil("buildAccountResourceCommon"))
	}

	if v, ok := d.GetOk(constTenantTags); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constTenantedDeploymentParticipation); ok {
		account.TenantedDeploymentMode = v.(octopusdeploy.TenantedDeploymentMode)
	}

	if v, ok := d.GetOk(constTenants); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constTenantTags); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return account
}

func resourceAccountCreateCommon(d *schema.ResourceData, m interface{}, account octopusdeploy.IAccount) error {
	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.Add(account)
	if err != nil {
		return createResourceOperationError(errorCreatingAccount, account.GetName(), err)
	}

	d.SetId(account.GetID())

	return nil
}

func resourceAccountUpdateCommon(d *schema.ResourceData, m interface{}, accountResource *octopusdeploy.AccountResource) error {
	accountResource.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(accountResource)
	if err != nil {
		return createResourceOperationError(errorUpdatingAccount, d.Id(), err)
	}

	d.SetId(updatedAccount.GetID())

	return nil
}

func resourceAccountDeleteCommon(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountID := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.Accounts.DeleteByID(accountID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(constEmptyString)

	return nil
}
