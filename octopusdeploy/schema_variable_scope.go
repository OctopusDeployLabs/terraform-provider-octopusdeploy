package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandVariableScope(flattenedVariableScope interface{}) octopusdeploy.VariableScope {
	if flattenedVariableScope == nil {
		return octopusdeploy.VariableScope{}
	}

	list := flattenedVariableScope.([]interface{})
	if len(list) == 0 {
		return octopusdeploy.VariableScope{}
	}

	flattenedMap := list[0].(map[string]interface{})

	variableScope := octopusdeploy.VariableScope{
		Machines:          getSliceFromTerraformTypeList(flattenedMap["machines"]),
		ParentDeployments: getSliceFromTerraformTypeList(flattenedMap["parent_deployments"]),
		Private:           getSliceFromTerraformTypeList(flattenedMap["private"]),
		ProcessOwners:     getSliceFromTerraformTypeList(flattenedMap["process_owners"]),
		Projects:          getSliceFromTerraformTypeList(flattenedMap["projects"]),
		Roles:             getSliceFromTerraformTypeList(flattenedMap["roles"]),
		TargetRoles:       getSliceFromTerraformTypeList(flattenedMap["target_roles"]),
		Tenants:           getSliceFromTerraformTypeList(flattenedMap["tenants"]),
		TenantTags:        getSliceFromTerraformTypeList(flattenedMap["tenant_tags"]),
		Triggers:          getSliceFromTerraformTypeList(flattenedMap["triggers"]),
		Users:             getSliceFromTerraformTypeList(flattenedMap["users"]),
	}

	if v, ok := flattenedMap["actions"]; ok {
		variableScope.Actions = getSliceFromTerraformTypeList(v)
	}

	if v, ok := flattenedMap["channels"]; ok {
		variableScope.Channels = getSliceFromTerraformTypeList(v)
	}

	if v, ok := flattenedMap["environments"]; ok {
		variableScope.Environments = getSliceFromTerraformTypeList(v)
	}

	return variableScope
}

func flattenVariableScope(scope octopusdeploy.VariableScope) []interface{} {
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

	if len(scope.ParentDeployments) > 0 {
		flattenedScope["parent_deployments"] = scope.ParentDeployments
	}

	if len(scope.ProcessOwners) > 0 {
		flattenedScope["private"] = scope.ProcessOwners
	}

	if len(scope.Machines) > 0 {
		flattenedScope["process_owners"] = scope.Machines
	}

	if len(scope.Projects) > 0 {
		flattenedScope["projects"] = scope.Projects
	}

	if len(scope.Roles) > 0 {
		flattenedScope["roles"] = scope.Roles
	}

	if len(scope.TargetRoles) > 0 {
		flattenedScope["target_roles"] = scope.TargetRoles
	}

	if len(scope.Tenants) > 0 {
		flattenedScope["tenants"] = scope.Tenants
	}

	if len(scope.TenantTags) > 0 {
		flattenedScope["tenant_tags"] = scope.TenantTags
	}

	if len(scope.Triggers) > 0 {
		flattenedScope["triggers"] = scope.Triggers
	}

	if len(scope.Users) > 0 {
		flattenedScope["users"] = scope.Users
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
		"parent_deployments": {
			Description: "A list of parent deployments that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"private": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"process_owners": {
			Description: "A list of process owners that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"projects": {
			Description: "A list of projects that are scoped to this variable value.",
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
		"target_roles": {
			Description: "A list of target roles that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"tenants": {
			Description: "A list of tenants that are scoped to this variable value.",
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
		"triggers": {
			Description: "A list of triggers that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"users": {
			Description: "A list of users that are scoped to this variable value.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}
