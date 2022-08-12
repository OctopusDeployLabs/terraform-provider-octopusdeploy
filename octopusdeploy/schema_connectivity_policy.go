package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandConnectivityPolicy(connectivityPolicy []interface{}) *core.ConnectivityPolicy {
	connectivityPolicyMap := connectivityPolicy[0].(map[string]interface{})
	return &core.ConnectivityPolicy{
		AllowDeploymentsToNoTargets: connectivityPolicyMap["allow_deployments_to_no_targets"].(bool),
		ExcludeUnhealthyTargets:     connectivityPolicyMap["exclude_unhealthy_targets"].(bool),
		SkipMachineBehavior:         core.SkipMachineBehavior(connectivityPolicyMap["skip_machine_behavior"].(string)),
		TargetRoles:                 getSliceFromTerraformTypeList(connectivityPolicyMap["target_roles"]),
	}
}

func flattenConnectivityPolicy(connectivityPolicy *core.ConnectivityPolicy) []interface{} {
	if connectivityPolicy == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"allow_deployments_to_no_targets": connectivityPolicy.AllowDeploymentsToNoTargets,
		"exclude_unhealthy_targets":       connectivityPolicy.ExcludeUnhealthyTargets,
		"skip_machine_behavior":           connectivityPolicy.SkipMachineBehavior,
		"target_roles":                    connectivityPolicy.TargetRoles,
	}}
}

func getConnectivityPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_deployments_to_no_targets": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"exclude_unhealthy_targets": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"skip_machine_behavior": {
			Default:  "None",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"SkipUnavailableMachines",
				"None",
			}, false)),
		},
		"target_roles": {
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}
