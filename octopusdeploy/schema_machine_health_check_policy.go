package octopusdeploy

import (
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machinepolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandMachineHealthCheckPolicy(values interface{}) *machinepolicies.MachineHealthCheckPolicy {
	if values == nil {
		return nil
	}
	flattenedValues := values.(*schema.Set)
	if len(flattenedValues.List()) == 0 {
		return nil
	}

	flattenedMap := flattenedValues.List()[0].(map[string]interface{})

	machineHealthCheckPolicy := machinepolicies.NewMachineHealthCheckPolicy()

	if v, ok := flattenedMap["bash_health_check_policy"]; ok {
		if len(v.([]interface{})) > 0 {
			machineHealthCheckPolicy.BashHealthCheckPolicy = expandMachineScriptPolicy(v)
		}
	}

	if v, ok := flattenedMap["health_check_cron"]; ok {
		machineHealthCheckPolicy.HealthCheckCron = v.(string)
	}

	if v, ok := flattenedMap["health_check_cron_timezone"]; ok {
		if s := v.(string); len(s) > 0 {
			machineHealthCheckPolicy.HealthCheckCronTimezone = s
		}
	}

	if v, ok := flattenedMap["health_check_interval"]; ok {
		machineHealthCheckPolicy.HealthCheckInterval = time.Duration(v.(int))
	}

	if v, ok := flattenedMap["health_check_type"]; ok {
		machineHealthCheckPolicy.HealthCheckType = v.(string)
	}

	if v, ok := flattenedMap["powershell_health_check_policy"]; ok {
		if len(v.([]interface{})) > 0 {
			machineHealthCheckPolicy.PowerShellHealthCheckPolicy = expandMachineScriptPolicy(v)
		}
	}

	return machineHealthCheckPolicy
}

func flattenMachineHealthCheckPolicy(machineHealthCheckPolicy *machinepolicies.MachineHealthCheckPolicy) []interface{} {
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
			Required: true,
			Type:     schema.TypeList,
		},
		"health_check_cron": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"health_check_cron_timezone": {
			Default:  "UTC",
			Optional: true,
			Type:     schema.TypeString,
		},
		"health_check_interval": {
			Default:     24 * time.Hour,
			Optional:    true,
			Type:        schema.TypeInt,
			Description: "In nanoseconds.",
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
			Required: true,
			Type:     schema.TypeList,
		},
	}
}
