package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/workerpools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func flattenWorkerPool(workerPool *workerpools.WorkerPoolResource) map[string]interface{} {
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
		"space_id":     getSpaceIDSchema(),
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
		"description": getDescriptionSchema("worker pool"),
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
				"Ubuntu2204",
				"UbuntuDefault",
				"Windows2016",
				"Windows2019",
				"Windows2022",
				"WindowsDefault",
			}, false)),
		},
	}
}
