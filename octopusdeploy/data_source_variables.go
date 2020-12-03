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

// tfVariableScopetoODVariableScope converts a Terraform ResourceData into an OctopusDeploy VariableScope
func tfVariableScopetoODVariableScope(d *schema.ResourceData) *octopusdeploy.VariableScope {
	// Get the schema set. We specify a MaxItems of 1, so we will only ever have zero or one items
	// in our list.
	tfSchemaSetInterface, ok := d.GetOk(constScope)
	if !ok {
		return nil
	}

	tfSchemaSet := tfSchemaSetInterface.(*schema.Set)
	if len(tfSchemaSet.List()) == 0 {
		return nil
	}

	// Get the first element in the list, which is a map of the interfaces
	tfSchemaList := tfSchemaSet.List()[0].(map[string]interface{})

	// Use the getSliceFromTerraformTypeList helper to convert the data from the map into []string and
	// assign as the variable scopes we need
	var newScope octopusdeploy.VariableScope
	newScope.Action = getSliceFromTerraformTypeList(tfSchemaList["actions"])
	newScope.Channel = getSliceFromTerraformTypeList(tfSchemaList["channels"])
	newScope.Environment = getSliceFromTerraformTypeList(tfSchemaList["environments"])
	newScope.Machine = getSliceFromTerraformTypeList(tfSchemaList["machines"])
	newScope.Role = getSliceFromTerraformTypeList(tfSchemaList["roles"])
	newScope.TenantTag = getSliceFromTerraformTypeList(tfSchemaList["tenant_tags"])
	return &newScope
}

func dataSourceVariableReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectID := d.Get("project_id")
	name := d.Get("name")
	scope := tfVariableScopetoODVariableScope(d)

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
