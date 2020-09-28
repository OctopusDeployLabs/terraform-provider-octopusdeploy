package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceUsernamePassword() *schema.Resource {
	schemaMap := getCommonAccountsSchema()

	schemaMap[constUsername] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap[constPassword] = &schema.Schema{
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
	id := d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingUsernamePasswordAccount, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constAccount, m)

	d.Set(constName, resource.Name)
	d.Set(constDescription, resource.Description)
	d.Set(constEnvironments, resource.EnvironmentIDs)
	d.Set(constPassword, resource.Password)

	return nil
}

func buildUsernamePasswordResource(d *schema.ResourceData) (*model.Account, error) {
	name := d.Get(constName).(string)

	account, err := model.NewUsernamePasswordAccount(name)
	if err != nil {
		return nil, err
	}

	password := d.Get(constPassword).(string)
	if isEmpty(password) {
		log.Println("Key is nil. Must add in a password")
	}

	privateKey := model.NewSensitiveValue(password)
	account.Password = &privateKey

	if v, ok := d.GetOk(constTenantedDeploymentParticipation); ok {
		account.TenantedDeploymentParticipation, _ = enum.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk(constTenantTags); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return account, nil
}

func resourceUsernamePasswordCreate(d *schema.ResourceData, m interface{}) error {
	account, err := buildUsernamePasswordResource(d)
	if err != nil {
		return err
	}

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.Add(account)
	if err != nil {
		return createResourceOperationError(errorCreatingUsernamePasswordAccount, account.Name, err)
	}

	if isEmpty(resource.ID) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.ID)
	}

	return nil
}

func resourceUsernamePasswordUpdate(d *schema.ResourceData, m interface{}) error {
	account, err := buildUsernamePasswordResource(d)
	if err != nil {
		return err
	}
	account.ID = d.Id() // set ID so Octopus API knows which account to update

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.Update(*account)
	if err != nil {
		return createResourceOperationError(errorUpdatingUsernamePasswordAccount, d.Id(), err)
	}

	d.SetId(resource.ID)

	return nil
}
