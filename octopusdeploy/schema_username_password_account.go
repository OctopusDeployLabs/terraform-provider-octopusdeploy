package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandUsernamePasswordAccount(d *schema.ResourceData) *octopusdeploy.UsernamePasswordAccount {
	name := d.Get("name").(string)

	account, _ := octopusdeploy.NewUsernamePasswordAccount(name)
	account.ID = d.Id()

	account.Password = octopusdeploy.NewSensitiveValue(d.Get("password").(string))

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
		account.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenants"); ok {
		account.TenantIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("username"); ok {
		account.Username = v.(string)
	}

	return account
}

func setUsernamePasswordAccount(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.UsernamePasswordAccount) {
	setAccount(ctx, d, account)

	d.Set("account_type", "UsernamePassword")
	d.Set("username", account.Username)
	d.Set("password", account.Password.NewValue)

	d.SetId(account.GetID())
}

func getUsernamePasswordAccountSchema() map[string]*schema.Schema {
	schemaMap := getAccountSchema()
	schemaMap["account_type"] = &schema.Schema{
		Default:  "UsernamePassword",
		Optional: true,
		Type:     schema.TypeString,
	}
	schemaMap["password"] = &schema.Schema{
		Sensitive:    true,
		Optional:     true,
		Type:         schema.TypeString,
		ValidateFunc: validation.StringIsNotEmpty,
	}
	schemaMap["username"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	return schemaMap
}
