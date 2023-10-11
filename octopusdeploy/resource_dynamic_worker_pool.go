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

	client := m.(*client.Client)
	createdWorkerPool, err := workerpools.Add(client, workerPool.SpaceID, workerPool)
	if err != nil {
		return diag.FromErr(err)
	}

	dynamicWorkerPool := createdWorkerPool.(*workerpools.DynamicWorkerPool)
	if err := setDynamicWorkerPool(ctx, d, dynamicWorkerPool); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdWorkerPool.GetID())

	log.Printf("[INFO] dynamic worker pool created (%s)", d.Id())
	return nil
}

func resourceDynamicWorkerPoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting dynamic worker pool (%s)", d.Id())
	spaceID := d.Get("space_id").(string)

	client := m.(*client.Client)
	if err := workerpools.DeleteByID(client, spaceID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] dynamic worker pool deleted")
	return nil
}

func resourceDynamicWorkerPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading dynamic worker pool (%s)", d.Id())

	spaceID := d.Get("space_id").(string)
	client := m.(*client.Client)
	workerPoolResource, err := workerpools.GetByID(client, spaceID, d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "dynamic worker pool")
	}

	dynamicWorkerPool := workerPoolResource.(*workerpools.DynamicWorkerPool)
	if err := setDynamicWorkerPool(ctx, d, dynamicWorkerPool); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] dynamic worker pool read (%s)", d.Id())
	return nil
}

func resourceDynamicWorkerPoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	workerPool := expandDynamicWorkerPool(d)

	log.Printf("[INFO] updating dynamic worker pool (%s)", d.Id())
	spaceID := d.Get("space_id").(string)

	client := m.(*client.Client)
	updatedWorkerPool, err := workerpools.Update(client, spaceID, workerPool)
	if err != nil {
		return diag.FromErr(err)
	}

	dynamicWorkerPool := updatedWorkerPool.(*workerpools.DynamicWorkerPool)
	if err := setDynamicWorkerPool(ctx, d, dynamicWorkerPool); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] dynamic worker pool updated (%s)", d.Id())
	return nil
}
