package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenMachineUpdatePolicy(machineUpdatePolicy *octopusdeploy.MachineUpdatePolicy) []interface{} {
	if machineUpdatePolicy == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"calamari_update_behavior":   machineUpdatePolicy.CalamariUpdateBehavior,
		"tentacle_update_account_id": machineUpdatePolicy.TentacleUpdateAccountID,
		"tentacle_update_behavior":   machineUpdatePolicy.TentacleUpdateBehavior,
	}}
}

func getMachineUpdatePolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"calamari_update_behavior": {
			Default:  "UpdateOnDeployment",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validateValueFunc([]string{
				"UpdateAlways",
				"UpdateOnDeployment",
				"UpdateOnNewMachine",
			}),
		},
		"tentacle_update_account_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"tentacle_update_behavior": {
			Default:  "NeverUpdate",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validateValueFunc([]string{
				"NeverUpdate",
				"Update",
			}),
		},
	}
}
