package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenAccountResource(accountResource *accounts.AccountResource) map[string]interface{} {
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
		"space_id":     getQuerySpaceID(),
		"ids":          getQueryIDs(),
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
	}
}

func getAccountResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_key":   getAccessKeySchema(false),
		"account_type": getAccountTypeSchema(true),
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
		"client_secret": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"description":               getDescriptionSchema("account resource"),
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
