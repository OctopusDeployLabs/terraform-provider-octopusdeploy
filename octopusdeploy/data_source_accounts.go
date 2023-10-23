package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/accounts"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccounts() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Use an account-specific resource instead (i.e. octopusdeploy_aws_account, octopusdeploy_azure_service_principal, octopusdeploy_azure_subscription_account, octopusdeploy_ssh_key_account, octopusdeploy_token_account, octopusdeploy_username_password_account).",
		Description:        "Provides information about existing accounts.",
		ReadContext:        dataSourceAccountsRead,
		Schema:             getAccountResourceDataSchema(),
	}
}

func dataSourceAccountsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := accounts.AccountsQuery{
		AccountType: accounts.AccountType(d.Get("account_type").(string)),
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}
	spaceID := d.Get("space_id").(string)

	client := m.(*client.Client)
	existingAccounts, err := accounts.Get(client, spaceID, &query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedAccounts := []interface{}{}
	for _, account := range existingAccounts.Items {
		accountResource, err := accounts.ToAccountResource(account)
		if err != nil {
			return diag.FromErr(err)
		}

		flattenedAccounts = append(flattenedAccounts, flattenAccountResource(accountResource))
	}

	d.Set("accounts", flattenedAccounts)
	d.SetId("Accounts " + time.Now().UTC().String())

	return nil
}
