package octopusdeploy

import (
	"context"
	"fmt"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandSSHKeyAccount(d *schema.ResourceData) *accounts.SSHKeyAccount {
	name := d.Get("name").(string)
	username := d.Get("username").(string)
	privateKeyFile := core.NewSensitiveValue(d.Get("private_key_file").(string))

	account, _ := accounts.NewSSHKeyAccount(name, username, privateKeyFile)
	account.ID = d.Id()
	account.Description = d.Get("description").(string)

	if v, ok := d.GetOk("private_key_passphrase"); ok {
		account.SetPrivateKeyPassphrase(core.NewSensitiveValue(v.(string)))
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
			Optional:         true,
			Sensitive:        true,
			Type:             schema.TypeString,
			ValidateDiagFunc: warnIfSshPassphraseLooksLikeFile(),
		},
		"space_id":                          getSpaceIDSchema(),
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_tags":                       getTenantTagsSchema(),
		"username":                          getUsernameSchema(true),
	}
}

func warnIfSshPassphraseLooksLikeFile() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		value := v.(string)

		// Certificates often start with "---BEGIN".
		// If we can detect that, the chances are the user was setting this value to work around the bug

		if strings.HasPrefix(value, "LS0tLS1CRUdJTi") {

			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Certificate value used in \"private_key_passphrase\"",
				Detail: `The "private_key_passphrase" appears to be a certificate file.
This may be due to a previous bug with the provider which has been fixed (https://github.com/OctopusDeployLabs/terraform-provider-octopusdeploy/issues/343).
It is advised that you instead set this value on "private_key_file" and leave the passphrase "private_key_passphrase" field blank if none apply.

This warning can be ignored if the passphrase value is expected`,
				AttributePath: path,
			})
		}

		return diags
	}
}

func setSSHKeyAccount(ctx context.Context, d *schema.ResourceData, account *accounts.SSHKeyAccount) error {
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
