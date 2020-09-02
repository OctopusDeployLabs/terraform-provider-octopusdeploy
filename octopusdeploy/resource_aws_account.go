package model

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
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
	apiClient := m.(*client.Client)

	accountID := d.Id()
	account, err := apiClient.Accounts.Get(accountID)

	if err == client.ErrItemNotFound {
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

func buildAmazonWebServicesAccountResource(d *schema.ResourceData) *model.Account {
	account, err := model.NewAccount(d.Get("name").(string), enum.AmazonWebServicesAccount)
	if err != nil {
		return nil
	}

	account.Name = d.Get("name").(string)
	account.AccessKey = d.Get("access_key").(string)
	password := d.Get("secret_key").(string)
	account.Password = &model.SensitiveValue{NewValue: &password}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentParticipation, _ = enum.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return account
}

func resourceAmazonWebServicesAccountCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	newAccount := buildAmazonWebServicesAccountResource(d)
	account, err := apiClient.Accounts.Add(newAccount)

	if err != nil {
		return fmt.Errorf("error creating AWS account %s: %s", newAccount.Name, err.Error())
	}

	d.SetId(account.ID)

	return nil
}

func resourceAmazonWebServicesAccountUpdate(d *schema.ResourceData, m interface{}) error {
	account := buildAmazonWebServicesAccountResource(d)
	account.ID = d.Id()

	apiClient := m.(*client.Client)

	updatedAccount, err := apiClient.Accounts.Update(account)

	if err != nil {
		return fmt.Errorf("error updating aws acccount id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedAccount.ID)
	return nil
}

func resourceAmazonWebServicesAccountDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	accountID := d.Id()

	err := apiClient.Accounts.Delete(accountID)

	if err != nil {
		return fmt.Errorf("error deleting AWS account id %s: %s", accountID, err.Error())
	}

	d.SetId("")
	return nil
}
