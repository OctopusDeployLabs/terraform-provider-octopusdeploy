package octopusdeploy

import (
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machinepolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandMachineCleanupPolicy(values interface{}) *machinepolicies.MachineCleanupPolicy {
	if values == nil {
		return nil
	}
	flattenedValues := values.(*schema.Set)
	if len(flattenedValues.List()) == 0 {
		return nil
	}

	flattenedMap := flattenedValues.List()[0].(map[string]interface{})

	machineCleanupPolicy := machinepolicies.NewMachineCleanupPolicy()

	if v, ok := flattenedMap["delete_machines_behavior"]; ok {
		machineCleanupPolicy.DeleteMachinesBehavior = v.(string)
	}

	if v, ok := flattenedMap["delete_machines_elapsed_timespan"]; ok {
		machineCleanupPolicy.DeleteMachinesElapsedTimeSpan = time.Duration(v.(int))
	}

	return machineCleanupPolicy
}

func flattenMachineCleanupPolicy(machineCleanupPolicy *machinepolicies.MachineCleanupPolicy) []interface{} {
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
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"DeleteUnavailableMachines",
				"DoNotDelete",
			}, false)),
		},
		"delete_machines_elapsed_timespan": {
			Computed:    true,
			Optional:    true,
			Type:        schema.TypeInt,
			Description: "In nanoseconds.",
		},
	}
}
