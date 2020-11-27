package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenAccount(accountResource *octopusdeploy.AccountResource) map[string]interface{} {
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
		"resource_management_endpoint":      accountResource.ResourceManagerEndpoint,
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

func getAccountDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_type": {
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validateValueFunc([]string{
				"AmazonWebServicesAccount",
				"AzureServicePrincipal",
				"AzureSubscription",
				"SshKeyPair",
				"Token",
				"UsernamePassword",
			}),
		},
		"ids": {
			Description: "Query and/or search by a list of IDs",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"partial_name": {
			Description: "Query and/or search by partial name",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"skip": {
			Default:     0,
			Description: "Indicates the number of items to skip in the response",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"take": {
			Default:     1,
			Description: "Indicates the number of items to take (or return) in the response",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"accounts": {
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"access_key": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"account_type": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"application_id": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"authentication_endpoint": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"azure_environment": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"certificate_thumbprint": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"description": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"environments": {
						Elem:     &schema.Schema{Type: schema.TypeString},
						Optional: true,
						Type:     schema.TypeList,
					},
					"id": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"name": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"space_id": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"resource_management_endpoint": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"subscription_id": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"tenant_id": {
						Optional: true,
						Type:     schema.TypeString,
					},
					"tenant_tags": {
						Elem:     &schema.Schema{Type: schema.TypeString},
						Optional: true,
						Type:     schema.TypeList,
					},
					"tenanted_deployment_participation": getTenantedDeploymentSchema(),
					"tenants": {
						Elem:     &schema.Schema{Type: schema.TypeString},
						Optional: true,
						Type:     schema.TypeList,
					},
					"username": {
						Optional: true,
						Type:     schema.TypeString,
					},
				},
			},
			Type: schema.TypeList,
		},
	}
}
