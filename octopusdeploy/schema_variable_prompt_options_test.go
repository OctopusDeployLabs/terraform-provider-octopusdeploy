package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/stretchr/testify/require"
)

func TestExpandPromptedDisplaySettingsWithNilInput(t *testing.T) {
	result := expandPromptedDisplaySettings(nil)
	require.Nil(t, result)
}

func TestExpandPromptedDisplaySettingsWithEmptyInput(t *testing.T) {
	input := []interface{}{}
	result := expandPromptedDisplaySettings(input)
	require.Nil(t, result)
}

func TestExpandPromptedDisplaySettingsWithCheckbox(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"control_type": "Checkbox",
		},
	}
	result := expandPromptedDisplaySettings(input)
	require.NotNil(t, result)
	require.Equal(t, resources.ControlTypeCheckbox, result.ControlType)
}

func TestExpandPromptedDisplaySettingsWithSelect(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"control_type": "Select",
			"select_option": []interface{}{
				map[string]interface{}{
					"display_name": "Name-1",
					"value":        "Value-1",
				},
				map[string]interface{}{
					"display_name": "Name-2",
					"value":        "Value-2",
				},
			},
		},
	}
	result := expandPromptedDisplaySettings(input)
	require.NotNil(t, result)
	require.Equal(t, variables.ControlTypeSelect, result.ControlType)
	require.NotNil(t, result.SelectOptions)
	require.Len(t, result.SelectOptions, 2)
	require.Equal(t, "Name-1", result.SelectOptions[0].DisplayName)
	require.Equal(t, "Value-1", result.SelectOptions[0].Value)
	require.Equal(t, "Name-2", result.SelectOptions[1].DisplayName)
	require.Equal(t, "Value-2", result.SelectOptions[1].Value)
}

func TestExpandPromptedVariableSettingsWithNilInput(t *testing.T) {
	result := expandPromptedVariableSettings(nil)
	require.Nil(t, result)
}

func TestExpandPromptedVariableSettingsWithEmptyInput(t *testing.T) {
	input := []interface{}{}
	result := expandPromptedVariableSettings(input)
	require.Nil(t, result)
}

func TestExpandSelectOptionsWithNilInput(t *testing.T) {
	result := expandSelectOptions(nil)
	require.Nil(t, result)
}

func TestExpandSelectOptionsWithEmptyInput(t *testing.T) {
	input := []interface{}{}
	result := expandSelectOptions(input)
	require.Nil(t, result)
}

func TestFlattenPromptedVariableDisplaySettingsWithNilInput(t *testing.T) {
	result := flattenPromptedVariableDisplaySettings(nil)
	require.Empty(t, result)
}

func TestFlattenPromptedVariableSettingsWithNilInput(t *testing.T) {
	result := flattenPromptedVariableSettings(nil)
	require.Empty(t, result)
}

func TestFlattenSelectOptionsWithNilInput(t *testing.T) {
	result := flattenSelectOptions(nil)
	require.Empty(t, result)
}
