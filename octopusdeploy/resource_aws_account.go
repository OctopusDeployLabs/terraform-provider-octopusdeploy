package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAmazonWebServicesAccount() *schema.Resource {
	validateSchema()
	schemaMap := getCommonAccountsSchema()
	schemaMap[constAccessKey] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	schemaMap[constSecretKey] = &schema.Schema{
		Type:      schema.TypeString,
		Required:  true,
		Sensitive: true,
	}
	return &schema.Resource{
		CreateContext: resourceAmazonWebServicesAccountCreate,
		DeleteContext: resourceAccountDeleteCommon,
		ReadContext:   resourceAmazonWebServicesAccountRead,
		UpdateContext: resourceAmazonWebServicesAccountUpdate,
		Schema:        schemaMap,
	}
}

func resourceAmazonWebServicesAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	d.Set(constAccessKey, account.AccessKey)
	d.Set(constDescription, account.Description)
	d.Set(constEnvironments, account.EnvironmentIDs)
	d.Set(constName, account.GetName())
	d.Set(constTenantedDeploymentParticipation, account.TenantedDeploymentMode)
	d.Set(constTenants, account.TenantIDs)
	d.Set(constTenantTags, account.TenantTags)

	// TODO: determine what to do here...
	// d.Set(constSecretKey, account.SecretKey)

	d.SetId(account.GetID())

	return nil
}

func buildAmazonWebServicesAccountResource(d *schema.ResourceData) (*octopusdeploy.AmazonWebServicesAccount, error) {
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

func resourceAmazonWebServicesAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	amazonWebServicesAccount, err := buildAmazonWebServicesAccountResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.Add(amazonWebServicesAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(account.GetID())

	return nil
}

func resourceAmazonWebServicesAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	amazonWebServicesAccount, err := buildAmazonWebServicesAccountResource(d)
	if err != nil {
		return diag.FromErr(err)
	}
	amazonWebServicesAccount.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.Update(amazonWebServicesAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(account.GetID())

	return nil
}
