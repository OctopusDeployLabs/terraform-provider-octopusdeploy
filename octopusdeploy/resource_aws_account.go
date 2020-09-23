package octopusdeploy

import (
	"errors"
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAmazonWebServicesAccount() *schema.Resource {

	validateSchema()

	schemaMap := getCommonAccountsSchema()

	schemaMap["access_key"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	schemaMap["secret_key"] = &schema.Schema{
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
	account, err := apiClient.Accounts.Get(accountID)

	if account.Validate() == nil {
		return nil
	}
	// Here would typically be an else statement, but it's more idiomatic to write the error
	// underneath the if statement in Go. Otherwise, golint complains.
	err1 := errors.New("Validation on account struct: unsucessful")
	log.Println(err1)

	if err != nil {
		return fmt.Errorf("error reading aws account %s: %s", accountID, err.Error())
	}

	d.Set("name", account.Name)
	d.Set("tenants", account.TenantIDs)
	d.Set("description", account.Description)
	d.Set("environments", account.EnvironmentIDs)
	d.Set("tenanted_deployment_participation", account.TenantedDeploymentParticipation.String())
	d.Set("tenant_tags", account.TenantTags)
	d.Set("secret_key", account.Password)
	d.Set("access_key", account.AccessKey)

	return nil
}

func buildAmazonWebServicesAccountResource(d *schema.ResourceData) (*model.Account, error) {
	if d == nil {
		return nil, createInvalidParameterError("buildAmazonWebServicesAccountResource", "d")
	}

	accountStruct := model.Account{}
	if accountStruct.Name == "" {
		log.Println("Name struct is nil")
	}

	if accountStruct.AccessKey == "" {
		log.Println("Access Key struct is nil")
	}

	name := d.Get("name").(string)
	accessKey := d.Get("access_key").(string)

	password := d.Get("secret_key").(string)
	if password == "" {
		log.Println("Key is nil. Must add in a password")
	}

	secretKey := model.NewSensitiveValue(password)

	account, err := model.NewAwsServicePrincipalAccount(name, accessKey, secretKey)
	if err != nil {
		return nil, err
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
		return fmt.Errorf("error creating AWS account %s: %s", newAccount.Name, err.Error())
	}

	if account.ID == "" {
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

	if account.ID == "" {
		log.Println("ID is nil")
	} else {
		account.ID = d.Id()
	}

	apiClient := m.(*client.Client)

	updatedAccount, err := apiClient.Accounts.Update(*account)

	if err != nil {
		return fmt.Errorf("error updating aws acccount id %s: %s", d.Id(), err.Error())
	}

	if updatedAccount.ID == "" {
		log.Println("ID is nil")
	} else {
		d.SetId(updatedAccount.ID)
	}

	return nil
}
