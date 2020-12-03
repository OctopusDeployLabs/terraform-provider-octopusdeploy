package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccounts() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing accounts.",
		ReadContext: dataSourceAccountsRead,
		Schema:      getAccountResourceDataSchema(),
	}
}

func dataSourceAccountsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.AccountsQuery{
		AccountType: octopusdeploy.AccountType(d.Get("account_type").(string)),
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	client := m.(*octopusdeploy.Client)
	accounts, err := client.Accounts.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedAccounts := []interface{}{}
	for _, account := range accounts.Items {
		accountResource, err := octopusdeploy.ToAccountResource(account)
		if err != nil {
			return diag.FromErr(err)
		}

		flattenedAccounts = append(flattenedAccounts, flattenAccountResource(accountResource))
	}

	d.Set("accounts", flattenedAccounts)
	d.SetId("Accounts " + time.Now().UTC().String())

	return nil
}
