package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTokenAccount() *schema.Resource {
	validateSchema()
	resourceTokenAccountSchema := getCommonAccountsSchema()
	resourceTokenAccountSchema[constToken] = &schema.Schema{
		Required:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	return &schema.Resource{
		CreateContext: resourceTokenAccountCreate,
		DeleteContext: resourceAccountDeleteCommon,
		ReadContext:   resourceTokenAccountRead,
		Schema:        resourceTokenAccountSchema,
		UpdateContext: resourceTokenAccountUpdate,
	}
}

func resourceTokenAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccount(accountResource.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	account := accountResource.(*octopusdeploy.TokenAccount)

	d.Set(constDescription, account.Description)
	d.Set(constEnvironments, account.EnvironmentIDs)
	d.Set(constName, account.GetName())
	d.Set(constTenantedDeploymentParticipation, account.TenantedDeploymentMode)
	d.Set(constTenants, account.TenantIDs)
	d.Set(constTenantTags, account.TenantTags)

	// TODO: determine how to persist this sensitive value
	// d.Set(constToken, account.Token)

	d.SetId(account.GetID())

	return nil
}

func buildTokenAccountResource(d *schema.ResourceData) (*octopusdeploy.TokenAccount, error) {
	name := d.Get(constName).(string)
	token := d.Get(constToken).(string)
	tokenSensitiveValue := octopusdeploy.NewSensitiveValue(token)

	account, err := octopusdeploy.NewTokenAccount(name, tokenSensitiveValue)
	if err != nil {
		return nil, err
	}

	if v, ok := d.GetOk(constTenantedDeploymentParticipation); ok {
		account.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk(constTenantTags); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constTenants); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	return account, nil
}

func resourceTokenAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, err := buildTokenAccountResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	return nil
}

func resourceTokenAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, err := buildTokenAccountResource(d)
	if err != nil {
		return diag.FromErr(err)
	}
	account.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(updatedAccount.GetID())

	return nil
}
