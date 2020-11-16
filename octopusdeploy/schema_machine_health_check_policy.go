package octopusdeploy

import (
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandMachineHealthCheckPolicy(d *schema.ResourceData) *octopusdeploy.MachineHealthCheckPolicy {
	machineHealthCheckPolicy := octopusdeploy.NewMachineHealthCheckPolicy()

	if v, ok := d.GetOk("bash_health_check_policy"); ok {
		machineHealthCheckPolicy.BashHealthCheckPolicy = expandMachineScriptPolicy(v)
	}

	if v, ok := d.GetOk("health_check_cron"); ok {
		machineHealthCheckPolicy.HealthCheckCron = v.(string)
	}

	if v, ok := d.GetOk("health_check_cron_timezone"); ok {
		machineHealthCheckPolicy.HealthCheckCronTimezone = v.(string)
	}

	if v, ok := d.GetOk("health_check_interval"); ok {
		duration, _ := time.ParseDuration(v.(string))
		machineHealthCheckPolicy.HealthCheckInterval = duration
	}

	if v, ok := d.GetOk("health_check_type"); ok {
		machineHealthCheckPolicy.HealthCheckType = v.(string)
	}

	if v, ok := d.GetOk("powershell_health_check_policy"); ok {
		machineHealthCheckPolicy.PowerShellHealthCheckPolicy = expandMachineScriptPolicy(v)
	}

	return machineHealthCheckPolicy
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
			ValidateDiagFunc: validateValueFunc([]string{
				"OnlyConnectivity",
				"RunScript",
			}),
		},
		"powershell_health_check_policy": {
			Elem:     &schema.Resource{Schema: getMachineScriptPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}
