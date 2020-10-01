package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		CreateContext: resourceUsernamePasswordCreate,
		ReadContext:   resourceUsernamePasswordRead,
		UpdateContext: resourceUsernamePasswordUpdate,
		DeleteContext: resourceAccountDeleteCommon,
		Schema:        schemaMap,
	}
}

func resourceUsernamePasswordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	diagValidate()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.GetByID(id)
	if err != nil {
		// return createResourceOperationError(errorReadingUsernamePasswordAccount, id, err)
		return diag.FromErr(err)
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

func buildUsernamePasswordResource(d *schema.ResourceData) (*model.Account, diag.Diagnostics) {
	name := d.Get(constName).(string)

	account, err := model.NewUsernamePasswordAccount(name)
	if err != nil {
		return nil, diag.FromErr(err)
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

func resourceUsernamePasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, _ := buildUsernamePasswordResource(d)

	diagValidate()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.Add(account)
	if err != nil {
		// return createResourceOperationError(errorCreatingUsernamePasswordAccount, account.Name, err)
		return diag.FromErr(err)
	}

	if isEmpty(resource.ID) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.ID)
	}

	return nil
}

func resourceUsernamePasswordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, _ := buildUsernamePasswordResource(d)

	diagValidate()

	account.ID = d.Id() // set ID so Octopus API knows which account to update

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.Update(*account)
	if err != nil {
		// return createResourceOperationError(errorUpdatingUsernamePasswordAccount, d.Id(), err)
		return diag.FromErr(err)
	}

	d.SetId(resource.ID)

	return nil
}
