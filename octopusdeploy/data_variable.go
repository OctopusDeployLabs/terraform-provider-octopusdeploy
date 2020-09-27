package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataVariable() *schema.Resource {
	return &schema.Resource{
		Read: dataVariableReadByName,
		Schema: map[string]*schema.Schema{
			constProjectID: {
				Type:     schema.TypeString,
				Required: true,
			},
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constType: {
				Type:     schema.TypeString,
				Computed: true,
			},
			constValue: {
				Type:     schema.TypeString,
				Computed: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Computed: true,
			},
			constScope: schemaVariableScope,
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
			constEnvironments: schemaVariableScopeValue,
			constMachines:     schemaVariableScopeValue,
			constActions:      schemaVariableScopeValue,
			constRoles:        schemaVariableScopeValue,
			constChannels:     schemaVariableScopeValue,
			constTenantTags:   schemaVariableScopeValue,
		},
	},
}

// tfVariableScopetoODVariableScope converts a Terraform ResourceData into an OctopusDeploy VariableScope
func tfVariableScopetoODVariableScope(d *schema.ResourceData) *model.VariableScope {
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
	var newScope model.VariableScope
	newScope.Environment = getSliceFromTerraformTypeList(tfSchemaList[constEnvironments])
	newScope.Action = getSliceFromTerraformTypeList(tfSchemaList[constActions])
	newScope.Role = getSliceFromTerraformTypeList(tfSchemaList[constRoles])
	newScope.Channel = getSliceFromTerraformTypeList(tfSchemaList[constChannels])
	newScope.Machine = getSliceFromTerraformTypeList(tfSchemaList[constMachines])
	newScope.TenantTag = getSliceFromTerraformTypeList(tfSchemaList[constTenantTags])

	return &newScope
}

func dataVariableReadByName(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	varProject := d.Get(constProjectID)
	varName := d.Get(constName)
	varScope := tfVariableScopetoODVariableScope(d)

	varItems, err := apiClient.Variables.GetByName(varProject.(string), varName.(string), varScope)
	if err != nil {
		return fmt.Errorf("error reading variable from project %s with name %s: %s", varProject, varName, err.Error())
	}
	if varItems == nil {
		return nil
	}
	if len(varItems) > 1 {
		return fmt.Errorf("found %v variables for project %s with name %s, should match exactly 1", len(varItems), varProject, varName)
	}

	d.SetId(varItems[0].ID)
	d.Set(constName, varItems[0].Name)
	d.Set(constType, varItems[0].Type)
	d.Set(constValue, varItems[0].Value)
	d.Set(constDescription, varItems[0].Description)

	return nil
}
