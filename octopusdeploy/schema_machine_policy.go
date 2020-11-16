package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenMachinePolicy(machinePolicy *octopusdeploy.MachinePolicy) map[string]interface{} {
	if machinePolicy == nil {
		return nil
	}

	return map[string]interface{}{
		"connection_connect_timeout":      machinePolicy.ConnectionConnectTimeout,
		"connection_retry_count_limit":    machinePolicy.ConnectionRetryCountLimit,
		"connection_retry_sleep_interval": machinePolicy.ConnectionRetrySleepInterval,
		"connection_retry_time_limit":     machinePolicy.ConnectionRetryTimeLimit,
		"description":                     machinePolicy.Description,
		"id":                              machinePolicy.GetID(),
		"is_default":                      machinePolicy.IsDefault,
		"machine_cleanup_policy":          flattenMachineCleanupPolicy(machinePolicy.MachineCleanupPolicy),
		"machine_connectivity_policy":     flattenMachineConnectivityPolicy(machinePolicy.MachineConnectivityPolicy),
		"machine_health_check_policy":     flattenMachineHealthCheckPolicy(machinePolicy.MachineHealthCheckPolicy),
		"machine_update_policy":           flattenMachineUpdatePolicy(machinePolicy.MachineUpdatePolicy),
		"name":                            machinePolicy.Name,
		"polling_request_maximum_message_processing_timeout": machinePolicy.PollingRequestMaximumMessageProcessingTimeout,
		"polling_request_queue_timeout":                      machinePolicy.PollingRequestMaximumMessageProcessingTimeout,
		"space_id":                                           machinePolicy.SpaceID,
	}

}

func getMachinePolicyDataSchema() map[string]*schema.Schema {
	machinePolicySchema := getMachinePolicySchema()
	for _, field := range machinePolicySchema {
		field.Computed = true
		field.Default = nil
		field.MaxItems = 0
		field.MinItems = 0
		field.Optional = false
		field.Required = false
		field.ValidateDiagFunc = nil
		field.ValidateFunc = nil
	}

	return map[string]*schema.Schema{
		"ids": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"machine_policies": {
			Computed: true,
			Elem:     &schema.Resource{Schema: machinePolicySchema},
			Type:     schema.TypeList,
		},
		"partial_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"skip": {
			Default:  0,
			Type:     schema.TypeInt,
			Optional: true,
		},
		"take": {
			Default:  1,
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
}

func getMachinePolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_connect_timeout": {
			Optional: true,
			Type:     schema.TypeInt,
		},
		"connection_retry_count_limit": {
			Optional: true,
			Type:     schema.TypeInt,
		},
		"connection_retry_sleep_interval": {
			Optional: true,
			Type:     schema.TypeInt,
		},
		"connection_retry_time_limit": {
			Optional: true,
			Type:     schema.TypeInt,
		},
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"is_default": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"machine_cleanup_policy": {
			Elem:     &schema.Resource{Schema: getMachineCleanupPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"machine_connectivity_policy": {
			Elem:     &schema.Resource{Schema: getMachineConnectivityPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"machine_health_check_policy": {
			Elem:     &schema.Resource{Schema: getMachineHealthCheckPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"machine_update_policy": {
			Elem:     &schema.Resource{Schema: getMachineUpdatePolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"polling_request_maximum_message_processing_timeout": {
			Optional: true,
			Type:     schema.TypeInt,
		},
		"polling_request_queue_timeout": {
			Optional: true,
			Type:     schema.TypeInt,
		},
		"space_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
	}
}
