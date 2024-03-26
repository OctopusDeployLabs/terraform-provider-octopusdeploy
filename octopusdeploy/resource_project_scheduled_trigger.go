package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceProjectScheduledTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectScheduledTriggerCreate,
		DeleteContext: resourceProjectScheduledTriggerDelete,
		Importer:      getImporter(),
		ReadContext:   resourceProjectScheduledTriggerRead,
		Schema:        getProjectScheduledTriggerSchema(),
		UpdateContext: resourceProjectScheduledTriggerUpdate,
	}
}

func resourceProjectScheduledTriggerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	scheduledTrigger, err := client.ProjectTriggers.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if scheduledTrigger == nil {
		d.SetId("")
		return nil
	}

	flattenedScheduledTrigger := flattenProjectScheduledTrigger(scheduledTrigger)
	for key, value := range flattenedScheduledTrigger {
		err := d.Set(key, value)
		if err != nil {
			return nil
		}
	}

	return nil
}

func resourceProjectScheduledTriggerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	expandedScheduledTrigger, err := expandProjectScheduledTrigger(d, client)

	if err != nil {
		return diag.FromErr(err)
	}

	scheduledTrigger, err := client.ProjectTriggers.Add(expandedScheduledTrigger)

	if err != nil {
		return diag.FromErr(err)
	}

	if isEmpty(scheduledTrigger.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(scheduledTrigger.GetID())
	}

	return resourceProjectScheduledTriggerRead(ctx, d, m)
}

func resourceProjectScheduledTriggerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	expandedScheduledTrigger, err := expandProjectScheduledTrigger(d, client)

	if err != nil {
		return diag.FromErr(err)
	}

	expandedScheduledTrigger.ID = d.Id()

	if err != nil {
		return diag.FromErr(err)
	}

	scheduledTrigger, err := client.ProjectTriggers.Update(expandedScheduledTrigger)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(scheduledTrigger.GetID())

	return resourceProjectScheduledTriggerRead(ctx, d, m)
}

func resourceProjectScheduledTriggerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	err := client.ProjectTriggers.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
