package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataAwsAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAwsAccountReadByName,
		Schema: map[string]*schema.Schema{
			constAccessKey: {
				Optional: true,
				Type:     schema.TypeString,
			},
			constAccountType: {
				Default:  constAccountTypeAWS,
				Optional: true,
				Type:     schema.TypeString,
			},
			constDescription: {
				Optional: true,
				Type:     schema.TypeString,
			},
			constEnvironments: {
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Type:     schema.TypeList,
			},
			constName: {
				Required: true,
				Type:     schema.TypeString,
			},
			constSecretKey: {
				Optional: true,
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
		},
	}
}

func dataAwsAccountReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	logResource(constAccount, m)

	account := accounts.Items[0].(*octopusdeploy.AmazonWebServicesAccount)

	d.SetId(account.GetID())
	d.Set(constAccessKey, account.AccessKey)
	d.Set(constAccountType, account.AccountType)
	d.Set(constDescription, account.Description)
	d.Set(constEnvironments, account.EnvironmentIDs)
	d.Set(constName, account.GetName())
	d.Set(constSpaceID, account.SpaceID)
	d.Set(constTenants, account.TenantIDs)
	d.Set(constTenantTags, account.TenantTags)
	d.Set(constTenantedDeploymentParticipation, account.TenantedDeploymentMode)

	return nil
}
