package schemas

import (
	"testing"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/resources"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stretchr/testify/require"
)

func TestExpandPromptedDisplaySettingsWithNilInput(t *testing.T) {
	result := MapToDisplaySettings(types.ListNull(types.ObjectType{AttrTypes: VariableDisplaySettingsObjectType()}))
	require.Nil(t, result)
}

func TestExpandPromptedDisplaySettingsWithEmptyInput(t *testing.T) {
	input := types.ListValueMust(types.ObjectType{AttrTypes: VariableDisplaySettingsObjectType()}, []attr.Value{})
	result := MapToDisplaySettings(input)
	require.Nil(t, result)
}

func TestExpandPromptedDisplaySettingsWithCheckbox(t *testing.T) {
	// input := []interface{}{
	// 	map[string]interface{}{
	// 		"control_type": "Checkbox",
	// 	},
	// }
	input := types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{"control_type": types.StringType}},
		[]attr.Value{types.ObjectValueMust(
			map[string]attr.Type{"control_type": types.StringType},
			map[string]attr.Value{"control_type": types.StringValue("Checkbox")},
		)},
	)

	result := MapToDisplaySettings(input)
	require.NotNil(t, result)
	require.Equal(t, resources.ControlTypeCheckbox, result.ControlType)
}

func TestExpandPromptedDisplaySettingsWithSelect(t *testing.T) {
	// input := []interface{}{
	// 	map[string]interface{}{
	// 		"control_type": "Select",
	// 		"select_option": []interface{}{
	// 			map[string]interface{}{
	// 				"display_name": "Name-1",
	// 				"value":        "Value-1",
	// 			},
	// 			map[string]interface{}{
	// 				"display_name": "Name-2",
	// 				"value":        "Value-2",
	// 			},
	// 		},
	// 	},
	// }
	input := types.ListValueMust(
		types.ObjectType{AttrTypes: VariableDisplaySettingsObjectType()},
		[]attr.Value{
			types.ObjectValueMust(
				VariableDisplaySettingsObjectType(),
				map[string]attr.Value{
					"control_type": types.StringValue("Select"),
					"select_option": types.ListValueMust(
						types.ObjectType{AttrTypes: VariableSelectOptionsObjectType()},
						[]attr.Value{
							types.ObjectValueMust(
								VariableSelectOptionsObjectType(),
								map[string]attr.Value{
									"display_name": types.StringValue("Name-1"),
									"value":        types.StringValue("Value-1"),
								},
							),
							types.ObjectValueMust(
								VariableSelectOptionsObjectType(),
								map[string]attr.Value{
									"display_name": types.StringValue("Name-2"),
									"value":        types.StringValue("Value-2"),
								},
							),
						},
					),
				},
			),
		},
	)

	result := MapToDisplaySettings(input)
	require.NotNil(t, result)
	require.Equal(t, resources.ControlTypeSelect, result.ControlType)
	require.NotNil(t, result.SelectOptions)
	require.Len(t, result.SelectOptions, 2)
	require.Equal(t, "Name-1", result.SelectOptions[0].DisplayName)
	require.Equal(t, "Value-1", result.SelectOptions[0].Value)
	require.Equal(t, "Name-2", result.SelectOptions[1].DisplayName)
	require.Equal(t, "Value-2", result.SelectOptions[1].Value)
}

func TestExpandPromptedVariableSettingsWithNilInput(t *testing.T) {
	// result := expandPromptedVariableSettings(nil)
	result := MapToVariablePromptOptions(types.ListNull(types.ObjectType{AttrTypes: VariablePromptOptionsObjectType()}))
	require.Nil(t, result)
}

func TestExpandPromptedVariableSettingsWithEmptyInput(t *testing.T) {
	// input := []interface{}{}
	// result := expandPromptedVariableSettings(input)
	input := types.ListValueMust(types.ObjectType{AttrTypes: VariablePromptOptionsObjectType()}, []attr.Value{})
	result := MapToVariablePromptOptions(input)
	require.Nil(t, result)
}

func TestExpandSelectOptionsWithNilInput(t *testing.T) {
	// result := expandSelectOptions(nil)
	result := MapToSelectOptions(types.ListNull(types.ObjectType{AttrTypes: VariableSelectOptionsObjectType()}))
	require.Nil(t, result)
}

func TestExpandSelectOptionsWithEmptyInput(t *testing.T) {
	// input := []interface{}{}
	// result := expandSelectOptions(input)
	input := types.ListValueMust(types.ObjectType{AttrTypes: VariableSelectOptionsObjectType()}, []attr.Value{})
	result := MapToSelectOptions(input)
	require.Nil(t, result)
}

func TestFlattenPromptedVariableDisplaySettingsWithNilInput(t *testing.T) {
	result := MapFromDisplaySettings(nil)
	require.Empty(t, result)
}

func TestFlattenPromptedVariableSettingsWithNilInput(t *testing.T) {
	result := MapFromVariablePromptOptions(nil)
	require.Empty(t, result)
}

func TestFlattenSelectOptionsWithNilInput(t *testing.T) {
	result := MapFromSelectOptions(nil)
	require.Empty(t, result)
}
