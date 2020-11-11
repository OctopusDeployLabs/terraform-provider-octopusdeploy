package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountsRead,
		Schema: map[string]*schema.Schema{
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
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
							Type:     schema.TypeList,
						},
						"tenanted_deployment_participation": {
							Optional: true,
							Type:     schema.TypeString,
						},
						"tenants": {
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
		},
	}
}

func dataSourceAccountsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.AccountsQuery{
		AccountType: octopusdeploy.AccountType(d.Get("account_type").(string)),
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	client := m.(*octopusdeploy.Client)
	accounts, err := client.Accounts.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedAccounts := []interface{}{}
	for _, account := range accounts.Items {
		accountResource, err := octopusdeploy.ToAccountResource(account)
		if err != nil {
			return diag.FromErr(err)
		}

		flattenedAccount := map[string]interface{}{
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
			flattenedAccount["application_id"] = applicationID.String()
		}

		if subscriptionID := accountResource.SubscriptionID; subscriptionID != nil {
			flattenedAccount["subscription_id"] = subscriptionID.String()
		}

		if tenantID := accountResource.TenantID; tenantID != nil {
			flattenedAccount["tenant_id"] = tenantID.String()
		}

		flattenedAccounts = append(flattenedAccounts, flattenedAccount)
	}

	d.Set("accounts", flattenedAccounts)
	d.SetId("Accounts " + time.Now().UTC().String())

	return nil
}
