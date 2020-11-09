package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEnvironment() *schema.Resource {
	dataSourceEnvironmentSchema := map[string]*schema.Schema{
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
	}

	return &schema.Resource{
		ReadContext: dataSourceEnvironmentRead,
		Schema:      dataSourceEnvironmentSchema,
	}
}

func dataSourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := d.Get("name").(string)

	client := m.(*octopusdeploy.Client)
	environments, err := client.Environments.GetByName(name)
	if err != nil {
		return diag.FromErr(err)
	}

	environment := environments[0]

	flattenEnvironment(ctx, d, environment)
	return nil
}
