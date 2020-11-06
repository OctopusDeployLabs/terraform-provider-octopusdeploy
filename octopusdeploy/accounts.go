package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

func resourceAccountCreateCommon(ctx context.Context, d *schema.ResourceData, m interface{}, account octopusdeploy.IAccount) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	return nil
}

func resourceAccountUpdateCommon(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountResource := buildAccountResource(d)
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
	client := m.(*octopusdeploy.Client)
	err := client.Accounts.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(constEmptyString)

	return nil
}
