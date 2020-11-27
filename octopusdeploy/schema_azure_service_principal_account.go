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

func getAzureServicePrincipalAccountSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
		},
		"application_password": {
			Required:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"authentication_endpoint": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"azure_environment": getAzureEnvironmentSchema(),
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
		"name": {
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_manager_endpoint": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"space_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"subscription_id": {
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
		},
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"tenant_id": {
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
		},
		"tenant_tags": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}

func setAzureServicePrincipalAccount(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.AzureServicePrincipalAccount) {
	d.Set("application_id", account.ApplicationID.String())
	d.Set("authentication_endpoint", account.AuthenticationEndpoint)
	d.Set("azure_environment", account.AzureEnvironment)
	d.Set("description", account.GetDescription())
	d.Set("environments", account.GetEnvironmentIDs())
	d.Set("id", account.GetID())
	d.Set("name", account.GetName())
	d.Set("resource_manager_endpoint", account.ResourceManagerEndpoint)
	d.Set("space_id", account.GetSpaceID())
	d.Set("subscription_id", account.SubscriptionID.String())
	d.Set("tenanted_deployment_participation", account.GetTenantedDeploymentMode())
	d.Set("tenants", account.GetTenantIDs())
	d.Set("tenant_id", account.TenantID.String())
	d.Set("tenant_tags", account.GetTenantTags())

	d.SetId(account.GetID())
}
