package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataVariable() *schema.Resource {
	return &schema.Resource{
		Read: dataVariableReadByName,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scope": schemaVariableScope,
		},
	}
}

var schemaVariableScopeValue = &schema.Schema{
	Type: schema.TypeList,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Optional: true,
}

var schemaVariableScope = &schema.Schema{
	Type:     schema.TypeSet,
	MaxItems: 1,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"environments": schemaVariableScopeValue,
			"machines":     schemaVariableScopeValue,
			"actions":      schemaVariableScopeValue,
			"roles":        schemaVariableScopeValue,
			"channels":     schemaVariableScopeValue,
			"tenant_tags":  schemaVariableScopeValue,
		},
	},
}

// tfVariableScopetoODVariableScope converts a Terraform ResourceData into an OctopusDeploy VariableScope
func tfVariableScopetoODVariableScope(d *schema.ResourceData) *octopusdeploy.VariableScope {
	//Get the schema set. We specify a MaxItems of 1, so we will only ever have zero or one items
	//in our list.
	tfSchemaSetInterface, ok := d.GetOk("scope")
	if !ok {
		return nil
	}

	tfSchemaSet := tfSchemaSetInterface.(*schema.Set)
	if len(tfSchemaSet.List()) == 0 {
		return nil
	}

	//Get the first element in the list, which is a map of the interfaces
	tfSchemaList := tfSchemaSet.List()[0].(map[string]interface{})

	//Use the getSliceFromTerraformTypeList helper to convert the data from the map into []string and
	//assign as the variable scopes we need
	var newScope octopusdeploy.VariableScope
	newScope.Environment = getSliceFromTerraformTypeList(tfSchemaList["environments"])
	newScope.Action = getSliceFromTerraformTypeList(tfSchemaList["actions"])
	newScope.Role = getSliceFromTerraformTypeList(tfSchemaList["roles"])
	newScope.Channel = getSliceFromTerraformTypeList(tfSchemaList["channels"])
	newScope.Machine = getSliceFromTerraformTypeList(tfSchemaList["machines"])
	newScope.TenantTag = getSliceFromTerraformTypeList(tfSchemaList["tenant_tags"])

	return &newScope
}

func dataVariableReadByName(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	varProject := d.Get("project_id")
	varName := d.Get("name")
	varScope := tfVariableScopetoODVariableScope(d)

	varItems, err := client.Variable.GetByName(varProject.(string), varName.(string), varScope)

	if err == octopusdeploy.ErrItemNotFound {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading variable from project %s with name %s: %s", varProject, varName, err.Error())
	}

	if len(varItems) > 1 {
		return fmt.Errorf("found %v variables for project %s with name %s, should match exactly 1", len(varItems), varProject, varName)
	}

	d.SetId(varItems[0].ID)
	d.Set("name", varItems[0].Name)
	d.Set("type", varItems[0].Type)
	d.Set("value", varItems[0].Value)
	d.Set("description", varItems[0].Description)

	return nil
}
