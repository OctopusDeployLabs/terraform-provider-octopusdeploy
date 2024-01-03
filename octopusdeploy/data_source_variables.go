package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
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
	spaceID := d.Get("space_id").(string)

	ownerID := d.Get("owner_id").(string)

	name := d.Get("name").(string)

	scope := expandVariableScope(d.Get("scope"))

	client := m.(*client.Client)
	variables, err := variables.GetByName(client, spaceID, ownerID, name, &scope)
	if err != nil {
		return diag.Errorf("error reading variable with owner ID %s with name %s: %s", ownerID, name, err.Error())
	}
	if variables == nil {
		return nil
	}

	if len(variables) != 1 {
		return diag.Errorf("error could not find variable by name, expected to find 1 variable but got %d", len(variables))
	}

	setVariable(ctx, d, variables[0])

	return nil
}
