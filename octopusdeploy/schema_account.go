package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"partial_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"skip": {
			Default:  0,
			Type:     schema.TypeInt,
			Optional: true,
		},
		"take": {
			Default:  1,
			Type:     schema.TypeInt,
			Optional: true,
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

func getAccountSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_type": {
			Default:  "None",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validateValueFunc([]string{
				"AmazonWebServicesAccount",
				"AzureServicePrincipal",
				"AzureSubscription",
				"None",
				"SshKeyPair",
				"Token",
				"UsernamePassword",
			}),
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
			Computed: true,
			Type:     schema.TypeString,
		},
		"name": &schema.Schema{
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"space_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"tenant_tags": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}

func setAccount(ctx context.Context, d *schema.ResourceData, account octopusdeploy.IAccount) {
	d.Set("account_type", account.GetAccountType())
	d.Set("description", account.GetDescription())
	d.Set("environments", account.GetEnvironmentIDs())
	d.Set("name", account.GetName())
	d.Set("space_id", account.GetSpaceID())
	d.Set("tenanted_deployment_participation", account.GetTenantedDeploymentMode())
	d.Set("tenants", account.GetTenantIDs())
	d.Set("tenant_tags", account.GetTenantTags())
}
