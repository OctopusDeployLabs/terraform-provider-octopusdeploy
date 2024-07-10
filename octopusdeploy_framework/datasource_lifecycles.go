package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/lifecycles"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
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
		IDs:         util.GetStringSlice(data.IDs),
		PartialName: data.PartialName.ValueString(),
		Skip:        int(data.Skip.ValueInt64()),
		Take:        int(data.Take.ValueInt64()),
	}

	lifecyclesResult, err := lifecycles.Get(l.Config.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read lifecycles, got error: %s", err))
		return
	}

	data.Lifecycles = flattenLifecycles(lifecyclesResult.Items)
	data.ID = types.StringValue("Lifecycles " + time.Now().UTC().String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenLifecycles(items []*lifecycles.Lifecycle) types.List {
	lifecyclesList := make([]attr.Value, 0, len(items))
	for _, lifecycle := range items {
		lifecycleMap := map[string]attr.Value{
			"id":                        types.StringValue(lifecycle.ID),
			"space_id":                  types.StringValue(lifecycle.SpaceID),
			"name":                      types.StringValue(lifecycle.Name),
			"description":               types.StringValue(lifecycle.Description),
			"phase":                     flattenPhases(lifecycle.Phases),
			"release_retention_policy":  flattenRetentionPeriod(lifecycle.ReleaseRetentionPolicy),
			"tentacle_retention_policy": flattenRetentionPeriod(lifecycle.TentacleRetentionPolicy),
		}
		lifecyclesList = append(lifecyclesList, types.ObjectValueMust(lifecycleObjectType(), lifecycleMap))
	}
	return types.ListValueMust(types.ObjectType{AttrTypes: lifecycleObjectType()}, lifecyclesList)
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
