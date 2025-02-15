package util

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

func GetProviderName() string {
	return "octopusdeploy"
}

func GetTypeName(name string) string {
	return fmt.Sprintf("%s_%s", GetProviderName(), name)
}

func GetResourceSchemaDescription(resourceName string) string {
	return fmt.Sprintf("This resource manages %ss in Octopus Deploy.", resourceName)
}

func GetDataSourceDescription(resourceName string) string {
	return fmt.Sprintf("Provides information about existing %s.", resourceName)
}

func GetStringOrEmpty(tfAttr interface{}) string {
	if tfAttr == nil {
		return ""
	}
	return tfAttr.(string)
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

func SetToStringArray(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	teams := make([]types.String, 0, len(set.Elements()))
	diags := diag.Diagnostics{}
	diags.Append(set.ElementsAs(ctx, &teams, true)...)
	if diags.HasError() {
		return nil, diags
	}
	convertedSet := make([]string, 0)
	for _, t := range teams {
		convertedSet = append(convertedSet, t.ValueString())
	}
	return convertedSet, diags
}

func ConvertStringMapToAttrStringMap(strMap map[string]string) map[string]attr.Value {
	attrMap := make(map[string]attr.Value, len(strMap))
	for key, val := range strMap {
		attrMap[key] = types.StringValue(val)
	}
	return attrMap
}

func ConvertAttrStringMapToStringMap(attrMap map[string]attr.Value) map[string]string {
	nativeMap := make(map[string]string, len(attrMap))
	for key, val := range attrMap {
		if val.IsNull() {
			nativeMap[key] = ""
		} else {
			nativeMap[key] = val.(types.String).ValueString()
		}
	}
	return nativeMap
}

func FlattenStringList(list []string) types.List {
	if list == nil {
		return types.ListValueMust(types.StringType, make([]attr.Value, 0))
	}

	elements := make([]attr.Value, 0, len(list))
	for _, s := range list {
		elements = append(elements, types.StringValue(s))
	}
	return types.ListValueMust(types.StringType, elements)
}

func Ternary[T interface{}](condition bool, whenTrue T, whenFalse T) T {
	if condition {
		return whenTrue
	}
	return whenFalse
}

func ToValueSlice(slice []string) []attr.Value {
	values := make([]attr.Value, len(slice))
	for i, s := range slice {
		values[i] = types.StringValue(s)
	}
	return values
}

func StringOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

func Map[T, V any](items []T, fn func(T) V) []V {
	result := make([]V, len(items))
	for i, t := range items {
		result[i] = fn(t)
	}
	return result
}

func GetIds(ids types.List) []string {
	var result = make([]string, 0, len(ids.Elements()))
	for _, id := range ids.Elements() {
		strVal, ok := id.(types.String)

		if !ok || strVal.IsNull() || strVal.IsUnknown() {
			continue
		}
		result = append(result, strVal.ValueString())
	}
	return result
}

func GetNumber(val types.Int64) int {
	v := 0
	if !val.IsNull() {
		v = int(val.ValueInt64())
	}

	return v
}

func ConvertMapStringToMapAttrValue(m map[string]string) map[string]attr.Value {
	result := make(map[string]attr.Value, len(m))
	for k, v := range m {
		result[k] = types.StringValue(v)
	}
	return result
}

func ConvertMapStringArrayToMapAttrValue(ctx context.Context, m map[string][]string) (map[string]attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := make(map[string]attr.Value, len(m))
	for k, v := range m {
		values := make([]attr.Value, len(v))
		for i, s := range v {
			values[i] = types.StringValue(s)
		}
		result[k], diags = types.SetValueFrom(ctx, types.StringType, v)
	}

	return result, diags
}

const sep = ":"

func BuildCompositeId(keys ...string) string {
	return strings.Join(keys, sep)
}
func SplitCompositeId(id string) []string {
	return strings.Split(id, sep)
}
