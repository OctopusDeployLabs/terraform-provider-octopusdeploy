package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandTokenAccount(d *schema.ResourceData) *octopusdeploy.TokenAccount {
	name := d.Get("name").(string)
	token := octopusdeploy.NewSensitiveValue(d.Get("token").(string))

	account, _ := octopusdeploy.NewTokenAccount(name, token)
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

func setTokenAccount(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.TokenAccount) {
	setAccount(ctx, d, account)

	d.Set("account_type", "Token")
	d.Set("token", account.Token.NewValue)

	d.SetId(account.GetID())
}

func getTokenAccountSchema() map[string]*schema.Schema {
	schemaMap := getAccountSchema()
	schemaMap["account_type"] = &schema.Schema{
		Optional: true,
		Default:  "Token",
		Type:     schema.TypeString,
	}
	schemaMap["token"] = &schema.Schema{
		Required:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	return schemaMap
}
