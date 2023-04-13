package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandRunbookRetentionPolicy(flattenedRetentionPeriod interface{}) *runbooks.RunbookRetentionPeriod {
	if flattenedRetentionPeriod == nil {
		return nil
	}

	retentionPeriodProperties := flattenedRetentionPeriod.([]interface{})
	if len(retentionPeriodProperties) == 1 {
		retentionPeriodMap := retentionPeriodProperties[0].(map[string]interface{})
		return &runbooks.RunbookRetentionPeriod{
			QuantityToKeep:    int32(retentionPeriodMap["quantity_to_keep"].(int)),
			ShouldKeepForever: retentionPeriodMap["should_keep_forever"].(bool),
		}
	}

	return nil
}

func flattenRunbookRetentionPeriod(r *runbooks.RunbookRetentionPeriod) []interface{} {
	retentionPeriod := make(map[string]interface{})
	retentionPeriod["quantity_to_keep"] = int(r.QuantityToKeep)
	retentionPeriod["should_keep_forever"] = r.ShouldKeepForever
	return []interface{}{retentionPeriod}
}

func getRunbookRetentionPeriodSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"quantity_to_keep": {
			Default:          30,
			Description:      "The number of days/releases to keep. The default value is `100`. If `0` then all are kept.",
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
	}
}
