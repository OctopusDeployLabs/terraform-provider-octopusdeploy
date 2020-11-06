package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTokenAccount() *schema.Resource {
	dataSourceTokenAccountSchema := map[string]*schema.Schema{
		constAccountType: {
			Default:  constAccountTypeToken,
			Optional: true,
			Type:     schema.TypeString,
		},
		constDescription: &schema.Schema{
			Computed: true,
			Type:     schema.TypeString,
		},
		constEnvironments: {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		constName: &schema.Schema{
			Required: true,
			Type:     schema.TypeString,
		},
		constTenantTags: {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		constTenantedDeploymentParticipation: getTenantedDeploymentSchema(),
		constToken: &schema.Schema{
			Required: true,
			Type:     schema.TypeString,
		},
	}

	return &schema.Resource{
		ReadContext: dataSourceTokenAccountRead,
		Schema:      dataSourceTokenAccountSchema,
	}
}

func dataSourceTokenAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	name := d.Get(constName).(string)
	query := octopusdeploy.AccountsQuery{
		PartialName: name,
		Take:        1,
	}

	accounts, err := client.Accounts.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}
	if accounts == nil || len(accounts.Items) == 0 {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("Unable to retrieve account (partial name: %s)", name))
	}

	account := accounts.Items[0].(*octopusdeploy.TokenAccount)

	d.SetId(account.GetID())
	d.Set(constAccountType, constAccountTypeToken)
	d.Set(constDescription, account.Description)
	d.Set(constEnvironments, account.EnvironmentIDs)
	d.Set(constName, account.GetName())
	d.Set(constSpaceID, account.SpaceID)
	d.Set(constTenants, account.TenantIDs)
	d.Set(constTenantTags, account.TenantTags)
	d.Set(constTenantedDeploymentParticipation, account.TenantedDeploymentMode)

	return nil
}
