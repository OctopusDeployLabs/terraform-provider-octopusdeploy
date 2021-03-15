package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// expandVariableScope converts a Terraform ResourceData into an OctopusDeploy VariableScope
func expandVariableScope(d *schema.ResourceData) *octopusdeploy.VariableScope {
	// Get the schema set. We specify a MaxItems of 1, so we will only ever have zero or one items
	// in our list.
	tfSchemaSetInterface, ok := d.GetOk("scope")
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
	newScope.Actions = getSliceFromTerraformTypeList(tfSchemaList["actions"])
	newScope.Channels = getSliceFromTerraformTypeList(tfSchemaList["channels"])
	newScope.Environments = getSliceFromTerraformTypeList(tfSchemaList["environments"])
	newScope.Machines = getSliceFromTerraformTypeList(tfSchemaList["machines"])
	newScope.ParentDeployments = getSliceFromTerraformTypeList(tfSchemaList["parent_deployments"])
	newScope.Private = getSliceFromTerraformTypeList(tfSchemaList["private"])
	newScope.ProcessOwners = getSliceFromTerraformTypeList(tfSchemaList["process_owners"])
	newScope.Projects = getSliceFromTerraformTypeList(tfSchemaList["projects"])
	newScope.Roles = getSliceFromTerraformTypeList(tfSchemaList["roles"])
	newScope.TargetRoles = getSliceFromTerraformTypeList(tfSchemaList["target_roles"])
	newScope.Tenants = getSliceFromTerraformTypeList(tfSchemaList["tenants"])
	newScope.TenantTags = getSliceFromTerraformTypeList(tfSchemaList["tenant_tags"])
	newScope.Triggers = getSliceFromTerraformTypeList(tfSchemaList["triggers"])
	newScope.Users = getSliceFromTerraformTypeList(tfSchemaList["users"])
	return &newScope
}

func flattenVariableScope(variableScope *octopusdeploy.VariableScope) []interface{} {
	if variableScope == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"actions":            variableScope.Actions,
		"channels":           variableScope.Channels,
		"environments":       variableScope.Environments,
		"machines":           variableScope.Machines,
		"parent_deployments": variableScope.ParentDeployments,
		"private":            variableScope.Private,
		"process_owners":     variableScope.ProcessOwners,
		"projects":           variableScope.Projects,
		"roles":              variableScope.Roles,
		"target_roles":       variableScope.TargetRoles,
		"tenants":            variableScope.Tenants,
		"tenant_tags":        variableScope.Tenants,
		"triggers":           variableScope.Triggers,
		"users":              variableScope.Users,
	}}
}

func getVariableScopeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"actions": {
			Computed:    true,
			Description: "A list of actions that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"channels": {
			Computed:    true,
			Description: "A list of channels that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"environments": {
			Computed:    true,
			Description: "A list of environments that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"machines": {
			Computed:    true,
			Description: "A list of machines that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"parent_deployments": {
			Computed: true,
			Description: "A list of parent deployments that are scoped to this variable value.",
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"private": {
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"process_owners": {
			Computed:    true,
			Description: "A list of process owners that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"projects": {
			Computed:    true,
			Description: "A list of projects that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"roles": {
			Computed:    true,
			Description: "A list of roles that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"target_roles": {
			Computed:    true,
			Description: "A list of target roles that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"tenants": {
			Computed:    true,
			Description: "A list of tenants that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"tenant_tags": {
			Computed:    true,
			Description: "A list of tenant tags that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"triggers": {
			Computed:    true,
			Description: "A list of triggers that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"users": {
			Computed:    true,
			Description: "A list of users that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}
