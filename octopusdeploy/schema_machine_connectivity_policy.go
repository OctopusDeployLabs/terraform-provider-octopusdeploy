package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandMachineConnectivityPolicy(d *schema.ResourceData) *octopusdeploy.MachineConnectivityPolicy {
	machineConnectivityPolicy := octopusdeploy.NewMachineConnectivityPolicy()

	if v, ok := d.GetOk("machine_connectivity_behavior"); ok {
		machineConnectivityPolicy.MachineConnectivityBehavior = v.(string)
	}

	return machineConnectivityPolicy
}

func flattenMachineConnectivityPolicy(machineConnectivityPolicy *octopusdeploy.MachineConnectivityPolicy) []interface{} {
	if machineConnectivityPolicy == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"machine_connectivity_behavior": machineConnectivityPolicy.MachineConnectivityBehavior,
	}}
}

func getMachineConnectivityPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"machine_connectivity_behavior": {
			Default:  "ExpectedToBeOnline",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validateValueFunc([]string{
				"ExpectedToBeOnline",
				"MayBeOfflineAndCanBeSkipped",
			}),
		},
	}
}
