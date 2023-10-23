package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandVariableScope(flattenedVariableScope interface{}) variables.VariableScope {
	if flattenedVariableScope == nil {
		return variables.VariableScope{}
	}

	list := flattenedVariableScope.([]interface{})
	if len(list) == 0 || list[0] == nil {
		return variables.VariableScope{}
	}

	if flattenedMap, ok := list[0].(map[string]interface{}); ok {
		return variables.VariableScope{
			Actions:      getSliceFromTerraformTypeList(flattenedMap["actions"]),
			Channels:     getSliceFromTerraformTypeList(flattenedMap["channels"]),
			Environments: getSliceFromTerraformTypeList(flattenedMap["environments"]),
			Machines:     getSliceFromTerraformTypeList(flattenedMap["machines"]),
			ProcessOwners: getSliceFromTerraformTypeList(flattenedMap["processes"]),
			Roles:        getSliceFromTerraformTypeList(flattenedMap["roles"]),
			TenantTags:   getSliceFromTerraformTypeList(flattenedMap["tenant_tags"]),
		}
	}

	return variables.VariableScope{}
}

func flattenVariableScope(scope variables.VariableScope) []interface{} {
	if scope.IsEmpty() {
		return nil
	}

	flattenedScope := map[string]interface{}{}

	if len(scope.Actions) > 0 {
		flattenedScope["actions"] = scope.Actions
	}

	if len(scope.Channels) > 0 {
		flattenedScope["channels"] = scope.Channels
	}

	if len(scope.Environments) > 0 {
		flattenedScope["environments"] = scope.Environments
	}

	if len(scope.Machines) > 0 {
		flattenedScope["machines"] = scope.Machines
	}

	if len(scope.ProcessOwners) > 0 {
		flattenedScope["processes"] = scope.ProcessOwners
	}

	if len(scope.Roles) > 0 {
		flattenedScope["roles"] = scope.Roles
	}

	if len(scope.TenantTags) > 0 {
		flattenedScope["tenant_tags"] = scope.TenantTags
	}

	return []interface{}{flattenedScope}
}

func getVariableScopeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"actions": {
			Description: "A list of actions that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"channels": {
			Description: "A list of channels that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"environments": {
			Description: "A list of environments that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"machines": {
			Description: "A list of machines that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
        "processes": {
			Description: "A list of processes that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"roles": {
			Description: "A list of roles that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"tenant_tags": {
			Description: "A list of tenant tags that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}
