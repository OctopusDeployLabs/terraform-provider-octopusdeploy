package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandStaticWorkerPool(d *schema.ResourceData) *octopusdeploy.StaticWorkerPool {
	name := d.Get("name").(string)

	dynamicWorkerPool := octopusdeploy.NewStaticWorkerPool(name)
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

func flattenStaticWorkerPool(staticWorkerPool *octopusdeploy.StaticWorkerPool) map[string]interface{} {
	if staticWorkerPool == nil {
		return nil
	}

	return map[string]interface{}{
		"can_add_workers": staticWorkerPool.CanAddWorkers,
		"description":     staticWorkerPool.Description,
		"id":              staticWorkerPool.GetID(),
		"is_default":      staticWorkerPool.IsDefault,
		"name":            staticWorkerPool.GetName(),
		"space_id":        staticWorkerPool.SpaceID,
		"sort_order":      staticWorkerPool.SortOrder,
	}
}

func getStaticWorkerPoolDataSchema() map[string]*schema.Schema {
	dataSchema := getStaticWorkerPoolSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"filter": getQueryFilter(),
		"id":     getDataSchemaID(),
		"ids":    getQueryIDs(),
		"skip":   getQuerySkip(),
		"take":   getQueryTake(),
		"static_worker_pools": {
			Computed:    true,
			Description: "A list of static worker pools that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}

func getStaticWorkerPoolSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"can_add_workers": {
			Computed: true,
			Type:     schema.TypeBool,
		},
		"description": getDescriptionSchema("static worker pool"),
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
	}
}

func setStaticWorkerPool(ctx context.Context, d *schema.ResourceData, staticWorkerPool *octopusdeploy.StaticWorkerPool) error {
	d.Set("can_add_workers", staticWorkerPool.CanAddWorkers)
	d.Set("description", staticWorkerPool.Description)
	d.Set("is_default", staticWorkerPool.IsDefault)
	d.Set("name", staticWorkerPool.Name)
	d.Set("space_id", staticWorkerPool.SpaceID)
	d.Set("sort_order", staticWorkerPool.SortOrder)

	d.SetId(staticWorkerPool.GetID())

	return nil
}
