package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAzureServicePrincipal() *schema.Resource {
	return &schema.Resource{
		Read: dataAzureServicePrincipalReadByName,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			"AzureServicePrincipals": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			constTenantTags: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constTenantedDeploymentParticipation: getTenantedDeploymentSchema(),
			constClientID: {
				Type:     schema.TypeString,
				Required: true,
			},
			constTenantID: {
				Type:     schema.TypeString,
				Required: true,
			},
			constSubscriptionNumber: {
				Type:     schema.TypeString,
				Required: true,
			},
			constKey: {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"azure_AzureServicePrincipal": {
				Type:     schema.TypeString,
				Optional: true,
			},
			constResourceManagementEndpointBaseURI: {
				Type:     schema.TypeString,
				Optional: true,
			},
			constActiveDirectoryEndpointBaseURI: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataAzureServicePrincipalReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)
	name := d.Get(constName).(string)
	query := octopusdeploy.AccountsQuery{
		PartialName: name,
		Take:        1,
	}

	accounts, err := client.Accounts.Get(query)
	if err != nil {
		return createResourceOperationError(errorReadingAzureServicePrincipal, name, err)
	}
	if accounts == nil || len(accounts.Items) == 0 {
		return fmt.Errorf("Unabled to retrieve account (partial name: %s)", name)
	}

	logResource(constAccount, m)

	account := accounts.Items[0].(*octopusdeploy.AzureServicePrincipalAccount)

	d.SetId(account.GetID())
	d.Set(constName, account.GetName())
	d.Set(constDescription, account.Description)

	return nil
}
