package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDynamicWorkerPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDynamicWorkerPoolCreate,
		DeleteContext: resourceDynamicWorkerPoolDelete,
		Description:   "This resource manages dynamic worker pools in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceDynamicWorkerPoolRead,
		Schema:        getDynamicWorkerPoolSchema(),
		UpdateContext: resourceDynamicWorkerPoolUpdate,
	}
}

func resourceDynamicWorkerPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workerPool := expandDynamicWorkerPool(d)

	log.Printf("[INFO] creating dynamic worker pool: %#v", workerPool)

	client := m.(*octopusdeploy.Client)
	createdWorkerPool, err := client.WorkerPools.Add(workerPool)
	if err != nil {
		return diag.FromErr(err)
	}

	dynamicWorkerPool := createdWorkerPool.(*octopusdeploy.DynamicWorkerPool)
	if err := setDynamicWorkerPool(ctx, d, dynamicWorkerPool); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdWorkerPool.GetID())

	log.Printf("[INFO] dynamic worker pool created (%s)", d.Id())
	return nil
}

func resourceDynamicWorkerPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting dynamic worker pool (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.WorkerPools.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] dynamic worker pool deleted")
	return nil
}

func resourceDynamicWorkerPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading dynamic worker pool (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	workerPoolResource, err := client.WorkerPools.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] dynamic worker pool (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	dynamicWorkerPool := workerPoolResource.(*octopusdeploy.DynamicWorkerPool)
	if err := setDynamicWorkerPool(ctx, d, dynamicWorkerPool); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] dynamic worker pool read (%s)", d.Id())
	return nil
}

func resourceDynamicWorkerPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workerPool := expandDynamicWorkerPool(d)

	log.Printf("[INFO] updating dynamic worker pool (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	updatedWorkerPool, err := client.WorkerPools.Update(workerPool)
	if err != nil {
		return diag.FromErr(err)
	}

	dynamicWorkerPool := updatedWorkerPool.(*octopusdeploy.DynamicWorkerPool)
	if err := setDynamicWorkerPool(ctx, d, dynamicWorkerPool); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] dynamic worker pool updated (%s)", d.Id())
	return nil
}
