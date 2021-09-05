package octopusdeploy

import (
	"context"
	"fmt"
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
		if len(v.(*schema.Set).List()) > 0 {
			machinePolicy.MachineCleanupPolicy = expandMachineCleanupPolicy(v)
		}
	}

	if v, ok := d.GetOk("machine_connectivity_policy"); ok {
		if len(v.(*schema.Set).List()) > 0 {
			machinePolicy.MachineConnectivityPolicy = expandMachineConnectivityPolicy(v)
		}
	}

	if v, ok := d.GetOk("machine_health_check_policy"); ok {
		if len(v.(*schema.Set).List()) > 0 {
			machinePolicy.MachineHealthCheckPolicy = expandMachineHealthCheckPolicy(v)
		}
	}

	if v, ok := d.GetOk("machine_update_policy"); ok {
		if len(v.(*schema.Set).List()) > 0 {
			machinePolicy.MachineUpdatePolicy = expandMachineUpdatePolicy(v)
		}
	}

	if v, ok := d.GetOk("name"); ok {
		machinePolicy.Name = v.(string)
	}

	if v, ok := d.GetOk("polling_request_maximum_message_processing_timeout"); ok {
		machinePolicy.PollingRequestMaximumMessageProcessingTimeout = time.Duration(v.(int))
	}

	if v, ok := d.GetOk("polling_request_queue_timeout"); ok {
		machinePolicy.PollingRequestQueueTimeout = time.Duration(v.(int))
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
		"polling_request_queue_timeout":                      machinePolicy.PollingRequestQueueTimeout,
		"space_id":                                           machinePolicy.SpaceID,
	}
}

func getMachinePolicyDataSchema() map[string]*schema.Schema {
	dataSchema := getMachinePolicySchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"ids": getQueryIDs(),
		"machine_policies": {
			Computed:    true,
			Description: "A list of machine policies that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
	}
}

func getMachinePolicySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"connection_connect_timeout": {
			Default:  time.Minute,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"connection_retry_count_limit": {
			Default:  5,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"connection_retry_sleep_interval": {
			Default:  time.Second,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"connection_retry_time_limit": {
			Default:  5 * time.Minute,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"description": getDescriptionSchema(),
		"id":          getIDSchema(),
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
		"name": getNameSchema(true),
		"polling_request_maximum_message_processing_timeout": {
			Default:  10 * time.Minute,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"polling_request_queue_timeout": {
			Default:  2 * time.Minute,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"space_id": getSpaceIDSchema(),
	}
}

func setMachinePolicy(ctx context.Context, d *schema.ResourceData, machinePolicy *octopusdeploy.MachinePolicy) error {
	d.Set("connection_connect_timeout", machinePolicy.ConnectionConnectTimeout)
	d.Set("connection_retry_count_limit", machinePolicy.ConnectionRetryCountLimit)
	d.Set("connection_retry_sleep_interval", machinePolicy.ConnectionRetrySleepInterval)
	d.Set("connection_retry_time_limit", machinePolicy.ConnectionRetryTimeLimit)
	d.Set("description", machinePolicy.Description)
	d.Set("id", machinePolicy.GetID())
	d.Set("is_default", machinePolicy.IsDefault)
	d.Set("name", machinePolicy.Name)
	d.Set("polling_request_maximum_message_processing_timeout", machinePolicy.PollingRequestMaximumMessageProcessingTimeout)
	d.Set("polling_request_queue_timeout", machinePolicy.PollingRequestQueueTimeout)
	d.Set("space_id", machinePolicy.SpaceID)

	if err := d.Set("machine_cleanup_policy", flattenMachineCleanupPolicy(machinePolicy.MachineCleanupPolicy)); err != nil {
		return fmt.Errorf("error setting machine_cleanup_policy: %s", err)
	}

	if err := d.Set("machine_connectivity_policy", flattenMachineConnectivityPolicy(machinePolicy.MachineConnectivityPolicy)); err != nil {
		return fmt.Errorf("error setting machine_connectivity_policy: %s", err)
	}

	if err := d.Set("machine_health_check_policy", flattenMachineHealthCheckPolicy(machinePolicy.MachineHealthCheckPolicy)); err != nil {
		return fmt.Errorf("error setting machine_health_check_policy: %s", err)
	}

	if err := d.Set("machine_update_policy", flattenMachineUpdatePolicy(machinePolicy.MachineUpdatePolicy)); err != nil {
		return fmt.Errorf("error setting machine_update_policy: %s", err)
	}

	return nil
}
