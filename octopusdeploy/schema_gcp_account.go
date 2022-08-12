package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandGoogleCloudPlatformAccount(d *schema.ResourceData) *accounts.GoogleCloudPlatformAccount {
	name := d.Get("name").(string)
	jsonKey := core.NewSensitiveValue(d.Get("json_key").(string))

	account, _ := accounts.NewGoogleCloudPlatformAccount(name, jsonKey)
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

func getGoogleCloudPlatformAccountSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Description: "A user-friendly description of this GCP account.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"environments": getEnvironmentsSchema(),
		"json_key": {
			Description: "The JSON key associated with this GCP account.",
			Required:    true,
			Sensitive:   true,
			Type:        schema.TypeString,
		},
		"name": {
			Description:      "The name of this GCP account.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 200)),
		},
		"space_id":                          getSpaceIDSchema(),
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_tags":                       getTenantTagsSchema(),
	}
}

func setGoogleCloudPlatformAccount(ctx context.Context, d *schema.ResourceData, account *accounts.GoogleCloudPlatformAccount) error {
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
