package schemas

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var variableScopeFieldNames = struct {
	Actions      string
	Channels     string
	Environments string
	Machines     string
	Processes    string
	Roles        string
	TenantTags   string
}{
	Actions:      "actions",
	Channels:     "channels",
	Environments: "environments",
	Machines:     "machines",
	Processes:    "processes",
	Roles:        "roles",
	TenantTags:   "tenant_tags",
}

func VariableScopeObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		variableScopeFieldNames.Actions:      types.ListType{ElemType: types.StringType},
		variableScopeFieldNames.Channels:     types.ListType{ElemType: types.StringType},
		variableScopeFieldNames.Environments: types.ListType{ElemType: types.StringType},
		variableScopeFieldNames.Machines:     types.ListType{ElemType: types.StringType},
		variableScopeFieldNames.Processes:    types.ListType{ElemType: types.StringType},
		variableScopeFieldNames.Roles:        types.ListType{ElemType: types.StringType},
		variableScopeFieldNames.TenantTags:   types.ListType{ElemType: types.StringType},
	}
}

func FlattenVariableScopes(variableScopes variables.VariableScope) attr.Value {
	if variableScopes.IsEmpty() {
		return types.ObjectNull(VariableScopeObjectType())
	}

	flattenedScopes := map[string]attr.Value{}
	flattenedScopes[variableScopeFieldNames.Actions] = util.Ternary(variableScopes.Actions != nil && len(variableScopes.Actions) > 0, util.FlattenStringList(variableScopes.Actions), types.ListNull(types.StringType))
	flattenedScopes[variableScopeFieldNames.Channels] = util.Ternary(variableScopes.Channels != nil, util.FlattenStringList(variableScopes.Channels), types.ListNull(types.StringType))
	flattenedScopes[variableScopeFieldNames.Environments] = util.Ternary(variableScopes.Environments != nil, util.FlattenStringList(variableScopes.Environments), types.ListNull(types.StringType))
	flattenedScopes[variableScopeFieldNames.Machines] = util.Ternary(variableScopes.Machines != nil, util.FlattenStringList(variableScopes.Machines), types.ListNull(types.StringType))
	flattenedScopes[variableScopeFieldNames.Processes] = util.Ternary(variableScopes.ProcessOwners != nil, util.FlattenStringList(variableScopes.ProcessOwners), types.ListNull(types.StringType))
	flattenedScopes[variableScopeFieldNames.Roles] = util.Ternary(variableScopes.Roles != nil, util.FlattenStringList(variableScopes.Roles), types.ListNull(types.StringType))
	flattenedScopes[variableScopeFieldNames.TenantTags] = util.Ternary(variableScopes.TenantTags != nil, util.FlattenStringList(variableScopes.TenantTags), types.ListNull(types.StringType))

	return types.ObjectValueMust(
		VariableScopeObjectType(),
		flattenedScopes,
	)
}

func ExpandVariableScopes(flattenedVariableScopes types.List) variables.VariableScope {
	if flattenedVariableScopes.IsNull() || len(flattenedVariableScopes.Elements()) == 0 {
		return variables.VariableScope{}
	}

	obj := flattenedVariableScopes.Elements()[0].(types.Object)
	attrs := obj.Attributes()
	scopes := variables.VariableScope{}

	scopes.Actions = util.ExpandStringList(attrs[variableScopeFieldNames.Actions].(types.List))
	scopes.Channels = util.ExpandStringList(attrs[variableScopeFieldNames.Channels].(types.List))
	scopes.Environments = util.ExpandStringList(attrs[variableScopeFieldNames.Environments].(types.List))
	scopes.Machines = util.ExpandStringList(attrs[variableScopeFieldNames.Machines].(types.List))
	scopes.ProcessOwners = util.ExpandStringList(attrs[variableScopeFieldNames.Processes].(types.List))
	scopes.Roles = util.ExpandStringList(attrs[variableScopeFieldNames.Roles].(types.List))
	scopes.TenantTags = util.ExpandStringList(attrs[variableScopeFieldNames.TenantTags].(types.List))

	return scopes
}

func getVariableScopeResourceSchema() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				variableScopeFieldNames.Actions:      getVariableScopeItemResourceSchema(variableScopeFieldNames.Actions),
				variableScopeFieldNames.Channels:     getVariableScopeItemResourceSchema(variableScopeFieldNames.Channels),
				variableScopeFieldNames.Environments: getVariableScopeItemResourceSchema(variableScopeFieldNames.Environments),
				variableScopeFieldNames.Machines:     getVariableScopeItemResourceSchema(variableScopeFieldNames.Machines),
				variableScopeFieldNames.Processes:    getVariableScopeItemResourceSchema(variableScopeFieldNames.Processes),
				variableScopeFieldNames.Roles:        getVariableScopeItemResourceSchema(variableScopeFieldNames.Roles),
				variableScopeFieldNames.TenantTags:   getVariableScopeItemResourceSchema(variableScopeFieldNames.TenantTags),
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func getVariableScopeItemResourceSchema(scopeDescription string) resourceSchema.ListAttribute {
	return resourceSchema.ListAttribute{
		Description: fmt.Sprintf("A list of %s that are scoped to this variable value.", scopeDescription),
		Optional:    true,
		ElementType: basetypes.StringType{},
	}
}

func getVariableScopeDatasourceSchema() datasourceSchema.ListNestedBlock {
	return datasourceSchema.ListNestedBlock{
		NestedObject: datasourceSchema.NestedBlockObject{
			Attributes: map[string]datasourceSchema.Attribute{
				variableScopeFieldNames.Actions:      getVariableScopeItemDatasourceSchema(variableScopeFieldNames.Actions),
				variableScopeFieldNames.Channels:     getVariableScopeItemDatasourceSchema(variableScopeFieldNames.Channels),
				variableScopeFieldNames.Environments: getVariableScopeItemDatasourceSchema(variableScopeFieldNames.Environments),
				variableScopeFieldNames.Machines:     getVariableScopeItemDatasourceSchema(variableScopeFieldNames.Machines),
				variableScopeFieldNames.Processes:    getVariableScopeItemDatasourceSchema(variableScopeFieldNames.Processes),
				variableScopeFieldNames.Roles:        getVariableScopeItemDatasourceSchema(variableScopeFieldNames.Roles),
				variableScopeFieldNames.TenantTags:   getVariableScopeItemDatasourceSchema(variableScopeFieldNames.TenantTags),
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func getVariableScopeItemDatasourceSchema(scopeDescription string) datasourceSchema.ListAttribute {
	return datasourceSchema.ListAttribute{
		Description: fmt.Sprintf("A list of %s that are scoped to this variable value.", scopeDescription),
		Optional:    true,
		ElementType: basetypes.StringType{},
	}
}
