package util

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetStringSlice(list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	result := make([]string, 0, len(list.Elements()))
	for _, element := range list.Elements() {
		if str, ok := element.(types.String); ok {
			result = append(result, str.ValueString())
		}
	}
	return result
}

func ExpandStringList(list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	result := make([]string, 0, len(list.Elements()))
	for _, elem := range list.Elements() {
		if str, ok := elem.(types.String); ok {
			result = append(result, str.ValueString())
		}
	}
	return result
}

func FlattenStringList(list []string) types.List {
	elements := make([]attr.Value, 0, len(list))
	for _, s := range list {
		elements = append(elements, types.StringValue(s))
	}
	return types.ListValueMust(types.StringType, elements)
}

func Ternary(condition bool, whenTrue, whenFalse attr.Value) attr.Value {
	if condition {
		return whenTrue
	}
	return whenFalse
}
