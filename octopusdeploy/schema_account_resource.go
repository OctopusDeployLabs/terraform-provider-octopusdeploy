package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandAccountResource(d *schema.ResourceData) *octopusdeploy.AccountResource {
	name := d.Get("name").(string)

	var accountType octopusdeploy.AccountType
	if v, ok := d.GetOk("account_type"); ok {
		accountType = octopusdeploy.AccountType(v.(string))
	}

	accountResource := octopusdeploy.NewAccountResource(name, accountType)
	accountResource.ID = d.Id()

	if v, ok := d.GetOk("access_key"); ok {
		accountResource.AccessKey = v.(string)
	}

	if v, ok := d.GetOk("active_directory_endpoint_base_uri"); ok {
		accountResource.AuthenticationEndpoint = v.(string)
	}

	if v, ok := d.GetOk("client_id"); ok {
		clientID := uuid.MustParse(v.(string))
		accountResource.ApplicationID = &clientID
	}

	if v, ok := d.GetOk("client_secret"); ok {
		accountResource.ApplicationPassword = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("azure_environment"); ok {
		accountResource.AzureEnvironment = v.(string)
	}

	if v, ok := d.GetOk("certificate_data"); ok {
		accountResource.CertificateBytes = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("certificate_thumbprint"); ok {
		accountResource.CertificateThumbprint = v.(string)
	}

	if v, ok := d.GetOk("environments"); ok {
		accountResource.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("password"); ok {
		accountResource.ApplicationPassword = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("private_key_file"); ok {
		accountResource.PrivateKeyFile = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("private_key_passphrase"); ok {
		accountResource.PrivateKeyFile = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("service_management_endpoint_suffix"); ok {
		accountResource.StorageEndpointSuffix = v.(string)
	}

	if v, ok := d.GetOk("resource_management_endpoint_base_uri"); ok {
		accountResource.ResourceManagerEndpoint = v.(string)
	}

	if v, ok := d.GetOk("secret_key"); ok {
		accountResource.SecretKey = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("space_id"); ok {
		accountResource.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("subscription_number"); ok {
		subscriptionID := uuid.MustParse(v.(string))
		accountResource.SubscriptionID = &subscriptionID
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		accountResource.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		tenantID := uuid.MustParse(v.(string))
		accountResource.TenantID = &tenantID
	}

	if v, ok := d.GetOk("tenants"); ok {
		accountResource.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		accountResource.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("token"); ok {
		accountResource.Token = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("username"); ok {
		accountResource.Username = v.(string)
	}

	return accountResource
}

func flattenAccountResource(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.AccountResource) {
	flattenAccount(ctx, d, account)

	d.Set("access_key", account.AccessKey)
	d.Set("active_directory_endpoint_base_uri", account.AuthenticationEndpoint)
	d.Set("azure_environment", account.AzureEnvironment)
	d.Set("certificate_bytes", account.CertificateBytes)
	d.Set("certificate_thumbprint", account.CertificateThumbprint)
	d.Set("client_id", account.ApplicationID.String())
	d.Set("client_secret", account.ApplicationPassword)
	d.Set("service_management_endpoint_base_uri", account.ManagementEndpoint)
	d.Set("password", account.ApplicationPassword)
	d.Set("private_key_file", account.PrivateKeyFile)
	d.Set("private_key_passphrase", account.PrivateKeyPassphrase)
	d.Set("resource_management_endpoint_base_uri", account.ResourceManagerEndpoint)
	d.Set("service_management_endpoint_suffix", account.StorageEndpointSuffix)
	d.Set("secret_key", account.SecretKey)
	d.Set("subscription_number", account.SubscriptionID.String())
	d.Set("tenant_id", account.TenantID.String())
	d.Set("token", account.Token.NewValue)
	d.Set("username", account.Username)

	d.SetId(account.GetID())
}

func getAccountResourceSchema() map[string]*schema.Schema {
	schemaMap := getAccountSchema()
	schemaMap["access_key"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	schemaMap["active_directory_endpoint_base_uri"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	schemaMap["azure_environment"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	schemaMap["certificate_bytes"] = &schema.Schema{
		Optional:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	schemaMap["certificate_thumbprint"] = &schema.Schema{
		Optional:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	schemaMap["client_id"] = &schema.Schema{
		Optional:         true,
		Type:             schema.TypeString,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	schemaMap["password"] = &schema.Schema{
		Optional:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	schemaMap["resource_management_endpoint_base_uri"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	schemaMap["secret_key"] = &schema.Schema{
		Optional:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	schemaMap["subscription_number"] = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	schemaMap["tenant_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return schemaMap
}
