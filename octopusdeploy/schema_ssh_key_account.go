package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandSSHKeyAccount(d *schema.ResourceData) *octopusdeploy.SSHKeyAccount {
	name := d.Get("name").(string)
	username := d.Get("username").(string)
	privateKeyFile := octopusdeploy.NewSensitiveValue(d.Get("private_key_passphrase").(string))

	account, _ := octopusdeploy.NewSSHKeyAccount(name, username, privateKeyFile)
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

func getSSHKeyAccountSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description":  getDescriptionSchema("SSH key account"),
		"environments": getEnvironmentsSchema(),
		"id":           getIDSchema(),
		"name":         getNameSchema(true),
		"private_key_file": {
			Required:         true,
			Sensitive:        true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"private_key_passphrase": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
		"space_id":                          getSpaceIDSchema(),
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_tags":                       getTenantTagsSchema(),
		"username":                          getUsernameSchema(true),
	}
}

func setSSHKeyAccount(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.SSHKeyAccount) error {
	d.Set("description", account.GetDescription())

	if err := d.Set("environments", account.GetEnvironmentIDs()); err != nil {
		return fmt.Errorf("error setting environments: %s", err)
	}

	d.Set("name", account.GetName())
	d.Set("space_id", account.GetSpaceID())
	d.Set("tenanted_deployment_participation", account.GetTenantedDeploymentMode())

	if err := d.Set("tenants", account.GetTenantIDs()); err != nil {
		return fmt.Errorf("error setting tenants: %s", err)
	}

	if err := d.Set("tenant_tags", account.TenantTags); err != nil {
		return fmt.Errorf("error setting tenant_tags: %s", err)
	}

	d.Set("username", account.Username)

	return nil
}
