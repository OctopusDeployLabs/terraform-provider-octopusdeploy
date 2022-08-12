package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandAzureServicePrincipalAccount(d *schema.ResourceData) *accounts.AzureServicePrincipalAccount {
	name := d.Get("name").(string)
	password := d.Get("password").(string)
	secretKey := core.NewSensitiveValue(password)

	applicationID, _ := uuid.Parse(d.Get("application_id").(string))
	tenantID, _ := uuid.Parse(d.Get("tenant_id").(string))
	subscriptionID, _ := uuid.Parse(d.Get("subscription_id").(string))

	account, _ := accounts.NewAzureServicePrincipalAccount(name, subscriptionID, tenantID, applicationID, secretKey)
	account.ID = d.Id()

	if v, ok := d.GetOk("authentication_endpoint"); ok {
		account.AuthenticationEndpoint = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		account.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("environments"); ok {
		account.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("name"); ok {
		account.SetName(v.(string))
	}

	if v, ok := d.GetOk("resource_manager_endpoint"); ok {
		account.ResourceManagerEndpoint = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		account.SetSpaceID(v.(string))
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentMode = core.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenants"); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	return account
}

func getAzureServicePrincipalAccountSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id":                    getApplicationIDSchema(true),
		"authentication_endpoint":           getAuthenticationEndpointSchema(false),
		"azure_environment":                 getAzureEnvironmentSchema(),
		"description":                       getDescriptionSchema("Azure service principal account"),
		"environments":                      getEnvironmentsSchema(),
		"id":                                getIDSchema(),
		"name":                              getNameSchema(true),
		"password":                          getPasswordSchema(true),
		"resource_manager_endpoint":         getResourceManagerEndpointSchema(false),
		"space_id":                          getSpaceIDSchema(),
		"subscription_id":                   getSubscriptionIDSchema(true),
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_id":                         getTenantIDSchema(true),
		"tenant_tags":                       getTenantTagsSchema(),
	}
}

func setAzureServicePrincipalAccount(ctx context.Context, d *schema.ResourceData, account *accounts.AzureServicePrincipalAccount) error {
	d.Set("application_id", account.ApplicationID.String())
	d.Set("authentication_endpoint", account.AuthenticationEndpoint)
	d.Set("azure_environment", account.AzureEnvironment)
	d.Set("description", account.GetDescription())
	d.Set("id", account.GetID())
	d.Set("name", account.GetName())
	d.Set("resource_manager_endpoint", account.ResourceManagerEndpoint)
	d.Set("space_id", account.GetSpaceID())
	d.Set("subscription_id", account.SubscriptionID.String())
	d.Set("tenanted_deployment_participation", account.GetTenantedDeploymentMode())
	d.Set("tenant_id", account.TenantID.String())

	if err := d.Set("environments", account.GetEnvironmentIDs()); err != nil {
		return fmt.Errorf("error setting environments: %s", err)
	}

	if err := d.Set("tenants", account.GetTenantIDs()); err != nil {
		return fmt.Errorf("error setting tenants: %s", err)
	}

	if err := d.Set("tenant_tags", account.TenantTags); err != nil {
		return fmt.Errorf("error setting tenant_tags: %s", err)
	}

	return nil
}
