package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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
			"account_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subscription_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_secret": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenanted_deployment_participation": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateValueFunc([]string{
					"Untenanted",
					"TenantedOrUntenanted",
					"Tenanted",
				}),
			},
			"token": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAccountRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountID := d.Id()
	account, err := client.Account.Get(accountID)

	if err == octopusdeploy.ErrItemNotFound {
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

func buildAccountResource(d *schema.ResourceData) *octopusdeploy.Account {
	accountName := d.Get("name").(string)

	var environments []string
	var accountType string
	var clientID string
	var tenantID string
	var subscriptionID string
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
		clientID = clientIDInterface.(string)
	}

	tenantIDInterface, ok := d.GetOk("tenant_id")
	if ok {
		tenantID = tenantIDInterface.(string)
	}

	subscriptionIDInterface, ok := d.GetOk("subscription_id")
	if ok {
		subscriptionID = subscriptionIDInterface.(string)
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

	accountTypeEnum, _ := octopusdeploy.ParseAccountType(accountType)
	var account = octopusdeploy.NewAccount(accountName, accountTypeEnum)
	account.EnvironmentIDs = environments
	account.ClientID = clientID
	account.TenantID = tenantID
	account.Password = octopusdeploy.SensitiveValue{
		NewValue: clientSecret,
	}
	account.SubscriptionNumber = subscriptionID
	account.TenantTags = tenantTags
	account.TenantedDeploymentParticipation, _ = octopusdeploy.ParseTenantedDeploymentMode(tenantedDeploymentParticipation)
	account.Token = octopusdeploy.SensitiveValue{
		NewValue: token,
	}

	return account
}

func resourceAccountCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newAccount := buildAccountResource(d)
	account, err := client.Account.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error creating account %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceAccountUpdate(d *schema.ResourceData, m interface{}) error {
	account := buildAccountResource(d)
	account.ID = d.Id() // set project struct ID so octopus knows which project to update

	client := m.(*octopusdeploy.Client)

	updatedAccount, err := client.Account.Update(account)

	if err != nil {
		return fmt.Errorf("error updating account id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceAccountDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountID := d.Id()

	err := client.Account.Delete(accountID)

	if err != nil {
		return fmt.Errorf("error deleting account id %s: %s", accountID, err.Error())
	}

	d.SetId("")
	return nil
}
