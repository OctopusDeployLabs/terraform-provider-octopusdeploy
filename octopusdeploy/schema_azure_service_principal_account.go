package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandAzureServicePrincipalAccount(d *schema.ResourceData) *octopusdeploy.AzureServicePrincipalAccount {
	name := d.Get("name").(string)
	password := d.Get("application_password").(string)
	secretKey := octopusdeploy.NewSensitiveValue(password)

	applicationID, _ := uuid.Parse(d.Get("application_id").(string))
	tenantID, _ := uuid.Parse(d.Get("tenant_id").(string))
	subscriptionID, _ := uuid.Parse(d.Get("subscription_id").(string))

	account, _ := octopusdeploy.NewAzureServicePrincipalAccount(name, subscriptionID, tenantID, applicationID, secretKey)
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
		account.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenants"); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	return account
}

func setAzureServicePrincipalAccount(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.AzureServicePrincipalAccount) {
	setAccount(ctx, d, account)

	d.Set("account_type", "AzureServicePrincipal")
	d.Set("application_id", account.ApplicationID.String())

	if account.ApplicationPassword != nil {
		d.Set("application_password", account.ApplicationPassword.NewValue)
	}

	d.Set("authentication_endpoint", account.AuthenticationEndpoint)
	d.Set("azure_environment", account.AzureEnvironment)
	d.Set("resource_manager_endpoint", account.ResourceManagerEndpoint)
	d.Set("subscription_id", account.SubscriptionID.String())
	d.Set("tenant_id", account.TenantID.String())

	d.SetId(account.GetID())
}

func getAzureServicePrincipalAccountSchema() map[string]*schema.Schema {
	schemaMap := getAccountSchema()
	schemaMap["account_type"] = &schema.Schema{
		Optional: true,
		Default:  "AzureServicePrincipal",
		Type:     schema.TypeString,
	}
	schemaMap["application_id"] = &schema.Schema{
		Required:         true,
		Type:             schema.TypeString,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	schemaMap["application_password"] = &schema.Schema{
		Required:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	schemaMap["authentication_endpoint"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	schemaMap["azure_environment"] = getAzureEnvironmentSchema()
	schemaMap["resource_manager_endpoint"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	schemaMap["subscription_id"] = &schema.Schema{
		Required:         true,
		Type:             schema.TypeString,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	schemaMap["tenant_id"] = &schema.Schema{
		Required:         true,
		Type:             schema.TypeString,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	return schemaMap
}
