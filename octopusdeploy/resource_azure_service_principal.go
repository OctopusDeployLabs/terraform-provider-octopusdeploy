package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAzureServicePrincipal() *schema.Resource {
	validateSchema()
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
		Type:             schema.TypeString,
		Required:         true,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
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
		Create:        resourceAzureServicePrincipalCreate,
		Read:          resourceAzureServicePrincipalRead,
		Update:        resourceAzureServicePrincipalUpdate,
		DeleteContext: resourceAccountDeleteCommon,
		Schema:        schemaMap,
	}
}

func buildAzureServicePrincipalResource(d *schema.ResourceData) (*octopusdeploy.AzureServicePrincipalAccount, error) {
	name := d.Get(constName).(string)

	password := d.Get(constKey).(string)
	if isEmpty(password) {
		log.Println("Key is nil. Must add in a password")
	}

	secretKey := octopusdeploy.NewSensitiveValue(password)

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

	account, err := octopusdeploy.NewAzureServicePrincipalAccount(name, subscriptionID, tenantID, applicationID, secretKey)
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
		account.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk(constTenantTags); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constTenants); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constResourceManagementEndpointBaseURI); ok {
		account.ResourceManagerEndpoint = v.(string)
	}

	if v, ok := d.GetOk(constActiveDirectoryEndpointBaseURI); ok {
		account.AuthenticationEndpoint = v.(string)
	}

	err = account.Validate()
	if err != nil {
		return nil, err
	}

	return account, nil
}

func resourceAzureServicePrincipalCreate(d *schema.ResourceData, m interface{}) error {
	account, err := buildAzureServicePrincipalResource(d)
	if err != nil {
		log.Println(err)
		return err
	}

	client := m.(*octopusdeploy.Client)
	resource, err := client.Accounts.Add(account)
	if err != nil {
		return createResourceOperationError(errorCreatingAzureServicePrincipal, account.GetName(), err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func resourceAzureServicePrincipalRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	resource, err := client.Accounts.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingAzureServicePrincipal, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constAccount, m)

	accountResource := resource.(*octopusdeploy.AccountResource)

	d.Set(constName, accountResource.Name)
	d.Set(constDescription, accountResource.Description)
	d.Set(constEnvironments, accountResource.EnvironmentIDs)
	d.Set(constTenantedDeploymentParticipation, accountResource.TenantedDeploymentMode)
	d.Set(constTenantTags, accountResource.TenantTags)
	d.Set(constClientID, accountResource.ApplicationID.String())
	d.Set(constTenantID, accountResource.TenantID.String())
	d.Set(constSubscriptionNumber, accountResource.SubscriptionID.String())

	// TODO: determine what to do here...
	// d.Set(constKey, accountResource.ApplicationPassword)

	d.Set(constAzureEnvironment, accountResource.AzureEnvironment)
	d.Set(constResourceManagementEndpointBaseURI, accountResource.ResourceManagerEndpoint)
	d.Set(constActiveDirectoryEndpointBaseURI, accountResource.AuthenticationEndpoint)

	return nil
}

func resourceAzureServicePrincipalUpdate(d *schema.ResourceData, m interface{}) error {
	account, err := buildAzureServicePrincipalResource(d)
	if err != nil {
		return err
	}
	account.ID = d.Id() // set ID so Octopus API knows which account to update

	client := m.(*octopusdeploy.Client)
	resource, err := client.Accounts.Update(account)
	if err != nil {
		return createResourceOperationError(errorUpdatingAzureServicePrincipal, d.Id(), err)
	}

	d.SetId(resource.GetID())

	return nil
}
