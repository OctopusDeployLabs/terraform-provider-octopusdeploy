package schemas

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestExpandVariableScope(t *testing.T) {
	// scope := ExpandVariableScopes(nil)
	scope := MapToVariableScope(types.ListNull(types.ObjectType{AttrTypes: VariableScopeObjectType()}))
	assert.True(t, scope.IsEmpty())
	assert.Equal(t, variables.VariableScope{}, scope)
	assert.Empty(t, scope.Channels)

	// flattenedVariableScope := []interface{}{}
	flattenedVariableScope := types.ListValueMust(types.ObjectType{AttrTypes: VariableScopeObjectType()}, []attr.Value{})
	scope = MapToVariableScope(flattenedVariableScope)
	assert.True(t, scope.IsEmpty())
	assert.Equal(t, variables.VariableScope{}, scope)
	assert.Empty(t, scope.Channels)

	// flattenedVariableScope = []interface{}{nil}
	flattenedVariableScope = types.ListValueMust(types.ObjectType{AttrTypes: VariableScopeObjectType()}, []attr.Value{types.ObjectNull(VariableScopeObjectType())})
	scope = MapToVariableScope(flattenedVariableScope)
	assert.True(t, scope.IsEmpty())
	assert.Equal(t, variables.VariableScope{}, scope)
	assert.Empty(t, scope.Channels)

	//flattenedVariableScope = []interface{}{"foo"}
	// flattenedVariableScope = types.ListValueMust(types.StringType, []attr.Value{types.StringValue("foo")})
	// scope = MapToVariableScope(flattenedVariableScope)
	// assert.True(t, scope.IsEmpty())
	// assert.Equal(t, variables.VariableScope{}, scope)
	// assert.Empty(t, scope.Channels)

	// flattenedVariableScope = []interface{}{map[string]interface{}{}}
	flattenedVariableScope = types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{}}, []attr.Value{})
	scope = MapToVariableScope(flattenedVariableScope)
	assert.True(t, scope.IsEmpty())
	assert.Equal(t, variables.VariableScope{}, scope)
	assert.Empty(t, scope.Channels)

	// flattenedVariableScope = []interface{}{map[string]interface{}{
	// 	"actions": []interface{}{"foo"},
	// }}
	flattenedVariableScope = types.ListValueMust(
		types.ObjectType{AttrTypes: map[string]attr.Type{"actions": types.ListType{ElemType: types.StringType}}},
		[]attr.Value{basetypes.NewObjectValueMust(
			map[string]attr.Type{"actions": types.ListType{ElemType: types.StringType}},
			map[string]attr.Value{"actions": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("foo")})},
		)})
	expectedScope := variables.VariableScope{
		Actions: []string{"foo"},
	}
	scope = MapToVariableScope(flattenedVariableScope)
	assert.False(t, scope.IsEmpty())
	assert.Equal(t, expectedScope, scope)
	assert.NotEmpty(t, scope.Actions)
	assert.Empty(t, scope.Channels)
}

func TestFlattenVariableScope(t *testing.T) {
	scopes := basetypes.NewObjectValueMust(
		map[string]attr.Type{
			"actions": types.ListType{ElemType: types.StringType},
		},
		map[string]attr.Value{
			"actions": types.ListValueMust(types.StringType, []attr.Value{types.StringValue("foo")}),
		},
	)

	flattenedVariableScope := types.ListValueMust(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"actions": types.ListType{ElemType: types.StringType},
			},
		},
		[]attr.Value{scopes},
	)
	expectedScope := variables.VariableScope{
		Actions: []string{"foo"},
	}

	scope := MapToVariableScope(flattenedVariableScope)
	assert.False(t, scope.IsEmpty())
	assert.Equal(t, expectedScope, scope)
	assert.NotEmpty(t, scope.Actions)
	assert.Empty(t, scope.Channels)

	flattenedVariableScope = types.ListValueMust(types.ObjectType{AttrTypes: VariableScopeObjectType()}, []attr.Value{MapFromVariableScope(scope)})

	assert.NotNil(t, flattenedVariableScope)
	assert.Len(t, flattenedVariableScope.Elements(), 1)
	actionScope := flattenedVariableScope.Elements()[0].(types.Object).Attributes()["actions"].(types.List)
	t.Logf("Action scope: %#v", actionScope)
	assert.Len(t, actionScope.Elements(), 1)
}
