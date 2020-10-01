package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/enum"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSSHKey() *schema.Resource {
	schemaMap := getCommonAccountsSchema()

	schemaMap[constUsername] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	schemaMap["passphrase"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	return &schema.Resource{
		CreateContext: resourceSSHKeyCreate,
		ReadContext:   resourceSSHKeyRead,
		UpdateContext: resourceSSHKeyUpdate,
		DeleteContext: resourceAccountDeleteCommon,
		Schema:        schemaMap,
	}
}

func resourceSSHKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	diagValidate()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.GetByID(id)
	if err != nil {
		// return createResourceOperationError(errorReadingSSHKeyPair, id, err)
		diag.FromErr(err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constAccount, m)

	d.Set(constName, resource.Name)
	d.Set("passphrase", resource.Password)
	d.Set(constTenants, resource.TenantIDs)

	return nil
}

func buildSSHKeyResource(d *schema.ResourceData) (*model.Account, error) {
	name := d.Get(constName).(string)
	username := d.Get(constUsername).(string)

	password := d.Get(constPassword).(string)
	if isEmpty(password) {
		log.Println("Key is nil. Must add in a password")
	}

	secretKey := model.NewSensitiveValue(password)

	account, err := model.NewSSHKeyAccount(name, username, secretKey)
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

func resourceSSHKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, err := buildSSHKeyResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	diagValidate()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.Add(account)
	if err != nil {
		// return createResourceOperationError(errorCreatingSSHKeyPair, account.Name, err)
		diag.FromErr(err)
	}

	if isEmpty(resource.ID) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.ID)
	}

	return nil
}

func resourceSSHKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, err := buildSSHKeyResource(d)
	if err != nil {
		return diag.FromErr(err)
	}
	account.ID = d.Id() // set ID so Octopus API knows which account to update

	diagValidate()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Accounts.Update(*account)
	if err != nil {
		// return createResourceOperationError(errorUpdatingSSHKeyPair, d.Id(), err)
		diag.FromErr(err)
	}

	d.SetId(resource.ID)

	return nil
}
