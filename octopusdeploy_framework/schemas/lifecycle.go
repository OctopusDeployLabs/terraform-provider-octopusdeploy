package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetLifecycleSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":          util.GetIdResourceSchema(),
			"space_id":    util.GetSpaceIdResourceSchema("lifecycle"),
			"name":        util.GetNameResourceSchema(true),
			"description": util.GetDescriptionResourceSchema("lifecycle"),
		},
		Blocks: map[string]schema.Block{
			"phase":                     getPhaseBlockSchema(),
			"release_retention_policy":  getRetentionPolicyBlockSchema(),
			"tentacle_retention_policy": getRetentionPolicyBlockSchema(),
		},
	}
}

func getPhaseBlockSchema() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Computed: true,
				},
				"name": schema.StringAttribute{
					Required: true,
				},
				"automatic_deployment_targets": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},
				"optional_deployment_targets": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},
				"minimum_environments_before_promotion": schema.Int64Attribute{
					Optional: true,
					Computed: true,
					Default:  int64default.StaticInt64(0),
				},
				"is_optional_phase": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(false),
				},
			},
			Blocks: map[string]schema.Block{
				"release_retention_policy":  getRetentionPolicyBlockSchema(),
				"tentacle_retention_policy": getRetentionPolicyBlockSchema(),
			},
		},
	}
}

func getRetentionPolicyBlockSchema() schema.ListNestedBlock {
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
