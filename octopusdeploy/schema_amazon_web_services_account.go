package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandAmazonWebServicesAccount(d *schema.ResourceData) *octopusdeploy.AmazonWebServicesAccount {
	name := d.Get("name").(string)
	accessKey := d.Get("access_key").(string)
	secretKey := octopusdeploy.NewSensitiveValue(d.Get("secret_key").(string))

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
		"access_key":                        getAccessKeySchema(true),
		"description":                       getDescriptionSchema(),
		"environments":                      getEnvironmentsSchema(),
		"id":                                getIDSchema(),
		"name":                              getNameSchema(true),
		"secret_key":                        getSecretKeySchema(true),
		"space_id":                          getSpaceIDSchema(),
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_tags":                       getTenantTagsSchema(),
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
