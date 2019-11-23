package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAzureServicePrincipal() *schema.Resource {
	return &schema.Resource{
		Create: resourceAzureServicePrincipalCreate,
		Read:   resourceAzureServicePrincipalRead,
		Update: resourceAzureServicePrincipalUpdate,
		Delete: resourceAzureServicePrincipalDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environments": {
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
			"azure_environment": {
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

func resourceAzureServicePrincipalRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountID := d.Id()
	account, err := client.Account.Get(accountID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading azure service principal %s: %s", accountID, err.Error())
	}

	d.Set("name", account.Name)
	d.Set("description", account.Description)
	d.Set("environments", account.EnvironmentIDs)
	d.Set("tenanted_deployment_participation", account.TenantedDeploymentParticipation.String())
	d.Set("tenant_tags", account.TenantTags)

	d.Set("client_id", account.ClientID)
	d.Set("tenant_id", account.TenantID)
	d.Set("subscription_number", account.SubscriptionNumber)
	d.Set("key", nil)
	d.Set("azure_environment", account.AzureEnvironment)
	d.Set("resource_management_endpoint_base_uri", account.ResourceManagementEndpointBaseURI)
	d.Set("active_directory_endpoint_base_uri", account.ActiveDirectoryEndpointBaseURI)

	return nil
}

func buildAzureServicePrincipalResource(d *schema.ResourceData) *octopusdeploy.Account {
	var account = octopusdeploy.NewAccount(d.Get("name").(string), octopusdeploy.AzureServicePrincipal)

	// Required fields
	account.ClientID = d.Get("client_id").(string)
	account.TenantID = d.Get("tenant_id").(string)
	account.SubscriptionNumber = d.Get("subscription_number").(string)
	account.Password = octopusdeploy.SensitiveValue{NewValue: d.Get("key").(string)}

	// Optional Fields
	if v, ok := d.GetOk("description"); ok {
		account.Description = v.(string)
	}

	if v, ok := d.GetOk("environments"); ok {
		account.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentParticipation, _ = octopusdeploy.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("resource_management_endpoint_base_uri"); ok {
		account.ResourceManagementEndpointBaseURI = v.(string)
	}

	if v, ok := d.GetOk("active_directory_endpoint_base_uri"); ok {
		account.ActiveDirectoryEndpointBaseURI = v.(string)
	}

	return account
}

func resourceAzureServicePrincipalCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newAccount := buildAzureServicePrincipalResource(d)
	account, err := client.Account.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error creating azure service principal %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceAzureServicePrincipalUpdate(d *schema.ResourceData, m interface{}) error {
	account := buildAzureServicePrincipalResource(d)
	account.ID = d.Id()

	client := m.(*octopusdeploy.Client)

	updatedAccount, err := client.Account.Update(account)

	if err != nil {
		return fmt.Errorf("error updating azure service principal id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceAzureServicePrincipalDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountID := d.Id()

	err := client.Account.Delete(accountID)

	if err != nil {
		return fmt.Errorf("error deleting azure service principal id %s: %s", accountID, err.Error())
	}

	d.SetId("")
	return nil
}
