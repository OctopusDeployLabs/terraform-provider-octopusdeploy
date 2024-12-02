package octopusdeploy

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandGenericOpenIDConnectAccount(d *schema.ResourceData) *accounts.GenericOIDCAccount {
	name := d.Get("name").(string)

	account, _ := accounts.NewGenericOIDCAccount(name)
	account.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		account.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("environments"); ok {
		account.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("name"); ok {
		account.SetName(v.(string))
	}

	if v, ok := d.GetOk("space_id"); ok {
		account.SetSpaceID(v.(string))
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

	if v, ok := d.GetOk("audience"); ok {
		account.Audience = v.(string)
	}

	if v, ok := d.GetOk("execution_subject_keys"); ok {
		account.DeploymentSubjectKeys = getSliceFromTerraformTypeList(v)
	}

	return account
}

func getGenericOpenIdConnectAccountSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description":                       getDescriptionSchema("Generic OpenID Connect account"),
		"environments":                      getEnvironmentsSchema(),
		"id":                                getIDSchema(),
		"name":                              getNameSchema(true),
		"space_id":                          getSpaceIDSchema(),
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_tags":                       getTenantTagsSchema(),
		"execution_subject_keys":            getSubjectKeysSchema(SchemaSubjectKeysDescriptionExecution),
		"audience":                          getOidcAudienceSchema(),
	}
}

func setGenericOpenIDConnectAccount(ctx context.Context, d *schema.ResourceData, account *accounts.GenericOIDCAccount) error {
	d.Set("description", account.GetDescription())
	d.Set("id", account.GetID())
	d.Set("name", account.GetName())
	d.Set("space_id", account.GetSpaceID())
	d.Set("tenanted_deployment_participation", account.GetTenantedDeploymentMode())
	d.Set("audience", account.Audience)

	if err := d.Set("environments", account.GetEnvironmentIDs()); err != nil {
		return fmt.Errorf("error setting environments: %s", err)
	}

	if err := d.Set("tenants", account.GetTenantIDs()); err != nil {
		return fmt.Errorf("error setting tenants: %s", err)
	}

	if err := d.Set("tenant_tags", account.TenantTags); err != nil {
		return fmt.Errorf("error setting tenant_tags: %s", err)
	}

	if err := d.Set("execution_subject_keys", account.DeploymentSubjectKeys); err != nil {
		return fmt.Errorf("error setting execution_subject_keys: %s", err)
	}

	return nil
}
