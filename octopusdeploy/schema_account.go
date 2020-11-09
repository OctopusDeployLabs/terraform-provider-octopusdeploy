package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenAccount(ctx context.Context, d *schema.ResourceData, account octopusdeploy.IAccount) {
	d.Set("account_type", account.GetAccountType())
	d.Set("description", account.GetDescription())
	d.Set("environments", account.GetEnvironmentIDs())
	d.Set("modified_by", account.GetModifiedBy())

	if modifiedOn := account.GetModifiedOn(); modifiedOn != nil {
		d.Set("modified_on", modifiedOn.Format(time.RFC3339))
	}

	d.Set("name", account.GetName())
	d.Set("space_id", account.GetSpaceID())
	d.Set("tenanted_deployment_participation", account.GetTenantedDeploymentMode())
	d.Set("tenants", account.GetTenantIDs())
	d.Set("tenant_tags", account.GetTenantTags())
}

func getAccountDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"account_type": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"description": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"environments": {
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Type: schema.TypeList,
		},
		"modified_by": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"modified_on": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"space_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"tenanted_deployment_participation": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"tenants": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Computed: true,
			Type:     schema.TypeList,
		},
		"tenant_tags": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Computed: true,
			Type:     schema.TypeList,
		},
	}
}

func getAccountSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_type": {
			Default:  "None",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validateValueFunc([]string{
				"None",
				"AmazonWebServicesAccount",
				"AzureServicePrincipal",
				"AzureSubscription",
				"SshKeyPair",
				"Token",
				"UsernamePassword",
			}),
		},
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"environments": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		"modified_by": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"modified_on": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"space_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"tenanted_deployment_participation": getTenantedDeploymentSchema(),
		"tenants": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		"tenant_tags": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}
