package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAccount() *schema.Resource {
	resourceAccountSchema := getAccountResourceSchema()
	resourceAccountSchema["access_key"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	resourceAccountSchema["active_directory_endpoint_base_uri"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	resourceAccountSchema["azure_environment"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	resourceAccountSchema["client_id"] = &schema.Schema{
		Optional:         true,
		Type:             schema.TypeString,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	resourceAccountSchema["client_secret"] = &schema.Schema{
		Optional:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	resourceAccountSchema["key"] = &schema.Schema{
		Optional:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	resourceAccountSchema["password"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	resourceAccountSchema["passphrase"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}
	resourceAccountSchema["secret_key"] = &schema.Schema{
		Optional:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	resourceAccountSchema["token"] = &schema.Schema{
		Optional:  true,
		Sensitive: true,
		Type:      schema.TypeString,
	}
	resourceAccountSchema["resource_management_endpoint_base_uri"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	resourceAccountSchema["subscription_id"] = &schema.Schema{
		Optional:         true,
		Type:             schema.TypeString,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	resourceAccountSchema["tenant_id"] = &schema.Schema{
		Optional:         true,
		Type:             schema.TypeString,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	resourceAccountSchema["username"] = &schema.Schema{
		Optional: true,
		Type:     schema.TypeString,
	}

	return &schema.Resource{
		CreateContext: resourceAccountCreate,
		DeleteContext: resourceAccountDeleteCommon,
		Importer:      getImporter(),
		ReadContext:   resourceAccountRead,
		Schema:        getAccountResourceSchema(),
		UpdateContext: resourceAccountUpdate,
	}
}

func resourceAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountResource := expandAccountResource(d)

	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.Add(accountResource)
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccountResource(account)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenAccountResource(ctx, d, accountResource)
	return nil
}

func resourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if account == nil {
		d.SetId("")
		return nil
	}

	accountResource := account.(*octopusdeploy.AccountResource)

	flattenAccountResource(ctx, d, accountResource)
	return nil
}

func resourceAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountResource := expandAccountResource(d)

	client := m.(*octopusdeploy.Client)
	account, err := client.Accounts.Update(accountResource)
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource, err = octopusdeploy.ToAccountResource(account)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenAccountResource(ctx, d, accountResource)
	return nil
}
