package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAccountReadByName,
		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func dataAccountReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	name := d.Get("name").(string)
	query := octopusdeploy.AccountsQuery{
		PartialName: name,
		Take:        1,
	}

	accounts, err := client.Accounts.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}
	if accounts == nil || len(accounts.Items) == 0 {
		d.SetId("")
		return diag.Errorf("unable to retrieve account (partial name: %s)", name)
	}

	logResource(constAccount, m)

	accountResource := accounts.Items[0].(*octopusdeploy.AccountResource)

	d.SetId(accountResource.GetID())
	d.Set(constAccessKey, accountResource.AccessKey)
	d.Set(constAccountType, accountResource.AccountType)
	d.Set(constClientID, accountResource.ApplicationID)
	d.Set(constActiveDirectoryEndpointBaseURI, accountResource.AuthenticationEndpoint)
	d.Set(constAzureEnvironment, accountResource.AzureEnvironment)
	d.Set(constCertificateData, accountResource.CertificateBytes)
	d.Set(constCertificateThumbprint, accountResource.CertificateThumbprint)
	d.Set("description", accountResource.Description)
	d.Set(constEnvironments, accountResource.EnvironmentIDs)
	d.Set(constName, accountResource.GetName())
	d.Set(constSpaceID, accountResource.SpaceID)
	d.Set(constResourceManagementEndpointBaseURI, accountResource.ResourceManagerEndpoint)
	d.Set(constSubscriptionNumber, accountResource.SubscriptionID)
	d.Set(constTenants, accountResource.TenantIDs)
	d.Set("tenant_tags", accountResource.TenantTags)
	d.Set("tenanted_deployment_participation", accountResource.TenantedDeploymentMode)
	d.Set(constToken, accountResource.Token)
	d.Set(constUsername, accountResource.Username)

	return nil
}
