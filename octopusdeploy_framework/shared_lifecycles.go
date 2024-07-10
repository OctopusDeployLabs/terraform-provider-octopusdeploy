package octopusdeploy_framework

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FlattenRetentionPeriod(retentionPeriod *core.RetentionPeriod) types.List {
	if retentionPeriod == nil {
		return types.ListNull(types.ObjectType{AttrTypes: GetRetentionPeriodAttrTypes()})
	}
	return types.ListValueMust(
		types.ObjectType{AttrTypes: GetRetentionPeriodAttrTypes()},
		[]attr.Value{
			types.ObjectValueMust(
				GetRetentionPeriodAttrTypes(),
				map[string]attr.Value{
					"quantity_to_keep":    types.Int64Value(int64(retentionPeriod.QuantityToKeep)),
					"should_keep_forever": types.BoolValue(retentionPeriod.ShouldKeepForever),
					"unit":                types.StringValue(retentionPeriod.Unit),
				},
			),
		},
	)
}

func ExpandRetentionPeriod(v types.List) *core.RetentionPeriod {
	if v.IsNull() || v.IsUnknown() || len(v.Elements()) == 0 {
		return nil
	}

	obj := v.Elements()[0].(types.Object)
	attrs := obj.Attributes()

	var quantityToKeep int32
	if qty, ok := attrs["quantity_to_keep"].(types.Int64); ok && !qty.IsNull() {
		quantityToKeep = int32(qty.ValueInt64())
	}

	var shouldKeepForever bool
	if keep, ok := attrs["should_keep_forever"].(types.Bool); ok && !keep.IsNull() {
		shouldKeepForever = keep.ValueBool()
	}

	var unit string
	if u, ok := attrs["unit"].(types.String); ok && !u.IsNull() {
		unit = u.ValueString()
	}

	return core.NewRetentionPeriod(quantityToKeep, unit, shouldKeepForever)
}

func GetRetentionPeriodAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"quantity_to_keep":    types.Int64Type,
		"should_keep_forever": types.BoolType,
		"unit":                types.StringType,
	}
}

func FlattenPhases(phases []*lifecycles.Phase) types.List {
	if phases == nil {
		return types.ListNull(types.ObjectType{AttrTypes: GetPhaseAttrTypes()})
	}
	phasesList := make([]attr.Value, 0, len(phases))
	for _, phase := range phases {
		phasesList = append(phasesList, types.ObjectValueMust(
			GetPhaseAttrTypes(),
			map[string]attr.Value{
				"id":                                    types.StringValue(phase.ID),
				"name":                                  types.StringValue(phase.Name),
				"automatic_deployment_targets":          FlattenStringList(phase.AutomaticDeploymentTargets),
				"optional_deployment_targets":           FlattenStringList(phase.OptionalDeploymentTargets),
				"minimum_environments_before_promotion": types.Int64Value(int64(phase.MinimumEnvironmentsBeforePromotion)),
				"is_optional_phase":                     types.BoolValue(phase.IsOptionalPhase),
				"release_retention_policy":              FlattenRetentionPeriod(phase.ReleaseRetentionPolicy),
				"tentacle_retention_policy":             FlattenRetentionPeriod(phase.TentacleRetentionPolicy),
			},
		))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: GetPhaseAttrTypes()}, phasesList)
}

func ExpandPhases(phases types.List) []*lifecycles.Phase {
	if phases.IsNull() || phases.IsUnknown() {
		return nil
	}

	result := make([]*lifecycles.Phase, 0, len(phases.Elements()))

	for _, phaseElem := range phases.Elements() {
		phaseObj := phaseElem.(types.Object)
		phaseAttrs := phaseObj.Attributes()

		phase := &lifecycles.Phase{}

		if v, ok := phaseAttrs["id"].(types.String); ok && !v.IsNull() {
			phase.ID = v.ValueString()
		}

		if v, ok := phaseAttrs["name"].(types.String); ok && !v.IsNull() {
			phase.Name = v.ValueString()
		}

		if v, ok := phaseAttrs["automatic_deployment_targets"].(types.List); ok && !v.IsNull() {
			phase.AutomaticDeploymentTargets = ExpandStringList(v)
		}

		if v, ok := phaseAttrs["optional_deployment_targets"].(types.List); ok && !v.IsNull() {
			phase.OptionalDeploymentTargets = ExpandStringList(v)
		}

		if v, ok := phaseAttrs["minimum_environments_before_promotion"].(types.Int64); ok && !v.IsNull() {
			phase.MinimumEnvironmentsBeforePromotion = int32(v.ValueInt64())
		}

		if v, ok := phaseAttrs["is_optional_phase"].(types.Bool); ok && !v.IsNull() {
			phase.IsOptionalPhase = v.ValueBool()
		}

		if v, ok := phaseAttrs["release_retention_policy"].(types.List); ok && !v.IsNull() {
			phase.ReleaseRetentionPolicy = ExpandRetentionPeriod(v)
		}

		if v, ok := phaseAttrs["tentacle_retention_policy"].(types.List); ok && !v.IsNull() {
			phase.TentacleRetentionPolicy = ExpandRetentionPeriod(v)
		}

		result = append(result, phase)
	}

	return result
}

func GetPhaseAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                                    types.StringType,
		"name":                                  types.StringType,
		"automatic_deployment_targets":          types.ListType{ElemType: types.StringType},
		"optional_deployment_targets":           types.ListType{ElemType: types.StringType},
		"minimum_environments_before_promotion": types.Int64Type,
		"is_optional_phase":                     types.BoolType,
		"release_retention_policy":              types.ObjectType{AttrTypes: GetRetentionPeriodAttrTypes()},
		"tentacle_retention_policy":             types.ObjectType{AttrTypes: GetRetentionPeriodAttrTypes()},
	}
}

func LifecycleObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                        types.StringType,
		"space_id":                  types.StringType,
		"name":                      types.StringType,
		"description":               types.StringType,
		"phase":                     types.ListType{ElemType: types.ObjectType{AttrTypes: GetPhaseAttrTypes()}},
		"release_retention_policy":  types.ListType{ElemType: types.ObjectType{AttrTypes: GetRetentionPeriodAttrTypes()}},
		"tentacle_retention_policy": types.ListType{ElemType: types.ObjectType{AttrTypes: GetRetentionPeriodAttrTypes()}},
	}
}
