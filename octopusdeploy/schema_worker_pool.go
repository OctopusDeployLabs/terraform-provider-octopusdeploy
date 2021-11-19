package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandWorkerPool(d *schema.ResourceData) *octopusdeploy.WorkerPoolResource {
	name := d.Get("name").(string)
	workerPoolType := octopusdeploy.WorkerPoolType(d.Get("worker_pool_type").(string))

	workerPool := octopusdeploy.NewWorkerPoolResource(name, workerPoolType)
	workerPool.ID = d.Id()

	if v, ok := d.GetOk("can_add_workers"); ok {
		workerPool.CanAddWorkers = v.(bool)
	}

	if v, ok := d.GetOk("description"); ok {
		workerPool.Description = v.(string)
	}

	if v, ok := d.GetOk("is_default"); ok {
		workerPool.IsDefault = v.(bool)
	}

	if v, ok := d.GetOk("sort_order"); ok {
		workerPool.SortOrder = v.(int)
	}

	if v, ok := d.GetOk("worker_type"); ok {
		workerPool.WorkerType = octopusdeploy.WorkerType(v.(string))
	}

	return workerPool
}

func flattenWorkerPool(workerPool *octopusdeploy.WorkerPoolResource) map[string]interface{} {
	if workerPool == nil {
		return nil
	}

	return map[string]interface{}{
		"can_add_workers":  workerPool.CanAddWorkers,
		"description":      workerPool.Description,
		"id":               workerPool.GetID(),
		"is_default":       workerPool.IsDefault,
		"name":             workerPool.GetName(),
		"space_id":         workerPool.SpaceID,
		"sort_order":       workerPool.SortOrder,
		"worker_pool_type": workerPool.WorkerPoolType,
		"worker_type":      workerPool.WorkerType,
	}
}

func getWorkerPoolDataSchema() map[string]*schema.Schema {
	dataSchema := getWorkerPoolSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"ids":          getQueryIDs(),
		"name":         getQueryName(),
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
		"worker_pools": {
			Computed:    true,
			Description: "A list of worker pools that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}

func getWorkerPoolSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"can_add_workers": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"description": getDescriptionSchema(),
		"id":          getIDSchema(),
		"is_default": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"name": getNameSchema(true),
		"sort_order": {
			Computed:    true,
			Description: "The order number to sort a dynamic worker pool.",
			Optional:    true,
			Type:        schema.TypeInt,
		},
		"space_id": getSpaceIDSchema(),
		"worker_pool_type": {
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"DynamicWorkerPool",
				"StaticWorkerPool",
			}, false)),
		},
		"worker_type": {
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"Ubuntu1804",
				"UbuntuDefault",
				"Windows2016",
				"Windows2019",
				"WindowsDefault",
			}, false)),
		},
	}
}

func setWorkerPool(ctx context.Context, d *schema.ResourceData, workerPool *octopusdeploy.WorkerPoolResource) error {
	d.Set("can_add_workers", workerPool.CanAddWorkers)
	d.Set("description", workerPool.Description)
	d.Set("is_default", workerPool.IsDefault)
	d.Set("name", workerPool.Name)
	d.Set("space_id", workerPool.SpaceID)
	d.Set("sort_order", workerPool.SortOrder)
	d.Set("worker_pool_type", workerPool.WorkerPoolType)
	d.Set("worker_type", workerPool.WorkerType)

	d.SetId(workerPool.GetID())

	return nil
}
