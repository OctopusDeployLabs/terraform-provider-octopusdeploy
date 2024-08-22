package schemas

import (
	"fmt"
	"strings"

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

func MapFromVariableScope(variableScope variables.VariableScope) attr.Value {
	if variableScope.IsEmpty() {
		return nil
	}

	flattenedScopes := map[string]attr.Value{
		variableScopeFieldNames.Actions:      util.Ternary(variableScope.Actions != nil && len(variableScope.Actions) > 0, util.FlattenStringList(variableScope.Actions), types.ListNull(types.StringType)),
		variableScopeFieldNames.Channels:     util.Ternary(variableScope.Channels != nil, util.FlattenStringList(variableScope.Channels), types.ListNull(types.StringType)),
		variableScopeFieldNames.Environments: util.Ternary(variableScope.Environments != nil, util.FlattenStringList(variableScope.Environments), types.ListNull(types.StringType)),
		variableScopeFieldNames.Machines:     util.Ternary(variableScope.Machines != nil, util.FlattenStringList(variableScope.Machines), types.ListNull(types.StringType)),
		variableScopeFieldNames.Processes:    util.Ternary(variableScope.ProcessOwners != nil, util.FlattenStringList(variableScope.ProcessOwners), types.ListNull(types.StringType)),
		variableScopeFieldNames.Roles:        util.Ternary(variableScope.Roles != nil, util.FlattenStringList(variableScope.Roles), types.ListNull(types.StringType)),
		variableScopeFieldNames.TenantTags:   util.Ternary(variableScope.TenantTags != nil, util.FlattenStringList(variableScope.TenantTags), types.ListNull(types.StringType)),
	}

	return types.ObjectValueMust(
		VariableScopeObjectType(),
		flattenedScopes,
	)
}

func MapToVariableScope(variableScope types.List) variables.VariableScope {
	if variableScope.IsNull() || len(variableScope.Elements()) == 0 {
		return variables.VariableScope{}
	}

	obj := variableScope.Elements()[0].(types.Object)
	attrs := obj.Attributes()

	scopes := variables.VariableScope{}

	if attrs[variableScopeFieldNames.Actions] != nil {
		scopes.Actions = util.ExpandStringList(attrs[variableScopeFieldNames.Actions].(types.List))
	}
	if attrs[variableScopeFieldNames.Channels] != nil {
		scopes.Channels = util.ExpandStringList(attrs[variableScopeFieldNames.Channels].(types.List))
	}
	if attrs[variableScopeFieldNames.Environments] != nil {
		scopes.Environments = util.ExpandStringList(attrs[variableScopeFieldNames.Environments].(types.List))
	}
	if attrs[variableScopeFieldNames.Machines] != nil {
		scopes.Machines = util.ExpandStringList(attrs[variableScopeFieldNames.Machines].(types.List))
	}
	if attrs[variableScopeFieldNames.Processes] != nil {
		scopes.ProcessOwners = util.ExpandStringList(attrs[variableScopeFieldNames.Processes].(types.List))
	}
	if attrs[variableScopeFieldNames.Roles] != nil {
		scopes.Roles = util.ExpandStringList(attrs[variableScopeFieldNames.Roles].(types.List))
	}
	if attrs[variableScopeFieldNames.TenantTags] != nil {
		scopes.TenantTags = util.ExpandStringList(attrs[variableScopeFieldNames.TenantTags].(types.List))
	}

	return scopes
}

func getVariableScopeResourceSchema() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				variableScopeFieldNames.Actions:      getVariableScopeFieldResourceSchema(variableScopeFieldNames.Actions),
				variableScopeFieldNames.Channels:     getVariableScopeFieldResourceSchema(variableScopeFieldNames.Channels),
				variableScopeFieldNames.Environments: getVariableScopeFieldResourceSchema(variableScopeFieldNames.Environments),
				variableScopeFieldNames.Machines:     getVariableScopeFieldResourceSchema(variableScopeFieldNames.Machines),
				variableScopeFieldNames.Processes:    getVariableScopeFieldResourceSchema(variableScopeFieldNames.Processes),
				variableScopeFieldNames.Roles:        getVariableScopeFieldResourceSchema(variableScopeFieldNames.Roles),
				variableScopeFieldNames.TenantTags:   getVariableScopeFieldResourceSchema(variableScopeFieldNames.TenantTags),
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func getVariableScopeFieldResourceSchema(scopeDescription string) resourceSchema.ListAttribute {
	return resourceSchema.ListAttribute{
		Description: fmt.Sprintf("A list of %s that are scoped to this variable value.", strings.ReplaceAll(scopeDescription, "_", " ")),
		Optional:    true,
		ElementType: basetypes.StringType{},
	}
}

func getVariableScopeDatasourceSchema() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Description: "As variable names can appear more than once under different scopes, a VariableScope must also be provided",
		Computed:    true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				variableScopeFieldNames.Actions:      getVariableScopeFieldDatasourceSchema(variableScopeFieldNames.Actions),
				variableScopeFieldNames.Channels:     getVariableScopeFieldDatasourceSchema(variableScopeFieldNames.Channels),
				variableScopeFieldNames.Environments: getVariableScopeFieldDatasourceSchema(variableScopeFieldNames.Environments),
				variableScopeFieldNames.Machines:     getVariableScopeFieldDatasourceSchema(variableScopeFieldNames.Machines),
				variableScopeFieldNames.Processes:    getVariableScopeFieldDatasourceSchema(variableScopeFieldNames.Processes),
				variableScopeFieldNames.Roles:        getVariableScopeFieldDatasourceSchema(variableScopeFieldNames.Roles),
				variableScopeFieldNames.TenantTags:   getVariableScopeFieldDatasourceSchema(variableScopeFieldNames.TenantTags),
			},
		},
		Validators: []validator.List{
			listvalidator.IsRequired(),
			listvalidator.SizeAtMost(1),
		},
	}
}

func getVariableScopeFieldDatasourceSchema(scopeDescription string) datasourceSchema.ListAttribute {
	return datasourceSchema.ListAttribute{
		Description: fmt.Sprintf("A list of %s that are scoped to this variable value.", strings.ReplaceAll(scopeDescription, "_", " ")),
		Optional:    true,
		ElementType: basetypes.StringType{},
	}
}
