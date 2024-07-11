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
				"id": resourceSchema.StringAttribute{
					Computed: true,
				},
				"name": resourceSchema.StringAttribute{
					Required: true,
				},
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
	return datasourceSchema.Schema{
		Attributes: map[string]datasourceSchema.Attribute{
			"id":           datasourceSchema.StringAttribute{Computed: true},
			"space_id":     datasourceSchema.StringAttribute{Optional: true},
			"ids":          datasourceSchema.ListAttribute{ElementType: types.StringType, Optional: true},
			"partial_name": datasourceSchema.StringAttribute{Optional: true},
			"skip":         datasourceSchema.Int64Attribute{Optional: true},
			"take":         datasourceSchema.Int64Attribute{Optional: true},
			"lifecycles": datasourceSchema.ListNestedAttribute{
				Computed: true,
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"id":                        datasourceSchema.StringAttribute{Computed: true},
						"space_id":                  datasourceSchema.StringAttribute{Computed: true},
						"name":                      datasourceSchema.StringAttribute{Computed: true},
						"description":               datasourceSchema.StringAttribute{Computed: true},
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
				"id":                                    datasourceSchema.StringAttribute{Computed: true},
				"name":                                  datasourceSchema.StringAttribute{Computed: true},
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
