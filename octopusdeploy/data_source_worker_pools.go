package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workerpools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceWorkerPools() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing worker pools.",
		ReadContext: dataSourceWorkerPoolsRead,
		Schema:      getWorkerPoolDataSchema(),
	}
}

func dataSourceWorkerPoolsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := workerpools.WorkerPoolsQuery{
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}
	name := d.Get("name").(string)

	client := m.(*client.Client)
	workerPools, err := workerpools.Get(client, d.Get("space_id").(string), query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedWorkerPools := []interface{}{}
	for _, workerPool := range workerPools.Items {
		workerPoolResource, err := workerpools.ToWorkerPoolResource(workerPool)
		if err != nil {
			return diag.FromErr(err)
		}

		// There is no name filter on the WorkerPools endpoint in the Octopus API,
		// so filter in-memory if the name field is specified
		if name != "" && workerPoolResource.Name != name {
			continue
		}
		flattenedWorkerPools = append(flattenedWorkerPools, flattenWorkerPool(workerPoolResource))
	}

	d.Set("worker_pools", flattenedWorkerPools)
	d.SetId("Worker Pools " + time.Now().UTC().String())

	return nil
}
