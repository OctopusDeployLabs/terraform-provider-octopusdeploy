package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataAzureServicePrincipal() *schema.Resource {
	return &schema.Resource{
		Read: dataAzureServicePrincipalReadByName,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
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
			"tenant_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenanted_deployment_participation": getTenantedDeploymentSchema(),
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subscription_number": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"azure_AzureServicePrincipal": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_management_endpoint_base_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"active_directory_endpoint_base_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataAzureServicePrincipalReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	AzureServicePrincipalName := d.Get("name")
	env, err := client.Account.GetByName(AzureServicePrincipalName.(string))

	if err == octopusdeploy.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading AzureServicePrincipal with name %s: %s", AzureServicePrincipalName, err.Error())
	}

	d.SetId(env.ID)

	d.Set("name", env.Name)
	d.Set("description", env.Description)

	return nil
}
