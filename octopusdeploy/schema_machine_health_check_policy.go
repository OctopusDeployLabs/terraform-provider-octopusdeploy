package octopusdeploy

import (
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandMachineHealthCheckPolicy(values interface{}) *octopusdeploy.MachineHealthCheckPolicy {
	flattenedValues := values.([]interface{})
	flattenedMap := flattenedValues[0].(map[string]interface{})

	duration, _ := time.ParseDuration(flattenedMap["health_check_interval"].(string))

	return &octopusdeploy.MachineHealthCheckPolicy{
		BashHealthCheckPolicy:       expandMachineScriptPolicy(flattenedMap["bash_health_check_policy"]),
		HealthCheckCron:             flattenedMap["health_check_cron"].(string),
		HealthCheckCronTimezone:     flattenedMap["health_check_cron_timezone"].(string),
		HealthCheckInterval:         duration,
		HealthCheckType:             flattenedMap["health_check_type"].(string),
		PowerShellHealthCheckPolicy: expandMachineScriptPolicy(flattenedMap["powershell_health_check_policy"]),
	}
}

func flattenMachineHealthCheckPolicy(machineHealthCheckPolicy *octopusdeploy.MachineHealthCheckPolicy) []interface{} {
	if machineHealthCheckPolicy == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"bash_health_check_policy":       flattenMachineScriptPolicy(machineHealthCheckPolicy.BashHealthCheckPolicy),
		"health_check_cron":              machineHealthCheckPolicy.HealthCheckCron,
		"health_check_cron_timezone":     machineHealthCheckPolicy.HealthCheckCronTimezone,
		"health_check_interval":          machineHealthCheckPolicy.HealthCheckInterval,
		"health_check_type":              machineHealthCheckPolicy.HealthCheckType,
		"powershell_health_check_policy": flattenMachineScriptPolicy(machineHealthCheckPolicy.PowerShellHealthCheckPolicy),
	}}
}

func getMachineHealthCheckPolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bash_health_check_policy": {
			Elem:     &schema.Resource{Schema: getMachineScriptPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"health_check_cron": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"health_check_cron_timezone": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"health_check_interval": {
			Optional: true,
			Type:     schema.TypeInt,
		},
		"health_check_type": {
			Default:  "RunScript",
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"OnlyConnectivity",
				"RunScript",
			}, false)),
		},
		"powershell_health_check_policy": {
			Elem:     &schema.Resource{Schema: getMachineScriptPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}
