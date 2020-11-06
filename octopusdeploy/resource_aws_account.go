package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAWSAccount() *schema.Resource {
	validateSchema()
	resourceAWSAccountSchema := getCommonAccountsSchema()
	resourceAWSAccountSchema[constAccessKey] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	resourceAWSAccountSchema[constSecretKey] = &schema.Schema{
		Type:      schema.TypeString,
		Required:  true,
		Sensitive: true,
	}
	return &schema.Resource{
		CreateContext: resourceAWSAccountCreate,
		DeleteContext: resourceAccountDeleteCommon,
		ReadContext:   resourceAWSAccountRead,
		Schema:        resourceAWSAccountSchema,
		UpdateContext: resourceAWSAccountUpdate,
	}
}

func resourceAWSAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccount(accountResource.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	account := accountResource.(*octopusdeploy.AmazonWebServicesAccount)
	flattenAWSAccount(ctx, d, account)

	return nil
}

func buildAWSAccountResource(d *schema.ResourceData) (*octopusdeploy.AmazonWebServicesAccount, error) {
	name := d.Get(constName).(string)
	accessKey := d.Get(constAccessKey).(string)
	password := d.Get(constSecretKey).(string)
	secretKey := octopusdeploy.NewSensitiveValue(password)

	account, err := octopusdeploy.NewAmazonWebServicesAccount(name, accessKey, secretKey)
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

func resourceAWSAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, err := buildAWSAccountResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		return diag.FromErr(err)
	}

	account = createdAccount.(*octopusdeploy.AmazonWebServicesAccount)
	flattenAWSAccount(ctx, d, account)

	return nil
}

func resourceAWSAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, err := buildAWSAccountResource(d)
	if err != nil {
		return diag.FromErr(err)
	}
	account.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		return diag.FromErr(err)
	}

	account = updatedAccount.(*octopusdeploy.AmazonWebServicesAccount)
	flattenAWSAccount(ctx, d, account)

	return nil
}
