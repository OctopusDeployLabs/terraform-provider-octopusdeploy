package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type lifecycleTypeResource struct {
	*Config
}

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

func (r *lifecycleTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "octopusdeploy_lifecycle"
}

func (r *lifecycleTypeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	description := "lifecycle"
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          util.GetIdResourceSchema(),
			"space_id":    util.GetSpaceIdResourceSchema(description),
			"name":        util.GetNameResourceSchema(true),
			"description": util.GetDescriptionResourceSchema(description),
		},
		Blocks: map[string]schema.Block{
			"phase": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id":                                    schema.StringAttribute{Computed: true},
						"name":                                  schema.StringAttribute{Required: true},
						"automatic_deployment_targets":          schema.ListAttribute{ElementType: types.StringType, Optional: true},
						"optional_deployment_targets":           schema.ListAttribute{ElementType: types.StringType, Optional: true},
						"minimum_environments_before_promotion": schema.Int64Attribute{Optional: true},
						"is_optional_phase":                     schema.BoolAttribute{Optional: true},
					},
					Blocks: map[string]schema.Block{

						"release_retention_policy":  getResourceRetentionPolicySchema(),
						"tentacle_retention_policy": getResourceRetentionPolicySchema(),
					},
				},
			},
			"release_retention_policy":  getResourceRetentionPolicySchema(),
			"tentacle_retention_policy": getResourceRetentionPolicySchema(),
		},
	}
}

func getResourceRetentionPolicySchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"quantity_to_keep": schema.Int64Attribute{
					Optional: true,
				},
				"should_keep_forever": schema.BoolAttribute{
					Optional: true,
				},
				"unit": schema.StringAttribute{
					Optional: true,
				},
			},
		},
	}
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

	newLifecycle := expandLifecycle(data)

	lifecycle, err := lifecycles.Add(r.Config.Client, newLifecycle)
	if err != nil {
		resp.Diagnostics.AddError("unable to create lifecycle", err.Error())
		return
	}
	data = flattenLifecycleResource(lifecycle)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *lifecycleTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *lifecycleTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lifecycle, err := lifecycles.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load lifecycle", err.Error())
		return
	}

	data = flattenLifecycleResource(lifecycle)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *lifecycleTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *lifecycleTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("updating lifecycle '%s'", data.ID.ValueString()))

	lifecycle := expandLifecycle(data)
	lifecycle.ID = state.ID.ValueString()

	updatedLifecycle, err := lifecycles.Update(r.Config.Client, lifecycle)
	if err != nil {
		resp.Diagnostics.AddError("unable to update lifecycle", err.Error())
		return
	}

	data = flattenLifecycleResource(updatedLifecycle)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

func expandLifecycle(data *lifecycleTypeResourceModel) *lifecycles.Lifecycle {
	lifecycle := lifecycles.NewLifecycle(data.Name.ValueString())
	lifecycle.Description = data.Description.ValueString()
	lifecycle.SpaceID = data.SpaceID.ValueString()

	lifecycle.Phases = expandPhases(data.Phase)
	lifecycle.ReleaseRetentionPolicy = expandRetentionPeriod(data.ReleaseRetentionPolicy)
	lifecycle.TentacleRetentionPolicy = expandRetentionPeriod(data.TentacleRetentionPolicy)

	return lifecycle
}

func expandPhases(phases types.List) []*lifecycles.Phase {
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
			phase.AutomaticDeploymentTargets = expandStringList(v)
		}

		if v, ok := phaseAttrs["optional_deployment_targets"].(types.List); ok && !v.IsNull() {
			phase.OptionalDeploymentTargets = expandStringList(v)
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

func expandStringList(list types.List) []string {
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

func flattenStringList(list []string) types.List {
	elements := make([]attr.Value, 0, len(list))
	for _, s := range list {
		elements = append(elements, types.StringValue(s))
	}
	return types.ListValueMust(types.StringType, elements)
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

func getRetentionPeriodAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"quantity_to_keep":    types.Int64Type,
		"should_keep_forever": types.BoolType,
		"unit":                types.StringType,
	}
}

func flattenPhases(phases []*lifecycles.Phase) types.List {
	if phases == nil {
		return types.ListNull(types.ObjectType{AttrTypes: getPhaseAttrTypes()})
	}
	phasesList := make([]attr.Value, 0, len(phases))
	for _, phase := range phases {
		phasesList = append(phasesList, types.ObjectValueMust(
			getPhaseAttrTypes(),
			map[string]attr.Value{
				"id":                                    types.StringValue(phase.ID),
				"name":                                  types.StringValue(phase.Name),
				"automatic_deployment_targets":          flattenStringList(phase.AutomaticDeploymentTargets),
				"optional_deployment_targets":           flattenStringList(phase.OptionalDeploymentTargets),
				"minimum_environments_before_promotion": types.Int64Value(int64(phase.MinimumEnvironmentsBeforePromotion)),
				"is_optional_phase":                     types.BoolValue(phase.IsOptionalPhase),
				"release_retention_policy":              flattenRetentionPeriod(phase.ReleaseRetentionPolicy),
				"tentacle_retention_policy":             flattenRetentionPeriod(phase.TentacleRetentionPolicy),
			},
		))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: getPhaseAttrTypes()}, phasesList)
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
