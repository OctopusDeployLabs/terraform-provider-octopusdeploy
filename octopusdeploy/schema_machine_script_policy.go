package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machinepolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandMachineScriptPolicy(values interface{}) *machinepolicies.MachineScriptPolicy {
	if values == nil {
		return nil
	}
	flattenedValues := values.([]interface{})
	if len(flattenedValues) == 0 {
		return nil
	}

	flattenedMap := flattenedValues[0].(map[string]interface{})

	machineScriptPolicy := machinepolicies.NewMachineScriptPolicy()

	if v, ok := flattenedMap["run_type"]; ok {
		machineScriptPolicy.RunType = v.(string)
	}

	if v, ok := flattenedMap["script_body"]; ok {
		scriptBody := v.(string)
		machineScriptPolicy.ScriptBody = &scriptBody
	}

	return machineScriptPolicy
}

func flattenMachineScriptPolicy(machineScriptPolicy *machinepolicies.MachineScriptPolicy) []interface{} {
	if machineScriptPolicy == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"run_type":    machineScriptPolicy.RunType,
		"script_body": machineScriptPolicy.ScriptBody,
	}}
}

func getMachineScriptPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"run_type": {
			Default:  "InheritFromDefault",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"InheritFromDefault",
				"Inline",
				"OnlyConnectivity",
			}, false)),
		},
		"script_body": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
