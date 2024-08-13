package modifiers

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type packagePropertiesModifier struct{}

func PackagePropertiesModifier() planmodifier.Map {
	return packagePropertiesModifier{}
}

func (p packagePropertiesModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	if req.PlanValue.IsUnknown() {
		properties := make(map[string]attr.Value)
		properties["Extract"] = types.StringValue("true")
		req.PlanValue = types.MapValueMust(types.StringType, properties)
		return
	}

	current := req.PlanValue.Elements()
	if _, ok := current["Extract"]; !ok {
		current["Extract"] = types.StringValue("true")
		req.PlanValue = types.MapValueMust(types.StringType, current)
		return
	}
}

func (p packagePropertiesModifier) Description(ctx context.Context) string {
	return "Package properties modifier for a given deployment process"
}

func (p packagePropertiesModifier) MarkdownDescription(ctx context.Context) string {
	return p.Description(ctx)
}

var _ planmodifier.Map = packagePropertiesModifier{}
