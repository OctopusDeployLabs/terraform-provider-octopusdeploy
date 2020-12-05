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

	if v, ok := d.GetOk("azure_environment"); ok {
		accountResource.AzureEnvironment = v.(string)
	}

	if v, ok := d.GetOk("certificate_data"); ok {
		accountResource.CertificateBytes = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("certificate_thumbprint"); ok {
		accountResource.CertificateThumbprint = v.(string)
	}

	if v, ok := d.GetOk("client_id"); ok {
		clientID := uuid.MustParse(v.(string))
		accountResource.ApplicationID = &clientID
	}

	if v, ok := d.GetOk("client_secret"); ok {
		accountResource.ApplicationPassword = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		accountResource.Description = v.(string)
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

	if v, ok := d.GetOk("resource_management_endpoint_base_uri"); ok {
		accountResource.ResourceManagerEndpoint = v.(string)
	}

	if v, ok := d.GetOk("service_management_endpoint_base_uri"); ok {
		accountResource.ManagementEndpoint = v.(string)
	}

	if v, ok := d.GetOk("service_management_endpoint_suffix"); ok {
		accountResource.StorageEndpointSuffix = v.(string)
	}

	if v, ok := d.GetOk("secret_key"); ok {
		accountResource.SecretKey = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("space_id"); ok {
		accountResource.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("subscription_id"); ok {
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

func flattenAccountResource(accountResource *octopusdeploy.AccountResource) map[string]interface{} {
	flattenedAccountResource := map[string]interface{}{
		"access_key":                        accountResource.AccessKey,
		"account_type":                      accountResource.AccountType,
		"authentication_endpoint":           accountResource.AuthenticationEndpoint,
		"azure_environment":                 accountResource.AzureEnvironment,
		"certificate_thumbprint":            accountResource.CertificateThumbprint,
		"description":                       accountResource.Description,
		"environments":                      accountResource.EnvironmentIDs,
		"id":                                accountResource.GetID(),
		"name":                              accountResource.Name,
		"space_id":                          accountResource.SpaceID,
		"resource_manager_endpoint":         accountResource.ResourceManagerEndpoint,
		"tenant_tags":                       accountResource.TenantTags,
		"tenanted_deployment_participation": accountResource.TenantedDeploymentMode,
		"tenants":                           accountResource.TenantIDs,
		"username":                          accountResource.Username,
	}

	if applicationID := accountResource.ApplicationID; applicationID != nil {
		flattenedAccountResource["application_id"] = applicationID.String()
	}

	if subscriptionID := accountResource.SubscriptionID; subscriptionID != nil {
		flattenedAccountResource["subscription_id"] = subscriptionID.String()
	}

	if tenantID := accountResource.TenantID; tenantID != nil {
		flattenedAccountResource["tenant_id"] = tenantID.String()
	}

	return flattenedAccountResource
}

func getAccountResourceDataSchema() map[string]*schema.Schema {
	dataSchema := getAccountResourceSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"account_type": getQueryAccountType(),
		"accounts": {
			Computed:    true,
			Description: "A list of accounts that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Type:        schema.TypeList,
		},
		"id":           getDataSchemaID(),
		"ids":          getQueryIDs(),
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
	}
}

func getAccountResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_key":   getAccessKeySchema(false),
		"account_type": getAccountTypeSchema(),
		"active_directory_endpoint_base_uri": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"application_id":          getApplicationIDSchema(false),
		"authentication_endpoint": getAuthenticationEndpointSchema(false),
		"azure_environment":       getAzureEnvironmentSchema(),
		"certificate_data": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"certificate_thumbprint": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"client_id": {
			Optional:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
		},
		"client_secret": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"description":               getDescriptionSchema(),
		"environments":              getEnvironmentsSchema(),
		"id":                        getIDSchema(),
		"name":                      getNameSchema(true),
		"password":                  getPasswordSchema(false),
		"resource_manager_endpoint": getResourceManagerEndpointSchema(false),
		"private_key_file": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"private_key_passphrase": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"secret_key": getSecretKeySchema(false),
		"service_management_endpoint_base_uri": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"service_management_endpoint_suffix": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"space_id":                          getSpaceIDSchema(),
		"subscription_id":                   getSubscriptionIDSchema(false),
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_id":                         getTenantIDSchema(false),
		"tenant_tags":                       getTenantTagsSchema(),
		"token":                             getTokenSchema(false),
		"username":                          getUsernameSchema(false),
	}
}

func setAccountResource(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.AccountResource) {
	d.Set("access_key", account.AccessKey)
	d.Set("account_type", account.GetAccountType())
	d.Set("active_directory_endpoint_base_uri", account.AuthenticationEndpoint)
	d.Set("azure_environment", account.AzureEnvironment)
	d.Set("certificate_thumbprint", account.CertificateThumbprint)
	d.Set("description", account.GetDescription())
	d.Set("environments", account.GetEnvironmentIDs())
	d.Set("name", account.GetName())
	d.Set("resource_manager_endpoint", account.ResourceManagerEndpoint)
	d.Set("secret_key", account.SecretKey)
	d.Set("service_management_endpoint_base_uri", account.ManagementEndpoint)
	d.Set("service_management_endpoint_suffix", account.StorageEndpointSuffix)
	d.Set("space_id", account.GetSpaceID())
	d.Set("tenanted_deployment_participation", account.GetTenantedDeploymentMode())
	d.Set("tenants", account.GetTenantIDs())
	d.Set("tenant_tags", account.GetTenantTags())
	d.Set("username", account.Username)

	if account.ApplicationID != nil {
		d.Set("client_id", account.ApplicationID.String())
	}

	if account.SubscriptionID != nil {
		d.Set("subscription_id", account.SubscriptionID.String())
	}

	if account.TenantID != nil {
		d.Set("tenant_id", account.TenantID.String())
	}

	d.SetId(account.GetID())
}
