package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAccount() *schema.Resource {
	schemaMap := getCommonAccountsSchema()
	schemaMap[constAccountType] = getAccountTypeSchema()
	schemaMap[constAccessKey] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap[constUsername] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap[constPassword] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap[constPassphrase] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap[constSecretKey] = &schema.Schema{
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	}
	schemaMap[constToken] = &schema.Schema{
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	}
	schemaMap[constClientID] = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	schemaMap[constClientSecret] = &schema.Schema{
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	}
	schemaMap[constTenantID] = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	schemaMap[constSubscriptionID] = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	schemaMap[constKey] = &schema.Schema{
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	}
	schemaMap[constAzureEnvironment] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap[constResourceManagementEndpointBaseURI] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap[constActiveDirectoryEndpointBaseURI] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	return &schema.Resource{
		CreateContext: resourceAccountCreate,
		DeleteContext: resourceAccountDeleteCommon,
		ReadContext:   resourceAccountRead,
		Schema:        schemaMap,
		UpdateContext: resourceAccountUpdate,
	}
}

func resourceAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*octopusdeploy.Client)

	var diags diag.Diagnostics
	accountID := d.Id()

	account, err := client.Accounts.GetByID(accountID)
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource := account.(*octopusdeploy.AccountResource)

	d.Set(constAccountType, accountResource.AccountType)
	d.Set(constActiveDirectoryEndpointBaseURI, accountResource.AuthenticationEndpoint)
	d.Set(constAzureEnvironment, accountResource.AzureEnvironment)
	d.Set(constClientID, accountResource.ApplicationID.String())
	d.Set(constDescription, accountResource.Description)
	d.Set(constEnvironments, accountResource.EnvironmentIDs)
	d.Set(constName, accountResource.Name)
	d.Set(constPassphrase, accountResource.PrivateKeyPassphrase)
	d.Set(constResourceManagementEndpointBaseURI, accountResource.ResourceManagerEndpoint)
	d.Set(constSubscriptionID, accountResource.SubscriptionID.String())
	d.Set(constTenantedDeploymentParticipation, accountResource.TenantedDeploymentMode)
	d.Set(constTenantID, accountResource.TenantID.String())
	d.Set(constTenants, accountResource.TenantIDs)
	d.Set(constTenantTags, accountResource.TenantTags)
	d.Set(constToken, accountResource.Token)
	d.Set(constUsername, accountResource.Username)

	if accountResource.ApplicationPassword.HasValue {
		if accountResource.ApplicationPassword.NewValue != nil {
			d.Set(constClientSecret, accountResource.ApplicationPassword.NewValue)
		}
	}

	return diags
}

func buildAccountResource(d *schema.ResourceData) *octopusdeploy.AccountResource {
	var name string
	if v, ok := d.GetOk(constName); ok {
		name = v.(string)
	}

	var accountType octopusdeploy.AccountType
	if v, ok := d.GetOk(constAccountType); ok {
		accountType = octopusdeploy.AccountType(v.(string))
	}

	account := octopusdeploy.NewAccountResource(name, accountType)

	if v, ok := d.GetOk(constAccessKey); ok {
		account.AccessKey = v.(string)
	}

	if v, ok := d.GetOk(constClientID); ok {
		clientID := uuid.MustParse(v.(string))
		account.ApplicationID = &clientID
	}

	if v, ok := d.GetOk(constClientSecret); ok {
		account.ApplicationPassword = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk(constActiveDirectoryEndpointBaseURI); ok {
		account.AuthenticationEndpoint = v.(string)
	}

	if v, ok := d.GetOk(constAzureEnvironment); ok {
		account.AzureEnvironment = v.(string)
	}

	if v, ok := d.GetOk(constCertificateData); ok {
		account.CertificateBytes = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk(constCertificateThumbprint); ok {
		account.CertificateThumbprint = v.(string)
	}

	if v, ok := d.GetOk(constDescription); ok {
		account.Description = v.(string)
	}

	if v, ok := d.GetOk(constEnvironmentIDs); ok {
		account.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constPrivateKeyFile); ok {
		account.PrivateKeyFile = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk(constPrivateKeyPassphrase); ok {
		account.PrivateKeyFile = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk(constResourceManagementEndpointBaseURI); ok {
		account.ResourceManagerEndpoint = v.(string)
	}

	if v, ok := d.GetOk(constSecretKey); ok {
		account.SecretKey = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk(constSpaceID); ok {
		account.SpaceID = v.(string)
	}

	if v, ok := d.GetOk(constSubscriptionID); ok {
		subscriptionID := uuid.MustParse(v.(string))
		account.SubscriptionID = &subscriptionID
	}

	if v, ok := d.GetOk(constTenantedDeploymentParticipation); ok {
		account.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk(constTenantID); ok {
		tenantID := uuid.MustParse(v.(string))
		account.TenantID = &tenantID
	}

	if v, ok := d.GetOk(constTenants); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constTenantTags); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constToken); ok {
		account.Token = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk(constUsername); ok {
		account.Username = v.(string)
	}

	return account
}

func resourceAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	accountResource := buildAccountResource(d)

	client := meta.(*octopusdeploy.Client)
	account, err := client.Accounts.Add(accountResource)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(account.GetID())

	return nil
}

func resourceAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	accountResource := buildAccountResource(d)
	accountResource.ID = d.Id()

	client := meta.(*octopusdeploy.Client)
	account, err := client.Accounts.Update(accountResource)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(account.GetID())

	return nil
}
