package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workerpools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandDynamicWorkerPool(d *schema.ResourceData) *workerpools.DynamicWorkerPool {
	name := d.Get("name").(string)
	workerType := workerpools.WorkerType(d.Get("worker_type").(string))

	dynamicWorkerPool := workerpools.NewDynamicWorkerPool(name, workerType)
	dynamicWorkerPool.ID = d.Id()

	if v, ok := d.GetOk("can_add_workers"); ok {
		dynamicWorkerPool.CanAddWorkers = v.(bool)
	}

	if v, ok := d.GetOk("description"); ok {
		dynamicWorkerPool.Description = v.(string)
	}

	if v, ok := d.GetOk("is_default"); ok {
		dynamicWorkerPool.IsDefault = v.(bool)
	}

	if v, ok := d.GetOk("sort_order"); ok {
		dynamicWorkerPool.SortOrder = v.(int)
	}

	if v, ok := d.GetOk("space_id"); ok {
		dynamicWorkerPool.SpaceID = v.(string)
	}

	return dynamicWorkerPool
}

func flattenDynamicWorkerPool(dynamicWorkerPool *workerpools.DynamicWorkerPool) map[string]interface{} {
	if dynamicWorkerPool == nil {
		return nil
	}

	return map[string]interface{}{
		"can_add_workers": dynamicWorkerPool.CanAddWorkers,
		"description":     dynamicWorkerPool.Description,
		"id":              dynamicWorkerPool.GetID(),
		"is_default":      dynamicWorkerPool.IsDefault,
		"name":            dynamicWorkerPool.GetName(),
		"space_id":        dynamicWorkerPool.SpaceID,
		"sort_order":      dynamicWorkerPool.SortOrder,
		"worker_type":     dynamicWorkerPool.WorkerType,
	}
}

func getDynamicWorkerPoolDataSchema() map[string]*schema.Schema {
	dataSchema := getDynamicWorkerPoolSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"filter": getQueryFilter(),
		"id":     getDataSchemaID(),
		"ids":    getQueryIDs(),
		"skip":   getQuerySkip(),
		"take":   getQueryTake(),
		"dynamic_worker_pools": {
			Computed:    true,
			Description: "A list of dynamic worker pools that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}

func getDynamicWorkerPoolSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"can_add_workers": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"description": getDescriptionSchema("dynamic worker pool"),
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
		"worker_type": {
			Required: true,
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

func setDynamicWorkerPool(ctx context.Context, d *schema.ResourceData, dynamicWorkerPool *workerpools.DynamicWorkerPool) error {
	d.Set("can_add_workers", dynamicWorkerPool.CanAddWorkers)
	d.Set("description", dynamicWorkerPool.Description)
	d.Set("is_default", dynamicWorkerPool.IsDefault)
	d.Set("name", dynamicWorkerPool.Name)
	d.Set("space_id", dynamicWorkerPool.SpaceID)
	d.Set("sort_order", dynamicWorkerPool.SortOrder)
	d.Set("worker_type", dynamicWorkerPool.WorkerType)

	d.SetId(dynamicWorkerPool.GetID())

	return nil
}
