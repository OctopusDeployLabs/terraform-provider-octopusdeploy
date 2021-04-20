package octopusdeploy

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/stretchr/testify/assert"
)

func TestExpandVariableScope(t *testing.T) {
	scope := expandVariableScope(nil)
	assert.True(t, scope.IsEmpty())
	assert.Equal(t, octopusdeploy.VariableScope{}, scope)
	assert.Empty(t, scope.Channels)

	flattenedVariableScope := []interface{}{}
	scope = expandVariableScope(flattenedVariableScope)
	assert.True(t, scope.IsEmpty())
	assert.Equal(t, octopusdeploy.VariableScope{}, scope)
	assert.Empty(t, scope.Channels)

	flattenedVariableScope = []interface{}{nil}
	scope = expandVariableScope(flattenedVariableScope)
	assert.True(t, scope.IsEmpty())
	assert.Equal(t, octopusdeploy.VariableScope{}, scope)
	assert.Empty(t, scope.Channels)

	flattenedVariableScope = []interface{}{"foo"}
	scope = expandVariableScope(flattenedVariableScope)
	assert.True(t, scope.IsEmpty())
	assert.Equal(t, octopusdeploy.VariableScope{}, scope)
	assert.Empty(t, scope.Channels)

	flattenedVariableScope = []interface{}{map[string]interface{}{}}
	scope = expandVariableScope(flattenedVariableScope)
	assert.True(t, scope.IsEmpty())
	assert.Equal(t, octopusdeploy.VariableScope{}, scope)
	assert.Empty(t, scope.Channels)

	flattenedVariableScope = []interface{}{map[string]interface{}{
		"actions": []interface{}{"foo"},
	}}
	expectedScope := octopusdeploy.VariableScope{
		Actions: []string{"foo"},
	}
	scope = expandVariableScope(flattenedVariableScope)
	assert.False(t, scope.IsEmpty())
	assert.Equal(t, expectedScope, scope)
	assert.NotEmpty(t, scope.Actions)
	assert.Empty(t, scope.Channels)
}

func TestFlattenVariableScope(t *testing.T) {
	flattenedVariableScope := []interface{}{map[string]interface{}{
		"actions": []interface{}{"foo"},
	}}
	expectedScope := octopusdeploy.VariableScope{
		Actions: []string{"foo"},
	}
	scope := expandVariableScope(flattenedVariableScope)
	assert.False(t, scope.IsEmpty())
	assert.Equal(t, expectedScope, scope)
	assert.NotEmpty(t, scope.Actions)
	assert.Empty(t, scope.Channels)

	flattenedVariableScope = flattenVariableScope(scope)
	assert.NotNil(t, flattenedVariableScope)
	assert.Len(t, flattenedVariableScope, 1)
	assert.Len(t, flattenedVariableScope[0], 1)
}
