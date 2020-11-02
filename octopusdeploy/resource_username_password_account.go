package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceUsernamePassword() *schema.Resource {
	schemaMap := getCommonAccountsSchema()
	schemaMap[constPassword] = &schema.Schema{
		Optional:     true,
		Type:         schema.TypeString,
		ValidateFunc: validation.StringIsNotEmpty,
	}
	schemaMap[constUsername] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}

	return &schema.Resource{
		CreateContext: resourceUsernamePasswordCreate,
		ReadContext:   resourceUsernamePasswordRead,
		UpdateContext: resourceUsernamePasswordUpdate,
		DeleteContext: resourceAccountDeleteCommon,
		Schema:        schemaMap,
	}
}

func resourceUsernamePasswordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	logResource(constAccount, m)

	accountResource := account.(*octopusdeploy.AccountResource)

	d.Set(constDescription, accountResource.Description)
	d.Set(constEnvironments, accountResource.EnvironmentIDs)
	d.Set(constName, accountResource.Name)
	d.Set(constUsername, accountResource.Username)

	if accountResource.ApplicationPassword.HasValue {
		if accountResource.ApplicationPassword.NewValue != nil {
			d.Set(constPassword, accountResource.ApplicationPassword.NewValue)
		}
	}

	return nil
}

func buildUsernamePasswordResource(d *schema.ResourceData) (*octopusdeploy.UsernamePasswordAccount, error) {
	name := d.Get(constName).(string)

	account, err := octopusdeploy.NewUsernamePasswordAccount(name)
	if err != nil {
		return nil, err
	}

	password := d.Get(constPassword).(string)
	privateKey := octopusdeploy.NewSensitiveValue(password)
	account.Password = privateKey

	if v, ok := d.GetOk(constTenantedDeploymentParticipation); ok {
		account.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk(constTenantTags); ok {
		account.TenantTags = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constUsername); ok {
		account.Username = v.(string)
	}

	return account, nil
}

func resourceUsernamePasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, err := buildUsernamePasswordResource(d)
	if err != nil {
		diag.FromErr(err)
	}

	client := m.(*octopusdeploy.Client)
	createdAccount, err := client.Accounts.Add(account)
	if err != nil {
		diag.FromErr(err)
	}

	d.SetId(createdAccount.GetID())

	return nil
}

func resourceUsernamePasswordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	account, err := buildUsernamePasswordResource(d)
	if err != nil {
		diag.FromErr(err)
	}
	account.ID = d.Id() // set ID so Octopus API knows which account to update

	client := m.(*octopusdeploy.Client)
	updatedAccount, err := client.Accounts.Update(account)
	if err != nil {
		diag.FromErr(err)
	}

	d.SetId(updatedAccount.GetID())

	return nil
}
