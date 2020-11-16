package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandMachineScriptPolicy(flattenedMachineScriptPolicy interface{}) *octopusdeploy.MachineScriptPolicy {
	flattenedMachineScriptPolicyMap := flattenedMachineScriptPolicy.([]interface{})
	for key, value := range flattenedMachineScriptPolicyMap {
		log.Println(key)
		log.Println(value)
	}

	machineScriptPolicy := octopusdeploy.NewMachineScriptPolicy()

	return machineScriptPolicy
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
