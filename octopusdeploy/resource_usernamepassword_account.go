package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceUsernamePassword() *schema.Resource {
	schemaMap := getCommonAccountsSchema()

	schemaMap["username"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap["password"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	return &schema.Resource{
		Create: resourceUsernamePasswordCreate,
		Read:   resourceUsernamePasswordRead,
		Update: resourceUsernamePasswordUpdate,
		Delete: resourceAccountDeleteCommon,
		Schema: schemaMap,
	}
}

func resourceUsernamePasswordRead(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceUsernamePasswordRead", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceUsernamePasswordRead", "m")
	}

	apiClient := m.(*client.Client)

	accountID := d.Id()
	account, err := apiClient.Accounts.Get(accountID)

	if err == client.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading username password account %s: %s", accountID, err.Error())
	}

	d.Set("name", account.Name)
	d.Set("description", account.Description)
	d.Set("environments", account.EnvironmentIDs)
	d.Set("password", account.Password)

	return nil
}

func buildUsernamePasswordResource(d *schema.ResourceData) (*model.Account, error) {
	accountStruct := model.Account{}
	if accountStruct.Username == "" {
		log.Println("Username struct is nil")
	}

	if d == nil {
		return nil, createInvalidParameterError("buildUsernamePasswordResource", "d")
	}

	name := d.Get("name").(string)

	account, err := model.NewUsernamePasswordAccount(name)

	if err != nil {
		return nil, err
	}

	password := d.Get("password").(string)
	if password == "" {
		log.Println("Key is nil. Must add in a password")
	}

	privateKey := model.NewSensitiveValue(password)
	account.Password = &privateKey

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentParticipation, _ = enum.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return account, nil
}

func resourceUsernamePasswordCreate(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceAzureServicePrincipalRead", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceAzureServicePrincipalRead", "m")
	}

	apiClient := m.(*client.Client)

	newAccount, err := buildUsernamePasswordResource(d)
	if err != nil {
		return err
	}

	account, err := apiClient.Accounts.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error creating username password account %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceUsernamePasswordUpdate(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceAzureServicePrincipalRead", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceAzureServicePrincipalRead", "m")
	}

	account, err := buildUsernamePasswordResource(d)
	if err != nil {
		return err
	}

	account.ID = d.Id()

	apiClient := m.(*client.Client)

	updatedAccount, err := apiClient.Accounts.Update(*account)

	if err != nil {
		return fmt.Errorf("error updating username password account id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}
