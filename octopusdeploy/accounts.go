package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getCommonAccountsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		constName: {
			Required: true,
			Type:     schema.TypeString,
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

func buildAccountResourceCommon(d *schema.ResourceData) *octopusdeploy.AccountResource {
	name := d.Get(constName).(string)
	accountType := d.Get(constAccountType).(string)

	var accountResource = octopusdeploy.NewAccountResource(name, octopusdeploy.AccountType(accountType))

	if v, ok := d.GetOk(constTenantedDeploymentParticipation); ok {
		accountResource.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk(constTenantTags); ok {
		accountResource.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constTenants); ok {
		accountResource.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	return accountResource
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

func resourceAccountUpdateCommon(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountResource := buildAccountResourceCommon(d)
	accountResource.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(accountResource)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(updatedAccount.GetID())

	return nil
}

func resourceAccountDeleteCommon(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.Accounts.DeleteByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(constEmptyString)

	return nil
}
