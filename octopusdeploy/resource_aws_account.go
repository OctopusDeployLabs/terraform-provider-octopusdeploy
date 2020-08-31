package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAmazonWebServicesAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAmazonWebServicesAccountCreate,
		Read:   resourceAmazonWebServicesAccountRead,
		Update: resourceAmazonWebServicesAccountUpdate,
		Delete: resourceAmazonWebServicesAccountDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "AWS",
			},
			"environments": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"tenant_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenanted_deployment_participation": getTenantedDeploymentSchema(),
			"secret_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_key": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAmazonWebServicesAccountRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountID := d.Id()
	account, err := client.Account.Get(accountID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading aws account %s: %s", accountID, err.Error())
	}

	d.Set("name", account.Name)
	d.Set("description", account.Description)
	d.Set("environments", account.EnvironmentIDs)
	d.Set("tenanted_deployment_participation", account.TenantedDeploymentParticipation.String())
	d.Set("tenant_tags", account.TenantTags)
	d.Set("secret_key", account.SecretKey)
	d.Set("access_key", account.AccessKey)

	return nil
}

func buildAmazonWebServicesAccountResource(d *schema.ResourceData) *octopusdeploy.Account {
	var account = octopusdeploy.NewAccount(d.Get("name").(string), octopusdeploy.AmazonWebServicesAccount)

	account.Name = d.Get("name").(string)
	account.AccessKey = d.Get("access_key").(string)
	account.Password = octopusdeploy.SensitiveValue{NewValue: d.Get("secret_key").(string)}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentParticipation, _ = octopusdeploy.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return account
}

func resourceAmazonWebServicesAccountCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newAccount := buildAmazonWebServicesAccountResource(d)
	account, err := client.Account.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error creating AWS account %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceAmazonWebServicesAccountUpdate(d *schema.ResourceData, m interface{}) error {
	account := buildAmazonWebServicesAccountResource(d)
	account.ID = d.Id()

	client := m.(*octopusdeploy.Client)

	updatedAccount, err := client.Account.Update(account)

	if err != nil {
		return fmt.Errorf("error updating aws acccount id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceAmazonWebServicesAccountDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	accountID := d.Id()

	err := client.Account.Delete(accountID)

	if err != nil {
		return fmt.Errorf("error deleting AWS account id %s: %s", accountID, err.Error())
	}

	d.SetId("")
	return nil
}
