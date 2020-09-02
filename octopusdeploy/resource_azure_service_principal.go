package model

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	uuid "github.com/google/uuid"
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
			"account_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Azure",
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

func buildAzureServicePrincipalResource(d *schema.ResourceData) *model.Account {
	account, err := model.NewAccount(d.Get("name").(string), enum.AzureServicePrincipal)
	if err != nil {
		return nil
	}

	clientID, err := uuid.Parse(d.Get("client_id").(string))
	tenantID, err := uuid.Parse(d.Get("tenant_id").(string))
	subscriptionNumber, err := uuid.Parse(d.Get("subscription_number").(string))

	// Required fields
	account.ClientID = &clientID
	account.TenantID = &tenantID
	account.SubscriptionNumber = &subscriptionNumber
	password := d.Get("key").(string)
	account.Password = &model.SensitiveValue{NewValue: &password}

	// Optional Fields
	if v, ok := d.GetOk("description"); ok {
		account.Description = v.(string)
	}

	if v, ok := d.GetOk("environments"); ok {
		account.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentParticipation, _ = enum.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("resource_management_endpoint_base_uri"); ok {
		account.ResourceManagementEndpointBase = v.(string)
	}

	if v, ok := d.GetOk("active_directory_endpoint_base_uri"); ok {
		account.ActiveDirectoryEndpointBase = v.(string)
	}

	return account
}

func resourceAzureServicePrincipalCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	newAccount := buildAzureServicePrincipalResource(d)
	account, err := apiClient.Accounts.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error creating azure service principal %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceAzureServicePrincipalRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	accountID := d.Id()
	account, err := apiClient.Accounts.Get(accountID)

	if err == client.ErrItemNotFound {
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
	d.Set("key", account.Password)
	d.Set("azure_environment", account.AzureEnvironment)
	d.Set("resource_management_endpoint_base_uri", account.ResourceManagementEndpointBase)
	d.Set("active_directory_endpoint_base_uri", account.ActiveDirectoryEndpointBase)

	return nil
}

func resourceAzureServicePrincipalUpdate(d *schema.ResourceData, m interface{}) error {
	account := buildAzureServicePrincipalResource(d)
	account.ID = d.Id()

	apiClient := m.(*client.Client)

	updatedAccount, err := apiClient.Accounts.Update(account)

	if err != nil {
		return fmt.Errorf("error updating azure service principal id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceAzureServicePrincipalDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	accountID := d.Id()

	err := apiClient.Accounts.Delete(accountID)

	if err != nil {
		return fmt.Errorf("error deleting azure service principal id %s: %s", accountID, err.Error())
	}

	d.SetId("")
	return nil
}
