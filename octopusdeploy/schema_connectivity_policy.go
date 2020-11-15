package octopusdeploy

import "github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"

func expandConnectivityPolicy(connectivityPolicy []interface{}) *octopusdeploy.ConnectivityPolicy {
	connectivityPolicyMap := connectivityPolicy[0].(map[string]interface{})
	return &octopusdeploy.ConnectivityPolicy{
		AllowDeploymentsToNoTargets: connectivityPolicyMap["allow_deployments_to_no_targets"].(bool),
		ExcludeUnhealthyTargets:     connectivityPolicyMap["exclude_unhealthy_targets"].(bool),
		SkipMachineBehavior:         octopusdeploy.SkipMachineBehavior(connectivityPolicyMap["skip_machine_behavior"].(string)),
		TargetRoles:                 getSliceFromTerraformTypeList(connectivityPolicyMap["target_roles"]),
	}
}

func flattenConnectivityPolicy(connectivityPolicy *octopusdeploy.ConnectivityPolicy) []interface{} {
	if connectivityPolicy == nil {
		return nil
	}

	flattenedConnectivityPolicy := make(map[string]interface{})
	flattenedConnectivityPolicy["allow_deployments_to_no_targets"] = connectivityPolicy.AllowDeploymentsToNoTargets
	flattenedConnectivityPolicy["exclude_unhealthy_targets"] = connectivityPolicy.ExcludeUnhealthyTargets
	flattenedConnectivityPolicy["skip_machine_behavior"] = connectivityPolicy.SkipMachineBehavior
	flattenedConnectivityPolicy["target_roles"] = connectivityPolicy.TargetRoles
	return []interface{}{flattenedConnectivityPolicy}
}
