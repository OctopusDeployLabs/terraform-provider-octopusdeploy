package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandMachineConnectivityPolicy(values interface{}) *octopusdeploy.MachineConnectivityPolicy {
	flattenedValues := values.([]interface{})
	flattenedMap := flattenedValues[0].(map[string]interface{})

	return &octopusdeploy.MachineConnectivityPolicy{
		MachineConnectivityBehavior: flattenedMap["machine_connectivity_behavior"].(string),
	}
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
