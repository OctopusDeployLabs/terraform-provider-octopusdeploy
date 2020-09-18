package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/enum"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSSHKey() *schema.Resource {
	schemaMap := getCommonAccountsSchema()

	schemaMap["username"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap["passphrase"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	return &schema.Resource{
		Create: resourceSSHKeyCreate,
		Read:   resourceSSHKeyRead,
		Update: resourceSSHKeyUpdate,
		Delete: resourceAccountDeleteCommon,
		Schema: schemaMap,
	}
}

func resourceSSHKeyRead(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceSSHKeyRead", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceSSHKeyRead", "m")
	}

	apiClient := m.(*client.Client)

	accountID := d.Id()
	account, err := apiClient.Accounts.Get(accountID)

	if err == client.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading SSH Key Pair %s: %s", accountID, err.Error())
	}

	d.Set("name", account.Name)
	d.Set("passphrase", account.Password)
	d.Set("tenants", account.TenantIDs)

	return nil
}

func buildSSHKeyResource(d *schema.ResourceData) (*model.Account, error) {
	accountStruct := model.Account{}
	if accountStruct.Name == "" {
		log.Println("Name struct is nil")
	}

	if d == nil {
		return nil, createInvalidParameterError("buildSSHKeyResource", "d")
	}

	name := d.Get("name").(string)
	userName := d.Get("username").(string)

	password := d.Get("password").(string)
	if password == "" {
		log.Println("Key is nil. Must add in a password")
	}

	secretKey := model.NewSensitiveValue(password)

	account, err := model.NewSSHKeyAccount(name, userName, secretKey)
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

func resourceSSHKeyCreate(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceAzureServicePrincipalRead", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceAzureServicePrincipalRead", "m")
	}

	if d == nil {
		return createInvalidParameterError("resourceSSHKeyCreate", "d")
	}

	apiClient := m.(*client.Client)

	newAccount, err := buildSSHKeyResource(d)
	if err != nil {
		return err
	}

	account, err := apiClient.Accounts.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error creating SSH Key Pair %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceSSHKeyUpdate(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceAzureServicePrincipalRead", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceAzureServicePrincipalRead", "m")
	}

	account, err := buildSSHKeyResource(d)
	if err != nil {
		return err
	}

	account.ID = d.Id()

	apiClient := m.(*client.Client)

	updatedAccount, err := apiClient.Accounts.Update(*account)

	if err != nil {
		return fmt.Errorf("error updating SSH Key Pair %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}
