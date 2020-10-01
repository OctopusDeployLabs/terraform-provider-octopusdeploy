package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func fetchAndReadAccount(d *schema.ResourceData, m interface{}) (*model.Account, diag.Diagnostics) {
	id := d.Id()

	var diags diag.Diagnostics

	if diags == nil {
		log.Println("diag package is empty")
	}

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.GetByID(id)
	if err != nil {
		// return nil, createResourceOperationError(errorReadingAccount, id, err)
		diag.FromErr(err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		// return nil, fmt.Errorf(errorAccountNotFound, id)
		diag.Errorf(errorAccountNotFound, id)
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

func resourceAccountCreateCommon(ctx context.Context, d *schema.ResourceData, m interface{}, account *model.Account) diag.Diagnostics {
	apiClient := m.(*client.Client)

	var diags diag.Diagnostics

	if diags == nil {
		log.Println("diag package is empty")
	}

	account, err := apiClient.Accounts.Add(account)
	if err != nil {
		// return createResourceOperationError(errorCreatingAccount, account.Name, err)
		diag.FromErr(err)
	}

	d.SetId(account.ID)

	return diags
}

func resourceAccountUpdateCommon(ctx context.Context, d *schema.ResourceData, m interface{}, account *model.Account) diag.Diagnostics {
	account.ID = d.Id()

	var diags diag.Diagnostics

	if diags == nil {
		log.Println("diag package is empty")
	}

	apiClient := m.(*client.Client)
	updatedAccount, err := apiClient.Accounts.Update(*account)
	if err != nil {
		// return createResourceOperationError(errorUpdatingAccount, d.Id(), err)
		diag.FromErr(err)
	}

	d.SetId(updatedAccount.ID)

	return nil
}

func resourceAccountDeleteCommon(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountID := d.Id()

	var diags diag.Diagnostics

	if diags == nil {
		log.Println("diag package is empty")
	}

	apiClient := m.(*client.Client)
	err := apiClient.Accounts.DeleteByID(accountID)
	if err != nil {
		// return createResourceOperationError(errorDeletingAccount, accountID, err)
		diag.FromErr(err)
	}

	d.SetId(constEmptyString)

	return nil
}
