package octopusdeploy

import (
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandMachineCleanupPolicy(d *schema.ResourceData) *octopusdeploy.MachineCleanupPolicy {
	machineCleanupPolicy := octopusdeploy.NewMachineCleanupPolicy()

	if v, ok := d.GetOk("delete_machines_behavior"); ok {
		machineCleanupPolicy.DeleteMachinesBehavior = v.(string)
	}

	if v, ok := d.GetOk("delete_machines_elapsed_timespan"); ok {
		machineCleanupPolicy.DeleteMachinesElapsedTimeSpan = v.(time.Duration)
	}

	return machineCleanupPolicy
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
			Optional: true,
			Type:     schema.TypeInt,
		},
	}
}
