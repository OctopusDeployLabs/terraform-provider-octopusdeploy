package octopusdeploy

import (
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandMachineCleanupPolicy(values interface{}) *octopusdeploy.MachineCleanupPolicy {
	flattenedValues := values.([]interface{})
	flattenedMap := flattenedValues[0].(map[string]interface{})

	return &octopusdeploy.MachineCleanupPolicy{
		DeleteMachinesBehavior:        flattenedMap["delete_machines_behavior"].(string),
		DeleteMachinesElapsedTimeSpan: time.Duration(flattenedMap["delete_machines_elapsed_timespan"].(int)),
	}
}

func flattenMachineCleanupPolicy(machineCleanupPolicy *octopusdeploy.MachineCleanupPolicy) []interface{} {
	if machineCleanupPolicy == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"delete_machines_behavior":         machineCleanupPolicy.DeleteMachinesBehavior,
		"delete_machines_elapsed_timespan": machineCleanupPolicy.DeleteMachinesElapsedTimeSpan,
	}}
}

func getMachineCleanupPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"delete_machines_behavior": {
			Default:  "DoNotDelete",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validateValueFunc([]string{
				"DeleteUnavailableMachines",
				"DoNotDelete",
			}),
		},
		"delete_machines_elapsed_timespan": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeInt,
		},
	}
}
