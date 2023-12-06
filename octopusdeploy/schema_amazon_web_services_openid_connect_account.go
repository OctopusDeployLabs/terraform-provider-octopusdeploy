package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandAmazonWebServicesOpenIDConnectAccount(d *schema.ResourceData) *accounts.AwsOIDCAccount {
	name := d.Get("name").(string)
	roleArn := d.Get("role_arn").(string)

	account, _ := accounts.NewAwsOIDCAccount(name, roleArn)
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

	if v, ok := d.GetOk("execution_subject_keys"); ok {
		account.DeploymentSubjectKeys = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("health_subject_keys"); ok {
		account.HealthCheckSubjectKeys = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("account_test_subject_keys"); ok {
		account.AccountTestSubjectKeys = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("session_duration"); ok {
		account.SessionDuration = v.(string)
	}

	return account
}

func getAmazonWebServicesOpenIDConnectAccountSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Description: "A user-friendly description of this AWS OIDC account.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"environments": getEnvironmentsSchema(),
		"name": {
			Description:      "The name of this AWS OIDC account.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, 200)),
		},
		"space_id":                          getSpaceIDSchema(),
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_tags":                       getTenantTagsSchema(),
		"execution_subject_keys":            getSubjectKeysSchema(),
		"health_subject_keys":               getSubjectKeysSchema(),
		"account_test_subject_keys":         getSubjectKeysSchema(),
		"role_arn": {
			Description: "The Amazon Resource Name (ARN) of the role that the caller is assuming.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"session_duration": {
			Description: "The duration, in seconds, of the role session.",
			Required:    false,
			Optional:    true,
			Type:        schema.TypeInt,
		},
	}
}

func setAmazonWebServicesOpenIDConnectAccount(ctx context.Context, d *schema.ResourceData, account *accounts.AwsOIDCAccount) error {
	d.Set("description", account.GetDescription())
	d.Set("name", account.GetName())
	d.Set("space_id", account.GetSpaceID())
	d.Set("tenanted_deployment_participation", account.GetTenantedDeploymentMode())
	d.Set("role_arn", account.RoleArn)

	if err := d.Set("environments", account.GetEnvironmentIDs()); err != nil {
		return fmt.Errorf("error setting environments: %s", err)
	}

	if err := d.Set("tenants", account.GetTenantIDs()); err != nil {
		return fmt.Errorf("error setting tenants: %s", err)
	}

	if err := d.Set("tenant_tags", account.GetTenantTags()); err != nil {
		return fmt.Errorf("error setting tenant_tags: %s", err)
	}

	if err := d.Set("execution_subject_keys", account.DeploymentSubjectKeys); err != nil {
		return fmt.Errorf("error setting execution_subject_keys: %s", err)
	}

	if err := d.Set("health_subject_keys", account.HealthCheckSubjectKeys); err != nil {
		return fmt.Errorf("error setting health_subject_keys: %s", err)
	}

	if err := d.Set("account_test_subject_keys", account.AccountTestSubjectKeys); err != nil {
		return fmt.Errorf("error setting account_test_subject_keys: %s", err)
	}

	return nil
}
