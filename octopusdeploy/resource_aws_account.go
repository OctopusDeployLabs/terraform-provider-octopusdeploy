package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAmazonWebServicesAccount() *schema.Resource {

	validateSchema()

	schemaMap := getCommonAccountsSchema()

	schemaMap[constAccessKey] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	schemaMap[constSecretKey] = &schema.Schema{
		Type:      schema.TypeString,
		Required:  true,
		Sensitive: true,
	}
	return &schema.Resource{
		Create: resourceAmazonWebServicesAccountCreate,
		Read:   resourceAmazonWebServicesAccountRead,
		Update: resourceAmazonWebServicesAccountUpdate,
		Delete: resourceAccountDeleteCommon,
		Schema: schemaMap,
	}
}

func resourceAmazonWebServicesAccountRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingAWSAccount, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constAccount, m)

	d.Set(constName, resource.Name)
	d.Set(constTenants, resource.TenantIDs)
	d.Set(constDescription, resource.Description)
	d.Set(constEnvironments, resource.EnvironmentIDs)
	d.Set(constTenantedDeploymentParticipation, resource.TenantedDeploymentParticipation.String())
	d.Set(constTenantTags, resource.TenantTags)
	d.Set(constSecretKey, resource.Password)
	d.Set(constAccessKey, resource.AccessKey)

	return nil
}

func buildAmazonWebServicesAccountResource(d *schema.ResourceData) (*model.Account, error) {
	if d == nil {
		return nil, createInvalidParameterError("buildAmazonWebServicesAccountResource", "d")
	}

	name := d.Get(constName).(string)
	accessKey := d.Get(constAccessKey).(string)
	password := d.Get(constSecretKey).(string)
	secretKey := model.NewSensitiveValue(password)

	account, err := model.NewAwsServicePrincipalAccount(name, accessKey, secretKey)
	if err != nil {
		return nil, err
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

	return account, nil
}

func resourceAmazonWebServicesAccountCreate(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceAmazonWebServicesAccountCreate", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceAmazonWebServicesAccountCreate", "m")
	}

	apiClient := m.(*client.Client)

	newAccount, err := buildAmazonWebServicesAccountResource(d)
	if err != nil {
		log.Println(err)
		return err
	}

	account, err := apiClient.Accounts.Add(newAccount)
	if err != nil {
		return createResourceOperationError(errorCreatingAWSAccount, newAccount.Name, err)
	}

	if account.ID == constEmptyString {
		log.Println("ID is nil")
	} else {
		d.SetId(account.ID)
	}

	return nil
}

func resourceAmazonWebServicesAccountUpdate(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceAmazonWebServicesAccountUpdate", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceAmazonWebServicesAccountUpdate", "m")
	}

	account, err := buildAmazonWebServicesAccountResource(d)
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
		return createResourceOperationError(errorUpdatingAWSAccount, d.Id(), err)
	}

	if updatedAccount.ID == constEmptyString {
		log.Println("ID is nil")
	} else {
		d.SetId(updatedAccount.ID)
	}

	return nil
}
