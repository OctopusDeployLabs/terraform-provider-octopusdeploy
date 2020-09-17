package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAzureServicePrincipal() *schema.Resource {
	schemaMap := getCommonAccountsSchema()

	schemaMap["client_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	schemaMap["tenant_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	schemaMap["subscription_number"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	schemaMap["key"] = &schema.Schema{
		Type:      schema.TypeString,
		Required:  true,
		Sensitive: true,
	}
	schemaMap["azure_environment"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap["resource_management_endpoint_base_uri"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap["active_directory_endpoint_base_uri"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	return &schema.Resource{
		Create: resourceAzureServicePrincipalCreate,
		Read:   resourceAzureServicePrincipalRead,
		Update: resourceAzureServicePrincipalUpdate,
		Delete: resourceAccountDeleteCommon,
		Schema: schemaMap,
	}
}

func buildAzureServicePrincipalResource(d *schema.ResourceData) (*model.Account, error) {
	if d == nil {
		return nil, createInvalidParameterError("buildAzureServicePrincipalResource", "d")
	}

	name := d.Get("name").(string)

	password := d.Get("key").(string)
	if password == "" {
		log.Println("Key is nil. Must add in a password")
	}

	secretKey := model.NewSensitiveValue(password)

	applicationID, err := uuid.Parse(d.Get("client_id").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	tenantID, err := uuid.Parse(d.Get("tenant_id").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	subscriptionID, err := uuid.Parse(d.Get("subscription_number").(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	account, err := model.NewAzureServicePrincipalAccount(name, subscriptionID, tenantID, applicationID, secretKey)
	if err != nil {
		log.Println(err)
		return nil, err
	}

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

	if v, ok := d.GetOk("tenants"); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("resource_management_endpoint_base_uri"); ok {
		account.ResourceManagementEndpointBase = v.(string)
	}

	if v, ok := d.GetOk("active_directory_endpoint_base_uri"); ok {
		account.ActiveDirectoryEndpointBase = v.(string)
	}

	return account, nil
}

func resourceAzureServicePrincipalCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	newAccount, err := buildAzureServicePrincipalResource(d)
	if err != nil {
		log.Println(err)
		return err
	}
	account, err := apiClient.Accounts.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error creating azure service principal %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceAzureServicePrincipalRead(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceAzureServicePrincipalRead", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceAzureServicePrincipalRead", "m")
	}

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

	d.Set("client_id", account.ApplicationID)
	d.Set("tenant_id", account.TenantIDs)
	d.Set("subscription_number", account.SubscriptionID)
	d.Set("key", account.Password)
	d.Set("azure_environment", account.AzureEnvironment)
	d.Set("resource_management_endpoint_base_uri", account.ResourceManagementEndpointBase)
	d.Set("active_directory_endpoint_base_uri", account.ActiveDirectoryEndpointBase)

	return nil
}

func resourceAzureServicePrincipalUpdate(d *schema.ResourceData, m interface{}) error {
	account, err := buildAzureServicePrincipalResource(d)
	if err != nil {
		return err
	}

	account.ID = d.Id()

	apiClient := m.(*client.Client)

	updatedAccount, err := apiClient.Accounts.Update(*account)

	if err != nil {
		return fmt.Errorf("error updating azure service principal id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}
