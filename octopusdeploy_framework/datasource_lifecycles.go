package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"time"
)

type lifecyclesDataSource struct {
	*Config
}

type lifecyclesDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	SpaceID     types.String `tfsdk:"space_id"`
	IDs         types.List   `tfsdk:"ids"`
	PartialName types.String `tfsdk:"partial_name"`
	Skip        types.Int64  `tfsdk:"skip"`
	Take        types.Int64  `tfsdk:"take"`
	Lifecycles  types.List   `tfsdk:"lifecycles"`
}

func NewLifecyclesDataSource() datasource.DataSource {
	return &lifecyclesDataSource{}
}

func (l *lifecyclesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	tflog.Debug(ctx, "lifecycles datasource Metadata")
	resp.TypeName = "octopusdeploy_lifecycles"
}

func (l *lifecyclesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	tflog.Debug(ctx, "lifecycles datasource Schema")
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true},
			"space_id":     schema.StringAttribute{Optional: true},
			"ids":          schema.ListAttribute{ElementType: types.StringType, Optional: true},
			"partial_name": schema.StringAttribute{Optional: true},
			"skip":         schema.Int64Attribute{Optional: true},
			"take":         schema.Int64Attribute{Optional: true},
			"lifecycles": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                        schema.StringAttribute{Computed: true},
						"space_id":                  schema.StringAttribute{Computed: true},
						"name":                      schema.StringAttribute{Computed: true},
						"description":               schema.StringAttribute{Computed: true},
						"phase":                     getPhasesSchema(),
						"release_retention_policy":  getRetentionPolicySchema(),
						"tentacle_retention_policy": getRetentionPolicySchema(),
					},
				},
			},
		},
	}
}

func (l *lifecyclesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Debug(ctx, "lifecycles datasource Configure")
	l.Config = DataSourceConfiguration(req, resp)
}

func (l *lifecyclesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "lifecycles datasource Read")
	var data lifecyclesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := lifecycles.Query{
		IDs:         getStringSlice(data.IDs),
		PartialName: data.PartialName.ValueString(),
		Skip:        int(data.Skip.ValueInt64()),
		Take:        int(data.Take.ValueInt64()),
	}

	lifecyclesResult, err := lifecycles.Get(l.Config.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read lifecycles, got error: %s", err))
		return
	}

	// Map the retrieved lifecycles to the data model
	lifecyclesList := make([]attr.Value, 0, len(lifecyclesResult.Items))
	for _, lifecycle := range lifecyclesResult.Items {
		lifecycleMap := map[string]attr.Value{
			"id":          types.StringValue(lifecycle.ID),
			"space_id":    types.StringValue(lifecycle.SpaceID),
			"name":        types.StringValue(lifecycle.Name),
			"description": types.StringValue(lifecycle.Description),
		}

		// Map phases
		phases := make([]attr.Value, 0, len(lifecycle.Phases))
		for _, phase := range lifecycle.Phases {
			phaseMap := map[string]attr.Value{
				"id":                                    types.StringValue(phase.ID),
				"name":                                  types.StringValue(phase.Name),
				"automatic_deployment_targets":          types.ListValueMust(types.StringType, toValueSlice(phase.AutomaticDeploymentTargets)),
				"optional_deployment_targets":           types.ListValueMust(types.StringType, toValueSlice(phase.OptionalDeploymentTargets)),
				"minimum_environments_before_promotion": types.Int64Value(int64(phase.MinimumEnvironmentsBeforePromotion)),
				"is_optional_phase":                     types.BoolValue(phase.IsOptionalPhase),
				"release_retention_policy":              mapRetentionPolicyList(phase.ReleaseRetentionPolicy),
				"tentacle_retention_policy":             mapRetentionPolicyList(phase.TentacleRetentionPolicy),
			}
			phases = append(phases, types.ObjectValueMust(phaseObjectType(), phaseMap))
		}
		lifecycleMap["phase"] = types.ListValueMust(types.ObjectType{AttrTypes: phaseObjectType()}, phases)

		// Map retention policies
		lifecycleMap["release_retention_policy"] = mapRetentionPolicyList(lifecycle.ReleaseRetentionPolicy)
		lifecycleMap["tentacle_retention_policy"] = mapRetentionPolicyList(lifecycle.TentacleRetentionPolicy)

		lifecyclesList = append(lifecyclesList, types.ObjectValueMust(lifecycleObjectType(), lifecycleMap))
	}

	data.Lifecycles = types.ListValueMust(types.ObjectType{AttrTypes: lifecycleObjectType()}, lifecyclesList)
	data.ID = types.StringValue("Lifecycles " + time.Now().UTC().String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getPhasesSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id":                                    schema.StringAttribute{Computed: true},
				"name":                                  schema.StringAttribute{Computed: true},
				"automatic_deployment_targets":          schema.ListAttribute{ElementType: types.StringType, Computed: true},
				"optional_deployment_targets":           schema.ListAttribute{ElementType: types.StringType, Computed: true},
				"minimum_environments_before_promotion": schema.Int64Attribute{Computed: true},
				"is_optional_phase":                     schema.BoolAttribute{Computed: true},
				"release_retention_policy":              getRetentionPolicySchema(),
				"tentacle_retention_policy":             getRetentionPolicySchema(),
			},
		},
	}
}

func getRetentionPolicySchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"quantity_to_keep":    schema.Int64Attribute{Computed: true},
				"should_keep_forever": schema.BoolAttribute{Computed: true},
				"unit":                schema.StringAttribute{Computed: true},
			},
		},
	}
}

func lifecycleObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                        types.StringType,
		"space_id":                  types.StringType,
		"name":                      types.StringType,
		"description":               types.StringType,
		"phase":                     types.ListType{ElemType: types.ObjectType{AttrTypes: phaseObjectType()}},
		"release_retention_policy":  types.ListType{ElemType: types.ObjectType{AttrTypes: retentionPolicyObjectType()}},
		"tentacle_retention_policy": types.ListType{ElemType: types.ObjectType{AttrTypes: retentionPolicyObjectType()}},
	}
}

func phaseObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                                    types.StringType,
		"name":                                  types.StringType,
		"automatic_deployment_targets":          types.ListType{ElemType: types.StringType},
		"optional_deployment_targets":           types.ListType{ElemType: types.StringType},
		"minimum_environments_before_promotion": types.Int64Type,
		"is_optional_phase":                     types.BoolType,
		"release_retention_policy":              types.ListType{ElemType: types.ObjectType{AttrTypes: retentionPolicyObjectType()}},
		"tentacle_retention_policy":             types.ListType{ElemType: types.ObjectType{AttrTypes: retentionPolicyObjectType()}},
	}
}

func retentionPolicyObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"quantity_to_keep":    types.Int64Type,
		"should_keep_forever": types.BoolType,
		"unit":                types.StringType,
	}
}

func toValueSlice(slice []string) []attr.Value {
	values := make([]attr.Value, len(slice))
	for i, s := range slice {
		values[i] = types.StringValue(s)
	}
	return values
}

func mapRetentionPolicyList(policy *core.RetentionPeriod) attr.Value {
	if policy == nil {
		return types.ListValueMust(types.ObjectType{AttrTypes: retentionPolicyObjectType()}, []attr.Value{})
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: retentionPolicyObjectType()}, []attr.Value{
		types.ObjectValueMust(retentionPolicyObjectType(), map[string]attr.Value{
			"quantity_to_keep":    types.Int64Value(int64(policy.QuantityToKeep)),
			"should_keep_forever": types.BoolValue(policy.ShouldKeepForever),
			"unit":                types.StringValue(policy.Unit),
		}),
	})
}

func getStringSlice(list types.List) []string {
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
