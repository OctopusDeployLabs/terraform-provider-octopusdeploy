package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandAmazonWebServicesAccount(d *schema.ResourceData) *octopusdeploy.AmazonWebServicesAccount {
	name := d.Get("name").(string)
	accessKey := d.Get("access_key").(string)
	password := d.Get("secret_key").(string)
	secretKey := octopusdeploy.NewSensitiveValue(password)

	account, _ := octopusdeploy.NewAmazonWebServicesAccount(name, accessKey, secretKey)
	account.ID = d.Id()

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenants"); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	return account
}

func getAmazonWebServicesAccountSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_key": {
			Required:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"environments": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"name": &schema.Schema{
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"secret_key": {
			Required:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"space_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"tenant_tags": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}

func setAmazonWebServicesAccount(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.AmazonWebServicesAccount) {
	d.Set("access_key", account.AccessKey)
	d.Set("description", account.GetDescription())
	d.Set("environments", account.GetEnvironmentIDs())
	d.Set("id", account.GetID())
	d.Set("name", account.GetName())
	d.Set("space_id", account.GetSpaceID())
	d.Set("tenanted_deployment_participation", account.GetTenantedDeploymentMode())
	d.Set("tenants", account.GetTenantIDs())
	d.Set("tenant_tags", account.GetTenantTags())

	d.SetId(account.GetID())
}
