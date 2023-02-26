package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandRetentionPeriod(retentionPeriod map[string]interface{}) *core.RetentionPeriod {
	retention := core.RetentionPeriod{
		QuantityToKeep:    int32(retentionPeriod["quantity_to_keep"].(int)),
		ShouldKeepForever: retentionPeriod["should_keep_forever"].(bool),
		Unit:              retentionPeriod["unit"].(string),
	}
	return &retention
}

func flattenRetentionPeriod(r *core.RetentionPeriod) []interface{} {
	retentionPeriod := make(map[string]interface{})
	retentionPeriod["quantity_to_keep"] = int(r.QuantityToKeep)
	retentionPeriod["should_keep_forever"] = r.ShouldKeepForever
	retentionPeriod["unit"] = r.Unit
	return []interface{}{retentionPeriod}
}

func getRetentionPeriodSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"quantity_to_keep": {
			Default:          30,
			Description:      "The number of days/releases to keep. The default value is `30`. If `0` then all are kept.",
			Optional:         true,
			Type:             schema.TypeInt,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		},
		"should_keep_forever": {
			Default:     false,
			Description: "Indicates if items should never be deleted. The default value is `false`.",
			Optional:    true,
			Type:        schema.TypeBool,
		},
		"unit": {
			Default:     "Days",
			Description: "The unit of quantity to keep. Valid units are `Days` or `Items`. The default value is `Days`.",
			Optional:    true,
			Type:        schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"Days",
				"Items",
			}, false)),
		},
	}
}
