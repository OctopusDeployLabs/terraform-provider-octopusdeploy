package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandAzureSubscriptionAccount(d *schema.ResourceData) *octopusdeploy.AzureSubscriptionAccount {
	name := d.Get("name").(string)
	subscriptionID, _ := uuid.Parse(d.Get("subscription_id").(string))

	account, _ := octopusdeploy.NewAzureSubscriptionAccount(name, subscriptionID)
	account.ID = d.Id()

	if v, ok := d.GetOk("azure_environment"); ok {
		account.AzureEnvironment = v.(string)
	}

	if v, ok := d.GetOk("certificate"); ok {
		account.CertificateBytes = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("certificate_thumbprint"); ok {
		account.CertificateThumbprint = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		account.Description = v.(string)
	}

	if v, ok := d.GetOk("environments"); ok {
		account.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("management_endpoint"); ok {
		account.ManagementEndpoint = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		account.Name = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		account.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("storage_endpoint_suffix"); ok {
		account.StorageEndpointSuffix = v.(string)
	}

	if v, ok := d.GetOk("subscription_id"); ok {
		subscriptionID, _ := uuid.Parse(v.(string))
		account.SubscriptionID = &subscriptionID
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

func getAzureSubscriptionAccountSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"azure_environment": getAzureEnvironmentSchema(),
		"certificate": {
			Computed:  true,
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"certificate_thumbprint": {
			Computed:  true,
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"description":  getDescriptionSchema("Azure subscription account"),
		"environments": getEnvironmentsSchema(),
		"management_endpoint": {
			Required:     true,
			Type:         schema.TypeString,
			RequiredWith: []string{"azure_environment"},
		},
		"name":     getNameSchema(true),
		"space_id": getSpaceIDSchema(),
		"storage_endpoint_suffix": {
			Description:  "The storage endpoint suffix associated with this Azure subscription account.",
			Required:     true,
			Type:         schema.TypeString,
			RequiredWith: []string{"azure_environment"},
		},
		"subscription_id":                   getSubscriptionIDSchema(true),
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_tags":                       getTenantTagsSchema(),
	}
}

func setAzureSubscriptionAccount(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.AzureSubscriptionAccount) error {
	d.Set("azure_environment", account.AzureEnvironment)
	d.Set("certificate_thumbprint", account.CertificateThumbprint)
	d.Set("description", account.GetDescription())
	d.Set("management_endpoint", account.ManagementEndpoint)
	d.Set("name", account.GetName())
	d.Set("space_id", account.GetSpaceID())
	d.Set("storage_endpoint_suffix", account.StorageEndpointSuffix)
	d.Set("subscription_id", account.SubscriptionID.String())
	d.Set("tenanted_deployment_participation", account.GetTenantedDeploymentMode())

	if err := d.Set("environments", account.GetEnvironmentIDs()); err != nil {
		return fmt.Errorf("error setting environments: %s", err)
	}

	if err := d.Set("tenants", account.GetTenantIDs()); err != nil {
		return fmt.Errorf("error setting tenants: %s", err)
	}

	if err := d.Set("tenant_tags", account.TenantTags); err != nil {
		return fmt.Errorf("error setting tenant_tags: %s", err)
	}

	d.SetId(account.GetID())

	return nil
}
