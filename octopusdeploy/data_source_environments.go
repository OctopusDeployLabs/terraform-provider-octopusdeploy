package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceEnvironment() *schema.Resource {
	dataSourceEnvironmentSchema := map[string]*schema.Schema{
		"name": &schema.Schema{
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
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

	setEnvironment(ctx, d, environment)
	return nil
}
