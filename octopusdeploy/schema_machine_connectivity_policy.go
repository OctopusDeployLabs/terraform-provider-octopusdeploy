package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machinepolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandMachineConnectivityPolicy(values interface{}) *machinepolicies.MachineConnectivityPolicy {
	if values == nil {
		return nil
	}
	flattenedValues := values.(*schema.Set)
	if len(flattenedValues.List()) == 0 {
		return nil
	}

	flattenedMap := flattenedValues.List()[0].(map[string]interface{})

	machineConnectivityPolicy := machinepolicies.NewMachineConnectivityPolicy()

	if v, ok := flattenedMap["machine_connectivity_behavior"]; ok {
		machineConnectivityPolicy.MachineConnectivityBehavior = v.(string)
	}

	return machineConnectivityPolicy
}

func flattenMachineConnectivityPolicy(machineConnectivityPolicy *machinepolicies.MachineConnectivityPolicy) []interface{} {
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
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"ExpectedToBeOnline",
				"MayBeOfflineAndCanBeSkipped",
			}, false)),
		},
	}
}
