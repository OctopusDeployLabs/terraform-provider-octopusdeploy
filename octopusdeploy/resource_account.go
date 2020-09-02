package model

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountCreate,
		Read:   resourceAccountRead,
		Update: resourceAccountUpdate,
		Delete: resourceAccountDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environments": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func resourceAccountRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	accountID := d.Id()
	account, err := apiClient.Accounts.Get(accountID)

	if err == client.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading account %s: %s", accountID, err.Error())
	}

	d.Set("name", account.Name)
	d.Set("environments", account.EnvironmentIDs)
	d.Set("account_type", account.AccountType.String())
	d.Set("client_id", account.ClientID)
	d.Set("tenant_id", account.TenantID)
	d.Set("subscription_id", account.SubscriptionNumber)
	d.Set("client_secret", account.Password)
	d.Set("tenant_tags", account.TenantTags)
	d.Set("tenanted_deployment_participation", account.TenantedDeploymentParticipation.String())
	d.Set("token", account.Token)

	return nil
}

func buildAccountResource(d *schema.ResourceData) *model.Account {
	accountName := d.Get("name").(string)

	var environments []string
	var accountType string
	var clientID uuid.UUID
	var tenantID uuid.UUID
	var subscriptionID uuid.UUID
	var clientSecret string
	var tenantTags []string
	var tenantedDeploymentParticipation string
	var token string

	environmentsInterface, ok := d.GetOk("environments")
	if ok {
		environments = getSliceFromTerraformTypeList(environmentsInterface)
	}

	accountTypeInterface, ok := d.GetOk("account_type")
	if ok {
		accountType = accountTypeInterface.(string)
	}

	clientIDInterface, ok := d.GetOk("client_id")
	if ok {
		clientID = clientIDInterface.(uuid.UUID)
	}

	tenantIDInterface, ok := d.GetOk("tenant_id")
	if ok {
		tenantID = tenantIDInterface.(uuid.UUID)
	}

	subscriptionIDInterface, ok := d.GetOk("subscription_id")
	if ok {
		subscriptionID = subscriptionIDInterface.(uuid.UUID)
	}

	clientSecretInterface, ok := d.GetOk("client_secret")
	if ok {
		clientSecret = clientSecretInterface.(string)
	}

	tenantedDeploymentParticipationInterface, ok := d.GetOk("tenanted_deployment_participation")
	if ok {
		tenantedDeploymentParticipation = tenantedDeploymentParticipationInterface.(string)
	}

	tenantTagsInterface, ok := d.GetOk("tenant_tags")
	if ok {
		tenantTags = getSliceFromTerraformTypeList(tenantTagsInterface)
	}

	if tenantTags == nil {
		tenantTags = []string{}
	}

	tokenInterface, ok := d.GetOk("token")
	if ok {
		token = tokenInterface.(string)
	}

	accountTypeEnum, _ := enum.ParseAccountType(accountType)

	account, err := model.NewAccount(accountName, accountTypeEnum)
	if err != nil {
		return nil
	}

	account.EnvironmentIDs = environments
	account.ClientID = &clientID
	account.TenantID = &tenantID
	account.Password = &model.SensitiveValue{
		NewValue: &clientSecret,
	}
	account.SubscriptionNumber = &subscriptionID
	account.TenantTags = tenantTags
	account.TenantedDeploymentParticipation, _ = enum.ParseTenantedDeploymentMode(tenantedDeploymentParticipation)
	account.Token = &model.SensitiveValue{
		NewValue: &token,
	}

	return account
}

func resourceAccountCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	newAccount := buildAccountResource(d)
	account, err := apiClient.Accounts.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error creating account %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceAccountUpdate(d *schema.ResourceData, m interface{}) error {
	account := buildAccountResource(d)
	account.ID = d.Id() // set project struct ID so octopus knows which project to update

	apiClient := m.(*client.Client)

	updatedAccount, err := apiClient.Accounts.Update(account)

	if err != nil {
		return fmt.Errorf("error updating account id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceAccountDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	accountID := d.Id()

	err := apiClient.Accounts.Delete(accountID)

	if err != nil {
		return fmt.Errorf("error deleting account id %s: %s", accountID, err.Error())
	}

	d.SetId("")
	return nil
}
