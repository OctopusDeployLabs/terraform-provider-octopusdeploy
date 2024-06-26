package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandWorker(d *schema.ResourceData) *machines.Worker {
	endpoint := expandEndpoint(d.Get("endpoint"))
	name := d.Get("name").(string)

	worker := machines.NewWorker(name, endpoint)
	worker.ID = d.Id()

	if v, ok := d.GetOk("machine_policy_id"); ok {
		worker.MachinePolicyID = v.(string)
	}

	if v, ok := d.GetOk("is_disabled"); ok {
		worker.IsDisabled = v.(bool)
	}

	if v, ok := d.GetOk("thumbprint"); ok {
		worker.Thumbprint = v.(string)
	}

	if v, ok := d.GetOk("uri"); ok {
		worker.URI = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		worker.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("thumbprint"); ok {
		worker.Thumbprint = v.(string)
	}

	worker.WorkerPoolIDs = getSliceFromTerraformTypeList(d.Get("worker_pool_ids"))

	return worker
}

func flattenWorker(worker *machines.Worker) map[string]interface{} {
	if worker == nil {
		return nil
	}

	endpointResource, _ := machines.ToEndpointResource(worker.Endpoint)

	return map[string]interface{}{
		"endpoint":            flattenEndpointResource(endpointResource),
		"has_latest_calamari": worker.HasLatestCalamari,
		"health_status":       worker.HealthStatus,
		"id":                  worker.GetID(),
		"is_disabled":         worker.IsDisabled,
		"is_in_process":       worker.IsInProcess,
		"machine_policy_id":   worker.MachinePolicyID,
		"name":                worker.Name,
		"operating_system":    worker.OperatingSystem,
		"shell_name":          worker.ShellName,
		"shell_version":       worker.ShellVersion,
		"space_id":            worker.SpaceID,
		"status":              worker.Status,
		"status_summary":      worker.StatusSummary,
		"thumbprint":          worker.Thumbprint,
		"uri":                 worker.URI,
		"worker_pool_ids":     worker.WorkerPoolIDs,
	}
}

func getWorkerDataSchema() map[string]*schema.Schema {
	dataSchema := getWorkerSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"communication_styles": getQueryCommunicationStyles(),
		"deployment_id":        getQueryDeploymentID(),
		"workers": {
			Computed:    true,
			Description: "A list of deployment targets that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"health_statuses": getQueryHealthStatuses(),
		"ids":             getQueryIDs(),
		"is_disabled":     getQueryIsDisabled(),
		"name":            getQueryName(),
		"partial_name":    getQueryPartialName(),
		"roles":           getQueryRoles(),
		"shell_names":     getQueryShellNames(),
		"skip":            getQuerySkip(),
		"take":            getQueryTake(),
		"thumbprint":      getQueryThumbprint(),
		"space_id":        getSpaceIDSchema(),
	}
}

func getWorkerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"endpoint": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getEndpointSchema()},
			MinItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"has_latest_calamari": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"health_status": getHealthStatusSchema(),
		"id":            getIDSchema(),
		"is_disabled": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeBool,
		},
		"is_in_process": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"machine_policy_id": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"name": getNameSchema(true),
		"operating_system": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"shell_name": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"shell_version": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"space_id":       getSpaceIDSchema(),
		"status":         getStatusSchema(),
		"status_summary": getStatusSummarySchema(),
		"thumbprint": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"uri": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"worker_pool_ids": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			MinItems: 1,
			Required: true,
			Type:     schema.TypeList,
		},
	}
}

func setWorker(ctx context.Context, d *schema.ResourceData, worker *machines.Worker) error {
	d.Set("has_latest_calamari", worker.HasLatestCalamari)
	d.Set("health_status", worker.HealthStatus)
	d.Set("is_disabled", worker.IsDisabled)
	d.Set("is_in_process", worker.IsInProcess)
	d.Set("machine_policy_id", worker.MachinePolicyID)
	d.Set("name", worker.Name)
	d.Set("operating_system", worker.OperatingSystem)
	d.Set("shell_name", worker.ShellName)
	d.Set("shell_version", worker.ShellVersion)
	d.Set("space_id", worker.SpaceID)
	d.Set("status", worker.Status)
	d.Set("status_summary", worker.StatusSummary)
	d.Set("thumbprint", worker.Thumbprint)
	d.Set("uri", worker.URI)
	d.Set("space_id", worker.SpaceID)
	d.Set("worker_pool_ids", worker.WorkerPoolIDs)

	endpointResource, err := machines.ToEndpointResource(worker.Endpoint)
	if err != nil {
		return fmt.Errorf("error setting endpoint: %s", err)
	}

	if d.Get("endpoint") != nil {
		if err := d.Set("endpoint", flattenEndpointResource(endpointResource)); err != nil {
			return fmt.Errorf("error setting endpoint: %s", err)
		}
	}

	d.SetId(worker.GetID())

	return nil
}
