package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAzureServicePrincipal() *schema.Resource {
	validateSchema()

	log.Println("Hello")
	schemaMap := getCommonAccountsSchema()

	schemaMap[constClientID] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	schemaMap[constTenantID] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	schemaMap[constSubscriptionNumber] = &schema.Schema{
		Type: schema.TypeString,
		//Computed:     true,
		Required:     true,
		ValidateFunc: validation.IsUUID,
	}
	schemaMap[constKey] = &schema.Schema{
		Type:      schema.TypeString,
		Required:  true,
		Sensitive: true,
	}
	schemaMap[constAzureEnvironment] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap[constResourceManagementEndpointBaseURI] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap[constActiveDirectoryEndpointBaseURI] = &schema.Schema{
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

	name := d.Get(constName).(string)

	password := d.Get(constKey).(string)
	if password == constEmptyString {
		log.Println("Key is nil. Must add in a password")
	}

	secretKey := model.NewSensitiveValue(password)

	applicationID, err := uuid.Parse(d.Get(constClientID).(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	tenantID, err := uuid.Parse(d.Get(constTenantID).(string))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	subscriptionID, err := uuid.Parse(d.Get(constSubscriptionNumber).(string))
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
	if v, ok := d.GetOk(constDescription); ok {
		account.Description = v.(string)
	}

	if v, ok := d.GetOk(constEnvironments); ok {
		account.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constTenantedDeploymentParticipation); ok {
		account.TenantedDeploymentParticipation, _ = enum.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk(constTenantTags); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constTenants); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constResourceManagementEndpointBaseURI); ok {
		account.ResourceManagementEndpointBase = v.(string)
	}

	if v, ok := d.GetOk(constActiveDirectoryEndpointBaseURI); ok {
		account.ActiveDirectoryEndpointBase = v.(string)
	}

	err = account.Validate()
	if err != nil {
		return nil, err
	}

	return account, nil
}

func resourceAzureServicePrincipalCreate(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceAzureServicePrincipalRead", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceAzureServicePrincipalRead", "m")
	}

	apiClient := m.(*client.Client)

	newAccount, err := buildAzureServicePrincipalResource(d)
	if err != nil {
		log.Println(err)
		return err
	}
	account, err := apiClient.Accounts.Add(newAccount)

	if err != nil {
		createResourceOperationError(errorCreatingAzureServicePrincipal, newAccount.Name, err)
	}

	if account.ID == constEmptyString {
		log.Println("ID is nil")
	} else {
		d.SetId(account.ID)
	}

	return nil
}

func resourceAzureServicePrincipalRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingAzureServicePrincipal, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constAccount, m)

	d.Set(constName, resource.Name)
	d.Set(constDescription, resource.Description)
	d.Set(constEnvironments, resource.EnvironmentIDs)
	d.Set(constTenantedDeploymentParticipation, resource.TenantedDeploymentParticipation.String())
	d.Set(constTenantTags, resource.TenantTags)
	d.Set(constClientID, resource.ApplicationID)
	d.Set(constTenantID, resource.TenantIDs)
	d.Set(constSubscriptionNumber, resource.SubscriptionID)
	d.Set(constKey, resource.Password)
	d.Set(constAzureEnvironment, resource.AzureEnvironment)
	d.Set(constResourceManagementEndpointBaseURI, resource.ResourceManagementEndpointBase)
	d.Set(constActiveDirectoryEndpointBaseURI, resource.ActiveDirectoryEndpointBase)

	return nil
}

func resourceAzureServicePrincipalUpdate(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceAzureServicePrincipalRead", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceAzureServicePrincipalRead", "m")
	}

	account, err := buildAzureServicePrincipalResource(d)
	if err != nil {
		return err
	}

	if account.ID == constEmptyString {
		log.Println("ID is nil")
	} else {
		account.ID = d.Id()
	}

	apiClient := m.(*client.Client)

	updatedAccount, err := apiClient.Accounts.Update(*account)

	if err != nil {
		return createResourceOperationError(errorUpdatingAzureServicePrincipal, d.Id(), err)
	}

	if updatedAccount.ID == constEmptyString {
		log.Println("ID is nil")
	} else {
		d.SetId(updatedAccount.ID)
	}

	return nil
}
