package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

type lifecycleTypeResource struct {
	*Config
}

var _ resource.Resource = &lifecycleTypeResource{}
var _ resource.ResourceWithImportState = &lifecycleTypeResource{}

type lifecycleTypeResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	SpaceID                 types.String `tfsdk:"space_id"`
	Name                    types.String `tfsdk:"name"`
	Description             types.String `tfsdk:"description"`
	Phase                   types.List   `tfsdk:"phase"`
	ReleaseRetentionPolicy  types.List   `tfsdk:"release_retention_policy"`
	TentacleRetentionPolicy types.List   `tfsdk:"tentacle_retention_policy"`
}

func NewLifecycleResource() resource.Resource {
	return &lifecycleTypeResource{}
}

func (r *lifecycleTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")
	lifecycleID := idParts[len(idParts)-1]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), lifecycleID)...)
	// Note: This implementation assumes that space_id is set at the provider level
}

func (r *lifecycleTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("lifecycle")
}

func (r *lifecycleTypeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetResourceLifecycleSchema()
}

func (r *lifecycleTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *lifecycleTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *lifecycleTypeResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	releaseRetentionPolicySet, tentacleRetentionPolicySet, defaultPolicy := setDefaultRetentionPolicy(data)

	newLifecycle := expandLifecycle(data)
	lifecycle, err := lifecycles.Add(r.Config.Client, newLifecycle)
	if err != nil {
		resp.Diagnostics.AddError("unable to create lifecycle", err.Error())
		return
	}
	data = flattenLifecycleResource(lifecycle)

	removeDefaultRetentionPolicy(releaseRetentionPolicySet, data, defaultPolicy, tentacleRetentionPolicySet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *lifecycleTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *lifecycleTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	releaseRetentionPolicySet, tentacleRetentionPolicySet, defaultPolicy := setDefaultRetentionPolicy(data)

	lifecycle, err := lifecycles.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load lifecycle", err.Error())
		return
	}
	data = flattenLifecycleResource(lifecycle)

	removeDefaultRetentionPolicy(releaseRetentionPolicySet, data, defaultPolicy, tentacleRetentionPolicySet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *lifecycleTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *lifecycleTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	releaseRetentionPolicySet, tentacleRetentionPolicySet, defaultPolicy := setDefaultRetentionPolicy(data)

	tflog.Debug(ctx, fmt.Sprintf("updating lifecycle '%s'", data.ID.ValueString()))

	lifecycle := expandLifecycle(data)
	lifecycle.ID = state.ID.ValueString()

	updatedLifecycle, err := lifecycles.Update(r.Config.Client, lifecycle)
	if err != nil {
		resp.Diagnostics.AddError("unable to update lifecycle", err.Error())
		return
	}

	data = flattenLifecycleResource(updatedLifecycle)

	removeDefaultRetentionPolicy(releaseRetentionPolicySet, data, defaultPolicy, tentacleRetentionPolicySet)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func removeDefaultRetentionPolicy(releaseRetentionPolicySet bool, data *lifecycleTypeResourceModel, defaultPolicy types.List, tentacleRetentionPolicySet bool) {
	// Remove default policies from data before setting state, but only if we added them
	if !releaseRetentionPolicySet && data.ReleaseRetentionPolicy.Equal(defaultPolicy) {
		data.ReleaseRetentionPolicy = types.ListNull(types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()})
	}
	if !tentacleRetentionPolicySet && data.TentacleRetentionPolicy.Equal(defaultPolicy) {
		data.TentacleRetentionPolicy = types.ListNull(types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()})
	}
}

func setDefaultRetentionPolicy(data *lifecycleTypeResourceModel) (bool, bool, types.List) {
	releaseRetentionPolicySet := !data.ReleaseRetentionPolicy.IsNull() && len(data.ReleaseRetentionPolicy.Elements()) > 0
	tentacleRetentionPolicySet := !data.TentacleRetentionPolicy.IsNull() && len(data.TentacleRetentionPolicy.Elements()) > 0

	// Set default policies only if they're not in the plan
	defaultPolicy := flattenRetentionPeriod(core.NewRetentionPeriod(30, "Days", false))
	if !releaseRetentionPolicySet {
		data.ReleaseRetentionPolicy = defaultPolicy
	}
	if !tentacleRetentionPolicySet {
		data.TentacleRetentionPolicy = defaultPolicy
	}
	return releaseRetentionPolicySet, tentacleRetentionPolicySet, defaultPolicy
}

func (r *lifecycleTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data lifecycleTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := lifecycles.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete lifecycle", err.Error())
		return
	}
}

func resourceConfiguration(req resource.ConfigureRequest, resp *resource.ConfigureResponse) *Config {
	if req.ProviderData == nil {
		return nil
	}

	p, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return p
}
func setDefaultRetentionPolicies(data *lifecycleTypeResourceModel) {
	defaultPolicy := flattenRetentionPeriod(core.NewRetentionPeriod(30, "Days", false))

	if data.ReleaseRetentionPolicy.IsNull() || len(data.ReleaseRetentionPolicy.Elements()) == 0 {
		data.ReleaseRetentionPolicy = defaultPolicy
	}

	if data.TentacleRetentionPolicy.IsNull() || len(data.TentacleRetentionPolicy.Elements()) == 0 {
		data.TentacleRetentionPolicy = defaultPolicy
	}
}

func flattenLifecycleResource(lifecycle *lifecycles.Lifecycle) *lifecycleTypeResourceModel {
	return &lifecycleTypeResourceModel{
		ID:                      types.StringValue(lifecycle.ID),
		SpaceID:                 types.StringValue(lifecycle.SpaceID),
		Name:                    types.StringValue(lifecycle.Name),
		Description:             types.StringValue(lifecycle.Description),
		Phase:                   flattenPhases(lifecycle.Phases),
		ReleaseRetentionPolicy:  flattenRetentionPeriod(lifecycle.ReleaseRetentionPolicy),
		TentacleRetentionPolicy: flattenRetentionPeriod(lifecycle.TentacleRetentionPolicy),
	}
}

func flattenPhases(phases []*lifecycles.Phase) types.List {
	if phases == nil {
		return types.ListNull(types.ObjectType{AttrTypes: getPhaseAttrTypes()})
	}
	phasesList := make([]attr.Value, 0, len(phases))

	for _, phase := range phases {
		attrs := map[string]attr.Value{
			"id":                                    types.StringValue(phase.ID),
			"name":                                  types.StringValue(phase.Name),
			"automatic_deployment_targets":          util.FlattenStringList(phase.AutomaticDeploymentTargets),
			"optional_deployment_targets":           util.FlattenStringList(phase.OptionalDeploymentTargets),
			"minimum_environments_before_promotion": types.Int64Value(int64(phase.MinimumEnvironmentsBeforePromotion)),
			"is_optional_phase":                     types.BoolValue(phase.IsOptionalPhase),
			"release_retention_policy":              util.Ternary(phase.ReleaseRetentionPolicy != nil, flattenRetentionPeriod(phase.ReleaseRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()})),
			"tentacle_retention_policy":             util.Ternary(phase.TentacleRetentionPolicy != nil, flattenRetentionPeriod(phase.TentacleRetentionPolicy), types.ListNull(types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()})),
		}
		phasesList = append(phasesList, types.ObjectValueMust(getPhaseAttrTypes(), attrs))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: getPhaseAttrTypes()}, phasesList)
}

func flattenRetentionPeriod(retentionPeriod *core.RetentionPeriod) types.List {
	if retentionPeriod == nil {
		return types.ListNull(types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()})
	}
	return types.ListValueMust(
		types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()},
		[]attr.Value{
			types.ObjectValueMust(
				getRetentionPeriodAttrTypes(),
				map[string]attr.Value{
					"quantity_to_keep":    types.Int64Value(int64(retentionPeriod.QuantityToKeep)),
					"should_keep_forever": types.BoolValue(retentionPeriod.ShouldKeepForever),
					"unit":                types.StringValue(retentionPeriod.Unit),
				},
			),
		},
	)
}

func expandLifecycle(data *lifecycleTypeResourceModel) *lifecycles.Lifecycle {
	if data == nil {
		return nil
	}

	lifecycle := lifecycles.NewLifecycle(data.Name.ValueString())
	lifecycle.Description = data.Description.ValueString()
	lifecycle.SpaceID = data.SpaceID.ValueString()
	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		lifecycle.ID = data.ID.ValueString()
	}

	lifecycle.Phases = expandPhases(data.Phase)
	lifecycle.ReleaseRetentionPolicy = expandRetentionPeriod(data.ReleaseRetentionPolicy)
	lifecycle.TentacleRetentionPolicy = expandRetentionPeriod(data.TentacleRetentionPolicy)

	return lifecycle
}

func expandPhases(phases types.List) []*lifecycles.Phase {
	if phases.IsNull() || phases.IsUnknown() || len(phases.Elements()) == 0 {
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
			phase.AutomaticDeploymentTargets = util.ExpandStringList(v)
		}

		if v, ok := phaseAttrs["optional_deployment_targets"].(types.List); ok && !v.IsNull() {
			phase.OptionalDeploymentTargets = util.ExpandStringList(v)
		}

		if v, ok := phaseAttrs["minimum_environments_before_promotion"].(types.Int64); ok && !v.IsNull() {
			phase.MinimumEnvironmentsBeforePromotion = int32(v.ValueInt64())
		}

		if v, ok := phaseAttrs["is_optional_phase"].(types.Bool); ok && !v.IsNull() {
			phase.IsOptionalPhase = v.ValueBool()
		}

		if v, ok := phaseAttrs["release_retention_policy"].(types.List); ok && !v.IsNull() {
			phase.ReleaseRetentionPolicy = expandRetentionPeriod(v)
		}

		if v, ok := phaseAttrs["tentacle_retention_policy"].(types.List); ok && !v.IsNull() {
			phase.TentacleRetentionPolicy = expandRetentionPeriod(v)
		}

		result = append(result, phase)
	}

	return result
}

func expandRetentionPeriod(v types.List) *core.RetentionPeriod {
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

func getRetentionPeriodAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"quantity_to_keep":    types.Int64Type,
		"should_keep_forever": types.BoolType,
		"unit":                types.StringType,
	}
}

func getPhaseAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                                    types.StringType,
		"name":                                  types.StringType,
		"automatic_deployment_targets":          types.ListType{ElemType: types.StringType},
		"optional_deployment_targets":           types.ListType{ElemType: types.StringType},
		"minimum_environments_before_promotion": types.Int64Type,
		"is_optional_phase":                     types.BoolType,
		"release_retention_policy":              types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()}},
		"tentacle_retention_policy":             types.ListType{ElemType: types.ObjectType{AttrTypes: getRetentionPeriodAttrTypes()}},
	}
}
