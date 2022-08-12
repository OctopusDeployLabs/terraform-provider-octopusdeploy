package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workerpools"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStaticWorkerPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStaticWorkerPoolCreate,
		DeleteContext: resourceStaticWorkerPoolDelete,
		Description:   "This resource manages static worker pools in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceStaticWorkerPoolRead,
		Schema:        getStaticWorkerPoolSchema(),
		UpdateContext: resourceStaticWorkerPoolUpdate,
	}
}

func resourceStaticWorkerPoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workerPool := expandStaticWorkerPool(d)

	log.Printf("[INFO] creating static worker pool: %#v", workerPool)

	client := m.(*client.Client)
	createdWorkerPool, err := client.WorkerPools.Add(workerPool)
	if err != nil {
		return diag.FromErr(err)
	}

	staticWorkerPool := createdWorkerPool.(*workerpools.StaticWorkerPool)
	if err := setStaticWorkerPool(ctx, d, staticWorkerPool); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdWorkerPool.GetID())

	log.Printf("[INFO] static worker pool created (%s)", d.Id())
	return nil
}

func resourceStaticWorkerPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting static worker pool (%s)", d.Id())

	client := m.(*client.Client)
	if err := client.WorkerPools.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] static worker pool deleted")
	return nil
}

func resourceStaticWorkerPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading static worker pool (%s)", d.Id())

	client := m.(*client.Client)
	workerPoolResource, err := client.WorkerPools.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "static worker pool")
	}

	staticWorkerPool := workerPoolResource.(*workerpools.StaticWorkerPool)
	if err := setStaticWorkerPool(ctx, d, staticWorkerPool); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] static worker pool read (%s)", d.Id())
	return nil
}

func resourceStaticWorkerPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workerPool := expandStaticWorkerPool(d)

	log.Printf("[INFO] updating static worker pool (%s)", d.Id())

	client := m.(*client.Client)
	updatedWorkerPool, err := client.WorkerPools.Update(workerPool)
	if err != nil {
		return diag.FromErr(err)
	}

	staticWorkerPool := updatedWorkerPool.(*workerpools.StaticWorkerPool)
	if err := setStaticWorkerPool(ctx, d, staticWorkerPool); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] static worker pool updated (%s)", d.Id())
	return nil
}
