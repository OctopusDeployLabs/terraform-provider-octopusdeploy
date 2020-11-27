package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandMachinePolicy(d *schema.ResourceData) *octopusdeploy.MachinePolicy {
	name := d.Get("name").(string)

	machinePolicy := octopusdeploy.NewMachinePolicy(name)
	machinePolicy.ID = d.Id()

	if v, ok := d.GetOk("connection_connect_timeout"); ok {
		machinePolicy.ConnectionConnectTimeout = time.Duration(v.(int))
	}

	if v, ok := d.GetOk("connection_retry_count_limit"); ok {
		machinePolicy.ConnectionRetryCountLimit = int32(v.(int))
	}

	if v, ok := d.GetOk("connection_retry_sleep_interval"); ok {
		machinePolicy.ConnectionRetrySleepInterval = time.Duration(v.(int))
	}

	if v, ok := d.GetOk("connection_retry_time_limit"); ok {
		machinePolicy.ConnectionRetryTimeLimit = time.Duration(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		machinePolicy.Description = v.(string)
	}

	if v, ok := d.GetOk("is_default"); ok {
		machinePolicy.IsDefault = v.(bool)
	}

	if v, ok := d.GetOk("machine_cleanup_policy"); ok {
		machinePolicy.MachineCleanupPolicy = expandMachineCleanupPolicy(v)
	}

	if v, ok := d.GetOk("machine_connectivity_policy"); ok {
		machinePolicy.MachineConnectivityPolicy = expandMachineConnectivityPolicy(v)
	}

	if v, ok := d.GetOk("machine_health_check_policy"); ok {
		machinePolicy.MachineHealthCheckPolicy = expandMachineHealthCheckPolicy(v)
	}

	if v, ok := d.GetOk("machine_update_policy"); ok {
		machinePolicy.MachineUpdatePolicy = expandMachineUpdatePolicy(v)
	}

	if v, ok := d.GetOk("name"); ok {
		machinePolicy.Name = v.(string)
	}

	if v, ok := d.GetOk("polling_request_maximum_message_processing_timeout"); ok {
		machinePolicy.PollingRequestMaximumMessageProcessingTimeout = time.Duration(v.(int))
	}

	if v, ok := d.GetOk("polling_request_queue_timeout"); ok {
		machinePolicy.PollingRequestMaximumMessageProcessingTimeout = time.Duration(v.(int))
	}

	if v, ok := d.GetOk("space_id"); ok {
		machinePolicy.SpaceID = v.(string)
	}

	return machinePolicy

}

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
	}

	return map[string]*schema.Schema{
		"ids": {
			Description: "Query and/or search by a list of IDs",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"machine_policies": {
			Computed: true,
			Elem:     &schema.Resource{Schema: machinePolicySchema},
			Type:     schema.TypeList,
		},
		"partial_name": {
			Description: "Query and/or search by partial name",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"skip": {
			Default:     0,
			Description: "Indicates the number of items to skip in the response",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"take": {
			Default:     1,
			Description: "Indicates the number of items to take (or return) in the response",
			Type:        schema.TypeInt,
			Optional:    true,
		},
	}
}

func getMachinePolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_connect_timeout": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"connection_retry_count_limit": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"connection_retry_sleep_interval": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"connection_retry_time_limit": {
			Computed: true,
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
			Computed: true,
			Elem:     &schema.Resource{Schema: getMachineCleanupPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"machine_connectivity_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getMachineConnectivityPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"machine_health_check_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getMachineHealthCheckPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeSet,
		},
		"machine_update_policy": {
			Computed: true,
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
			Computed: true,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"polling_request_queue_timeout": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"space_id": {
			Computed: true,
			Type:     schema.TypeString,
		},
	}
}

func setMachinePolicy(ctx context.Context, d *schema.ResourceData, machinePolicy *octopusdeploy.MachinePolicy) {
	if d == nil || machinePolicy == nil {
		return
	}

	d.Set("connection_connect_timeout", machinePolicy.ConnectionConnectTimeout)
	d.Set("connection_retry_count_limit", machinePolicy.ConnectionRetryCountLimit)
	d.Set("connection_retry_sleep_interval", machinePolicy.ConnectionRetrySleepInterval)
	d.Set("connection_retry_time_limit", machinePolicy.ConnectionRetryTimeLimit)
	d.Set("description", machinePolicy.Description)
	d.Set("id", machinePolicy.GetID())
	d.Set("is_default", machinePolicy.IsDefault)
	d.Set("machine_cleanup_policy", flattenMachineCleanupPolicy(machinePolicy.MachineCleanupPolicy))
	d.Set("machine_connectivity_policy", flattenMachineConnectivityPolicy(machinePolicy.MachineConnectivityPolicy))
	d.Set("machine_health_check_policy", flattenMachineHealthCheckPolicy(machinePolicy.MachineHealthCheckPolicy))
	d.Set("machine_update_policy", flattenMachineUpdatePolicy(machinePolicy.MachineUpdatePolicy))
	d.Set("name", machinePolicy.Name)
	d.Set("polling_request_maximum_message_processing_timeout", machinePolicy.PollingRequestMaximumMessageProcessingTimeout)
	d.Set("polling_request_queue_timeout", machinePolicy.PollingRequestMaximumMessageProcessingTimeout)
	d.Set("space_id", machinePolicy.SpaceID)
}
