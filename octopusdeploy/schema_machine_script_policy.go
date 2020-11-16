package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandMachineScriptPolicy(values interface{}) *octopusdeploy.MachineScriptPolicy {
	flattenedValues := values.([]interface{})
	flattenedMap := flattenedValues[0].(map[string]interface{})

	return &octopusdeploy.MachineScriptPolicy{
		RunType:    flattenedMap["run_type"].(string),
		ScriptBody: flattenedMap["script_body"].(*string),
	}
}

func flattenMachineScriptPolicy(machineScriptPolicy *octopusdeploy.MachineScriptPolicy) []interface{} {
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
			Optional: true,
			Type:     schema.TypeString,
		},
		"script_body": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}
