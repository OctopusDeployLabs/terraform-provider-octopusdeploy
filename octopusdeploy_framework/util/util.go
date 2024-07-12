package util

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetProviderName() string {
	return "octopusdeploy"
}

func GetTypeName(name string) string {
	return fmt.Sprintf("%s_%s", GetProviderName(), name)
}

func GetStringOrEmpty(tfAttr interface{}) string {
	if tfAttr == nil {
		return ""
	}
	return tfAttr.(string)
}

func ToStringArray(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
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
