package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machinepolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandMachineUpdatePolicy(values interface{}) *machinepolicies.MachineUpdatePolicy {
	if values == nil {
		return nil
	}
	flattenedValues := values.(*schema.Set)
	if len(flattenedValues.List()) == 0 {
		return nil
	}

	flattenedMap := flattenedValues.List()[0].(map[string]interface{})

	machineUpdatePolicy := machinepolicies.NewMachineUpdatePolicy()

	if v, ok := flattenedMap["calamari_update_behavior"]; ok {
		machineUpdatePolicy.CalamariUpdateBehavior = v.(string)
	}

	if v, ok := flattenedMap["tentacle_update_account_id"]; ok {
		machineUpdatePolicy.TentacleUpdateAccountID = v.(string)
	}

	if v, ok := flattenedMap["tentacle_update_behavior"]; ok {
		machineUpdatePolicy.TentacleUpdateBehavior = v.(string)
	}

	if v, ok := flattenedMap["kubernetes_agent_update_behavior"]; ok {
		machineUpdatePolicy.KubernetesAgentUpdateBehavior = v.(string)
	}

	return machineUpdatePolicy
}

func flattenMachineUpdatePolicy(machineUpdatePolicy *machinepolicies.MachineUpdatePolicy) []interface{} {
	if machineUpdatePolicy == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"calamari_update_behavior":         machineUpdatePolicy.CalamariUpdateBehavior,
		"tentacle_update_account_id":       machineUpdatePolicy.TentacleUpdateAccountID,
		"tentacle_update_behavior":         machineUpdatePolicy.TentacleUpdateBehavior,
		"kubernetes_agent_update_behavior": machineUpdatePolicy.KubernetesAgentUpdateBehavior,
	}}
}

func getMachineUpdatePolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"calamari_update_behavior": {
			Default:  "UpdateOnDeployment",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"UpdateAlways",
				"UpdateOnDeployment",
				"UpdateOnNewMachine",
			}, false)),
			Description: "The behaviour of how Calamari is updated. Valid values are `UpdateAlways`, `UpdateOnDeployment` and `UpdateOnNewMachine`.",
		},
		"tentacle_update_account_id": {
			Optional:    true,
			Type:        schema.TypeString,
			Description: "The Account ID to perform any Tentacle updates under.",
		},
		"tentacle_update_behavior": {
			Default:  "NeverUpdate",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"NeverUpdate",
				"Update",
			}, false)),
			Description: "The behaviour of how Tentacle machines are updated. Valid values are `NeverUpdate` and `Update`.",
		},
		"kubernetes_agent_update_behavior": {
			Default:  "Update",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"NeverUpdate",
				"Update",
			}, false)),
			Description: "The behaviour of how Kubernetes agent machines are updated. Valid values are `NeverUpdate` and `Update`.",
		},
	}
}
