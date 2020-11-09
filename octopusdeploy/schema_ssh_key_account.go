package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandSSHKeyAccount(d *schema.ResourceData) *octopusdeploy.SSHKeyAccount {
	name := d.Get("name").(string)
	username := d.Get("username").(string)
	privateKeyFile := octopusdeploy.NewSensitiveValue(d.Get("passphrase").(string))

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

func flattenSSHKeyAccount(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.SSHKeyAccount) {
	flattenAccount(ctx, d, account)

	d.Set("account_type", "SshKeyPair")
	d.Set("private_key_file", account.PrivateKeyFile)
	d.Set("private_key_passphrase", account.PrivateKeyPassphrase)
	d.Set("username", account.Username)

	d.SetId(account.GetID())
}

func getSSHKeyAccountDataSchema() map[string]*schema.Schema {
	schemaMap := getAccountDataSchema()
	schemaMap["account_type"] = &schema.Schema{
		Optional: true,
		Default:  "SshKeyPair",
		Type:     schema.TypeString,
	}
	schemaMap["username"] = &schema.Schema{
		Computed:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	schemaMap["passphrase"] = &schema.Schema{
		Computed:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	return schemaMap
}

func getSSHKeyAccountSchema() map[string]*schema.Schema {
	schemaMap := getAccountSchema()
	schemaMap["account_type"] = &schema.Schema{
		Optional: true,
		Default:  "SshKeyPair",
		Type:     schema.TypeString,
	}
	schemaMap["username"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	schemaMap["passphrase"] = &schema.Schema{
		Optional:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	return schemaMap
}
