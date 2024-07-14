package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetResourceLifecycleSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id":          util.GetIdResourceSchema(),
			"space_id":    util.GetSpaceIdResourceSchema("lifecycle"),
			"name":        util.GetNameResourceSchema(true),
			"description": util.GetDescriptionResourceSchema("lifecycle"),
		},
		Blocks: map[string]resourceSchema.Block{
			"phase":                     getResourcePhaseBlockSchema(),
			"release_retention_policy":  getResourceRetentionPolicyBlockSchema(),
			"tentacle_retention_policy": getResourceRetentionPolicyBlockSchema(),
		},
	}
}

func getResourcePhaseBlockSchema() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"id":   util.GetIdResourceSchema(),
				"name": util.GetNameResourceSchema(true),
				"automatic_deployment_targets": resourceSchema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},
				"optional_deployment_targets": resourceSchema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},
				"minimum_environments_before_promotion": resourceSchema.Int64Attribute{
					Optional: true,
					Computed: true,
					Default:  int64default.StaticInt64(0),
				},
				"is_optional_phase": resourceSchema.BoolAttribute{
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(false),
				},
			},
			Blocks: map[string]resourceSchema.Block{
				"release_retention_policy":  getResourceRetentionPolicyBlockSchema(),
				"tentacle_retention_policy": getResourceRetentionPolicyBlockSchema(),
			},
		},
	}
}

func getResourceRetentionPolicyBlockSchema() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"quantity_to_keep": resourceSchema.Int64Attribute{
					Optional: true,
				},
				"should_keep_forever": resourceSchema.BoolAttribute{
					Optional: true,
				},
				"unit": resourceSchema.StringAttribute{
					Optional: true,
				},
			},
		},
	}
}

func GetDatasourceLifecycleSchema() datasourceSchema.Schema {
	description := "lifecycle"
	return datasourceSchema.Schema{
		Attributes: map[string]datasourceSchema.Attribute{
			"id":           util.GetIdDatasourceSchema(),
			"space_id":     util.GetSpaceIdDatasourceSchema(description),
			"ids":          util.GetSpaceIdDatasourceSchema(description),
			"partial_name": util.GetQueryPartialNameDatasourceSchema(),
			"skip":         util.GetQuerySkipDatasourceSchema(),
			"take":         util.GetQueryTakeDatasourceSchema(),
			"lifecycles": datasourceSchema.ListNestedAttribute{
				Computed: true,
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"id":                        util.GetIdDatasourceSchema(),
						"space_id":                  util.GetSpaceIdDatasourceSchema(description),
						"name":                      util.GetNameDatasourceSchema(true),
						"description":               util.GetDescriptionDatasourceSchema(description),
						"phase":                     getDatasourcePhasesSchema(),
						"release_retention_policy":  getDatasourceRetentionPolicySchema(),
						"tentacle_retention_policy": getDatasourceRetentionPolicySchema(),
					},
				},
			},
		},
	}
}

func getDatasourcePhasesSchema() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"id":                                    util.GetIdDatasourceSchema(),
				"name":                                  util.GetNameDatasourceSchema(true),
				"automatic_deployment_targets":          datasourceSchema.ListAttribute{ElementType: types.StringType, Computed: true},
				"optional_deployment_targets":           datasourceSchema.ListAttribute{ElementType: types.StringType, Computed: true},
				"minimum_environments_before_promotion": datasourceSchema.Int64Attribute{Computed: true},
				"is_optional_phase":                     datasourceSchema.BoolAttribute{Computed: true},
				"release_retention_policy":              getDatasourceRetentionPolicySchema(),
				"tentacle_retention_policy":             getDatasourceRetentionPolicySchema(),
			},
		},
	}
}

func getDatasourceRetentionPolicySchema() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		Computed: true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: map[string]datasourceSchema.Attribute{
				"quantity_to_keep":    datasourceSchema.Int64Attribute{Computed: true},
				"should_keep_forever": datasourceSchema.BoolAttribute{Computed: true},
				"unit":                datasourceSchema.StringAttribute{Computed: true},
			},
		},
	}
}
