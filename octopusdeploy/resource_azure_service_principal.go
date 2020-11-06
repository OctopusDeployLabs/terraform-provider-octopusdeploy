package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAzureServicePrincipal() *schema.Resource {
	validateSchema()
	resourceAzureServicePrincipalSchema := getCommonAccountsSchema()
	resourceAzureServicePrincipalSchema[constClientID] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	resourceAzureServicePrincipalSchema[constTenantID] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	resourceAzureServicePrincipalSchema[constSubscriptionNumber] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	resourceAzureServicePrincipalSchema[constKey] = &schema.Schema{
		Type:      schema.TypeString,
		Required:  true,
		Sensitive: true,
	}
	resourceAzureServicePrincipalSchema[constAzureEnvironment] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	resourceAzureServicePrincipalSchema[constResourceManagementEndpointBaseURI] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	resourceAzureServicePrincipalSchema[constActiveDirectoryEndpointBaseURI] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	return &schema.Resource{
		CreateContext: resourceAzureServicePrincipalCreate,
		DeleteContext: resourceAccountDeleteCommon,
		ReadContext:   resourceAzureServicePrincipalRead,
		Schema:        resourceAzureServicePrincipalSchema,
		UpdateContext: resourceAzureServicePrincipalUpdate,
	}
}

func buildAzureServicePrincipalResource(d *schema.ResourceData) (*octopusdeploy.AzureServicePrincipalAccount, error) {
	name := d.Get(constName).(string)
	password := d.Get(constKey).(string)
	secretKey := octopusdeploy.NewSensitiveValue(password)

	applicationID, err := uuid.Parse(d.Get(constClientID).(string))
	if err != nil {
		return nil, err
	}

	tenantID, err := uuid.Parse(d.Get(constTenantID).(string))
	if err != nil {
		return nil, err
	}

	subscriptionID, err := uuid.Parse(d.Get(constSubscriptionNumber).(string))
	if err != nil {
		return nil, err
	}

	account, err := octopusdeploy.NewAzureServicePrincipalAccount(name, subscriptionID, tenantID, applicationID, secretKey)
	if err != nil {
		return nil, err
	}

	// Optional Fields
	if v, ok := d.GetOk(constDescription); ok {
		account.Description = v.(string)
	}

	if v, ok := d.GetOk(constEnvironments); ok {
		account.EnvironmentIDs = getSliceFromTerraformTypeList(v)
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

	if v, ok := d.GetOk(constResourceManagementEndpointBaseURI); ok {
		account.ResourceManagerEndpoint = v.(string)
	}

	if v, ok := d.GetOk(constActiveDirectoryEndpointBaseURI); ok {
		account.AuthenticationEndpoint = v.(string)
	}

	err = account.Validate()
	if err != nil {
		return nil, err
	}

	return account, nil
}

func resourceAzureServicePrincipalCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, err := buildAzureServicePrincipalResource(d)
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

func resourceAzureServicePrincipalRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	accountResource, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccount(accountResource.(*octopusdeploy.AccountResource))
	if err != nil {
		return diag.FromErr(err)
	}

	account := accountResource.(*octopusdeploy.AzureServicePrincipalAccount)

	d.Set(constName, account.Name)
	d.Set(constDescription, account.Description)
	d.Set(constEnvironments, account.EnvironmentIDs)
	d.Set(constTenantedDeploymentParticipation, account.TenantedDeploymentMode)
	d.Set(constTenantTags, account.TenantTags)
	d.Set(constClientID, account.ApplicationID.String())
	d.Set(constTenantID, account.TenantID.String())
	d.Set(constSubscriptionNumber, account.SubscriptionID.String())

	// TODO: determine what to do here...
	// d.Set(constKey, account.ApplicationPassword)

	d.Set(constAzureEnvironment, account.AzureEnvironment)
	d.Set(constResourceManagementEndpointBaseURI, account.ResourceManagerEndpoint)
	d.Set(constActiveDirectoryEndpointBaseURI, account.AuthenticationEndpoint)

	return nil
}

func resourceAzureServicePrincipalUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, err := buildAzureServicePrincipalResource(d)
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
