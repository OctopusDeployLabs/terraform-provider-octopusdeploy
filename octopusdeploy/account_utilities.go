package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func buildAccountResource(d *schema.ResourceData) *octopusdeploy.AccountResource {
	var name string
	if v, ok := d.GetOk(constName); ok {
		name = v.(string)
	}

	var accountType octopusdeploy.AccountType
	if v, ok := d.GetOk(constAccountType); ok {
		accountType = octopusdeploy.AccountType(v.(string))
	}

	account := octopusdeploy.NewAccountResource(name, accountType)

	if v, ok := d.GetOk(constAccessKey); ok {
		account.AccessKey = v.(string)
	}

	if v, ok := d.GetOk(constClientID); ok {
		clientID := uuid.MustParse(v.(string))
		account.ApplicationID = &clientID
	}

	if v, ok := d.GetOk(constClientSecret); ok {
		account.ApplicationPassword = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk(constActiveDirectoryEndpointBaseURI); ok {
		account.AuthenticationEndpoint = v.(string)
	}

	if v, ok := d.GetOk(constAzureEnvironment); ok {
		account.AzureEnvironment = v.(string)
	}

	if v, ok := d.GetOk(constCertificateData); ok {
		account.CertificateBytes = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk(constCertificateThumbprint); ok {
		account.CertificateThumbprint = v.(string)
	}

	if v, ok := d.GetOk(constDescription); ok {
		account.Description = v.(string)
	}

	if v, ok := d.GetOk(constEnvironmentIDs); ok {
		account.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constPrivateKeyFile); ok {
		account.PrivateKeyFile = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk(constPrivateKeyPassphrase); ok {
		account.PrivateKeyFile = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk(constResourceManagementEndpointBaseURI); ok {
		account.ResourceManagerEndpoint = v.(string)
	}

	if v, ok := d.GetOk(constSecretKey); ok {
		account.SecretKey = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk(constSpaceID); ok {
		account.SpaceID = v.(string)
	}

	if v, ok := d.GetOk(constSubscriptionID); ok {
		subscriptionID := uuid.MustParse(v.(string))
		account.SubscriptionID = &subscriptionID
	}

	if v, ok := d.GetOk(constTenantedDeploymentParticipation); ok {
		account.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk(constTenantID); ok {
		tenantID := uuid.MustParse(v.(string))
		account.TenantID = &tenantID
	}

	if v, ok := d.GetOk(constTenants); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constTenantTags); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constToken); ok {
		account.Token = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk(constUsername); ok {
		account.Username = v.(string)
	}

	return account
}

func flattenAWSAccount(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.AmazonWebServicesAccount) {
	d.Set(constAccessKey, account.AccessKey)
	d.Set(constDescription, account.Description)
	d.Set(constEnvironments, account.EnvironmentIDs)
	d.Set(constName, account.Name)
	// d.Set(constSecretKey, account.SecretKey)
	d.Set(constSpaceID, account.SpaceID)
	d.Set(constTenantedDeploymentParticipation, account.TenantedDeploymentMode)
	d.Set(constTenants, account.TenantIDs)
	d.Set(constTenantTags, account.TenantTags)
	d.SetId(account.GetID())
}

func flattenAccountResource(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.AccountResource) {
	d.Set(constAccountType, account.AccountType)
	d.Set(constAccessKey, account.AccessKey)
	d.Set(constActiveDirectoryEndpointBaseURI, account.AuthenticationEndpoint)
	d.Set(constAzureEnvironment, account.AzureEnvironment)
	d.Set(constClientID, account.ApplicationID.String())
	d.Set(constDescription, account.Description)
	d.Set(constEnvironments, account.EnvironmentIDs)
	d.Set(constName, account.Name)
	// d.Set(constPassphrase, account.PrivateKeyPassphrase)
	// d.Set(constPassword, account.ApplicationPassword)
	d.Set(constResourceManagementEndpointBaseURI, account.ResourceManagerEndpoint)
	// d.Set(constSecretKey, account.SecretKey)
	d.Set(constSpaceID, account.SpaceID)
	d.Set(constSubscriptionID, account.SubscriptionID.String())
	d.Set(constTenantedDeploymentParticipation, account.TenantedDeploymentMode)
	d.Set(constTenantID, account.TenantID.String())
	d.Set(constTenants, account.TenantIDs)
	d.Set(constTenantTags, account.TenantTags)
	// d.Set(constToken, account.Token)
	d.Set(constUsername, account.Username)
	d.SetId(account.GetID())
}

func getCommonAccountsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		constDescription: {
			Optional: true,
			Type:     schema.TypeString,
		},
		constEnvironments: {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		constName: {
			Required: true,
			Type:     schema.TypeString,
		},
		constSpaceID: {
			Computed: true,
			Type:     schema.TypeString,
		},
		constTenantedDeploymentParticipation: getTenantedDeploymentSchema(),
		constTenants: {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		constTenantTags: {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}
