package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnvironmentRead,
		Schema:      getEnvironmentDataSchema(),
	}
}

func dataSourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.EnvironmentsQuery{
		IDs:         expandArray(d.Get("ids").([]interface{})),
		Name:        d.Get("name").(string),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	client := m.(*octopusdeploy.Client)
	environments, err := client.Environments.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedEnvironments := []interface{}{}
	for _, environment := range environments.Items {
		flattenedEnvironments = append(flattenedEnvironments, flattenEnvironment(environment))
	}

	d.Set("environments", flattenedEnvironments)
	d.SetId("Environments " + time.Now().UTC().String())

	return nil
}
