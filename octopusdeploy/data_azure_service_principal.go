package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	name := d.Get(constName).(string)

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.GetByName(name)
	if err != nil {
		return createResourceOperationError(errorReadingAzureServicePrincipal, name, err)
	}
	if resource == nil {
		return nil
	}

	logResource(constAccount, m)

	d.SetId(resource.ID)
	d.Set(constName, resource.Name)
	d.Set(constDescription, resource.Description)

	return nil
}
