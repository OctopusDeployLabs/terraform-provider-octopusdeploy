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
	if d == nil {
		return createInvalidParameterError("resourceAmazonWebServicesAccountRead", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceAmazonWebServicesAccountRead", "m")
	}

	apiClient := m.(*client.Client)
	accountID := d.Id()
	account, err := apiClient.Accounts.GetByID(accountID)

	if err != nil {
		return createResourceOperationError(errorReadingAWSAccount, accountID, err)
	}
	if account == nil {
		d.SetId(constEmptyString)
		return nil
	}

	d.Set(constName, account.Name)
	d.Set(constTenants, account.TenantIDs)
	d.Set(constDescription, account.Description)
	d.Set(constEnvironments, account.EnvironmentIDs)
	d.Set(constTenantedDeploymentParticipation, account.TenantedDeploymentParticipation.String())
	d.Set(constTenantTags, account.TenantTags)
	d.Set(constSecretKey, account.Password)
	d.Set(constAccessKey, account.AccessKey)

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
