package schemas

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type reservedExecutionPropertiesValidator struct{}

func (v reservedExecutionPropertiesValidator) Description(ctx context.Context) string {
	return "execution properties must not contain reserved properties"
}

func (v reservedExecutionPropertiesValidator) MarkdownDescription(ctx context.Context) string {
	return "execution properties must not contain automatically generated properties or properties maintained by other attributes"
}

func (v reservedExecutionPropertiesValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	properties := make(map[string]types.String, len(req.ConfigValue.Elements()))
	diags := req.ConfigValue.ElementsAs(ctx, &properties, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	for key, _ := range properties {
		if warning, reserved := reservedExecutionProperties[key]; reserved {
			detail := fmt.Sprintf("This property managed automatically by the Octopus Deploy provider.\nAvoid configuring it manually as this can lead to discrepancies between planned and applied configurations, potentially causing unexpected behavior or drift.\n%s", warning.Suggestion)
			resp.Diagnostics.AddAttributeWarning(
				req.Path.AtMapKey(key),
				fmt.Sprintf("Process step execution property %q is reserved", key),
				detail,
			)
		}
	}
}

func warnAboutReservedExecutionProperties() validator.Map {
	return &reservedExecutionPropertiesValidator{}
}

type reservedExecutionPropertyWarning struct {
	Suggestion string
}

// Properties which managed by other attributes
var reservedExecutionProperties = map[string]reservedExecutionPropertyWarning{
	"Octopus.Action.Package.FeedId":             {Suggestion: "Consider to use 'primary_package' attribute to configure it's value"},
	"Octopus.Action.Package.PackageId":          {Suggestion: "Consider to use 'primary_package' attribute to configure it's value"},
	"Octopus.Action.Package.DownloadOnTentacle": {Suggestion: "Consider to use 'primary_package' attribute to configure it's value"},
}

func IsReservedExecutionProperty(name string) bool {
	_, exists := reservedExecutionProperties[name]
	return exists
}
