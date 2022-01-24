package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandUsernamePasswordAccount(d *schema.ResourceData) octopusdeploy.IUsernamePasswordAccount {
	name := d.Get("name").(string)

	account, _ := octopusdeploy.NewUsernamePasswordAccount(name)
	account.SetID(d.Id())
	account.SetPassword(octopusdeploy.NewSensitiveValue(d.Get("password").(string)))

	if v, ok := d.GetOk("description"); ok {
		account.SetDescription(v.(string))
	}

	if v, ok := d.GetOk("environments"); ok {
		account.SetEnvironmentIDs(getSliceFromTerraformTypeList(v))
	}

	if v, ok := d.GetOk("space_id"); ok {
		account.SetSpaceID(v.(string))
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		account.SetTenantedDeploymentMode(octopusdeploy.TenantedDeploymentMode(v.(string)))
	}

	if v, ok := d.GetOk("tenants"); ok {
		account.SetTenantIDs(getSliceFromTerraformTypeList(v))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		account.SetTenantTags(getSliceFromTerraformTypeList(v))
	}

	if v, ok := d.GetOk("username"); ok {
		account.SetUsername(v.(string))
	}

	return account
}

func setUsernamePasswordAccount(ctx context.Context, d *schema.ResourceData, account *octopusdeploy.UsernamePasswordAccount) error {
	d.Set("description", account.GetDescription())

	if err := d.Set("environments", account.GetEnvironmentIDs()); err != nil {
		return fmt.Errorf("error setting environments: %s", err)
	}

	d.Set("id", account.GetID())
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

	d.SetId(account.GetID())

	return nil
}

func getUsernamePasswordAccountSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description":                       getDescriptionSchema("username/password account"),
		"environments":                      getEnvironmentsSchema(),
		"id":                                getIDSchema(),
		"name":                              getNameSchema(true),
		"password":                          getPasswordSchema(false),
		"space_id":                          getSpaceIDSchema(),
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants":                           getTenantsSchema(),
		"tenant_tags":                       getTenantTagsSchema(),
		"username":                          getUsernameSchema(true),
	}
}
