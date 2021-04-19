package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVariable() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing variables.",
		ReadContext: dataSourceVariableReadByName,
		Schema:      getVariableDataSchema(),
	}
}

func dataSourceVariableReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectID := d.Get("project_id")
	name := d.Get("name")
	scope := octopusdeploy.VariableScope{}

	if v, ok := d.GetOk("scope"); ok {
		scope = expandVariableScope(v)
	}

	client := m.(*octopusdeploy.Client)
	variables, err := client.Variables.GetByName(projectID.(string), name.(string), scope)
	if err != nil {
		return diag.Errorf("error reading variable from project %s with name %s: %s", projectID, name, err.Error())
	}
	if variables == nil {
		return nil
	}
	if len(variables) > 1 {
		return diag.Errorf("found %v variables for project %s with name %s, should match exactly 1", len(variables), projectID, name)
	}

	d.SetId(variables[0].ID)
	d.Set("name", variables[0].Name)
	d.Set("type", variables[0].Type)
	d.Set("value", variables[0].Value)
	d.Set("description", variables[0].Description)

	return nil
}
