package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandAmazonWebServicesAccount(d *schema.ResourceData) *accounts.AmazonWebServicesAccount {
	name := d.Get("name").(string)
	accessKey := d.Get("access_key").(string)
	secretKey := core.NewSensitiveValue(d.Get("secret_key").(string))

	account, _ := accounts.NewAmazonWebServicesAccount(name, accessKey, secretKey)
	account.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		account.Description = v.(string)
	}

	if v, ok := d.GetOk("environments"); ok {
		account.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("space_id"); ok {
		account.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.TenantedDeploymentMode = core.TenantedDeploymentMode(v.(string))
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
			Description: "The access key associated with this AWS account.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"description": {
			Description: "A user-friendly description of this AWS account.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"environments": getEnvironmentsSchema(),
		"name": {
			Description:      "The name of this AWS account.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 200)),
		},
		"secret_key":                        getSecretKeySchema(true),
		"space_id":                          getSpaceIDSchema(),
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_tags":                       getTenantTagsSchema(),
	}
}

func setAmazonWebServicesAccount(ctx context.Context, d *schema.ResourceData, account *accounts.AmazonWebServicesAccount) error {
	d.Set("access_key", account.AccessKey)
	d.Set("description", account.GetDescription())
	d.Set("name", account.GetName())
	d.Set("space_id", account.GetSpaceID())
	d.Set("tenanted_deployment_participation", account.GetTenantedDeploymentMode())

	if err := d.Set("environments", account.GetEnvironmentIDs()); err != nil {
		return fmt.Errorf("error setting environments: %s", err)
	}

	if err := d.Set("tenants", account.GetTenantIDs()); err != nil {
		return fmt.Errorf("error setting tenants: %s", err)
	}

	if err := d.Set("tenant_tags", account.GetTenantTags()); err != nil {
		return fmt.Errorf("error setting tenant_tags: %s", err)
	}

	return nil
}
