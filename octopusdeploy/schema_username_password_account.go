package octopusdeploy

import (
	"context"
	"time"

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

	if v, ok := d.GetOk("modified_by"); ok {
		account.ModifiedBy = v.(string)
	}

	if v, ok := d.GetOk("modified_on"); ok {
		modifiedOnTime, _ := time.Parse(time.RFC3339, v.(string))
		account.ModifiedOn = &modifiedOnTime
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

func flattenUsernamePasswordAccount(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.UsernamePasswordAccount) {
	flattenAccount(ctx, d, account)

	d.Set("account_type", "UsernamePassword")
	d.Set("username", account.Username)
	d.Set("password", account.Password.NewValue)

	d.SetId(account.GetID())
}

func getUsernamePasswordAccountDataSchema() map[string]*schema.Schema {
	schemaMap := getAccountDataSchema()
	schemaMap["account_type"] = &schema.Schema{
		Optional: true,
		Default:  "UsernamePassword",
		Type:     schema.TypeString,
	}
	schemaMap["password"] = &schema.Schema{
		Computed:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	schemaMap["username"] = &schema.Schema{
		Computed:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	return schemaMap
}

func getUsernamePasswordAccountSchema() map[string]*schema.Schema {
	schemaMap := getAccountSchema()
	schemaMap["account_type"] = &schema.Schema{
		Optional: true,
		Default:  "UsernamePassword",
		Type:     schema.TypeString,
	}
	schemaMap["password"] = &schema.Schema{
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
