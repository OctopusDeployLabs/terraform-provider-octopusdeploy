package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAccount() *schema.Resource {
	resourceAccountSchema := getCommonAccountsSchema()
	resourceAccountSchema[constAccountType] = getAccountTypeSchema()
	resourceAccountSchema[constAccessKey] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	resourceAccountSchema[constUsername] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	resourceAccountSchema[constPassword] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	resourceAccountSchema[constPassphrase] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	resourceAccountSchema[constSecretKey] = &schema.Schema{
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	}
	resourceAccountSchema[constToken] = &schema.Schema{
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	}
	resourceAccountSchema[constClientID] = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	resourceAccountSchema[constClientSecret] = &schema.Schema{
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	}
	resourceAccountSchema[constTenantID] = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	resourceAccountSchema[constSubscriptionID] = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ValidateDiagFunc: validateDiagFunc(validation.IsUUID),
	}
	resourceAccountSchema[constKey] = &schema.Schema{
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
	}
	resourceAccountSchema[constAzureEnvironment] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	resourceAccountSchema[constResourceManagementEndpointBaseURI] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	resourceAccountSchema[constActiveDirectoryEndpointBaseURI] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	return &schema.Resource{
		CreateContext: resourceAccountCreate,
		DeleteContext: resourceAccountDeleteCommon,
		ReadContext:   resourceAccountRead,
		Schema:        resourceAccountSchema,
		UpdateContext: resourceAccountUpdate,
	}
}

func resourceAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*octopusdeploy.Client)
	account, err := client.Accounts.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	accountResource := account.(*octopusdeploy.AccountResource)
	flattenAccountResource(ctx, d, accountResource)

	return nil
}

func resourceAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	accountResource := buildAccountResource(d)

	client := meta.(*octopusdeploy.Client)
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

func resourceAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	accountResource := buildAccountResource(d)
	accountResource.ID = d.Id()

	client := meta.(*octopusdeploy.Client)
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
