package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceProjectScheduledTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectScheduledTriggerCreate,
		DeleteContext: resourceProjectScheduledTriggerDelete,
		Description:   "This resource manages a scheduled trigger for a project or runbook in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceProjectScheduledTriggerRead,
		Schema:        getProjectScheduledTriggerSchema(),
		UpdateContext: resourceProjectScheduledTriggerUpdate,
	}
}

func resourceProjectScheduledTriggerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	spaceId := d.Get("space_id").(string)
	spaceId = util.Ternary(len(spaceId) > 0, spaceId, client.GetSpaceID())

	scheduledTrigger, err := triggers.GetById(client, spaceId, d.Id())

	if scheduledTrigger == nil {
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

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
	projectId := d.Get("project_id").(string)
	spaceId := d.Get("space_id").(string)
	project, err := projects.GetByID(client, spaceId, projectId)
	if err != nil {
		return diag.FromErr(err)
	}

	expandedScheduledTrigger, err := expandProjectScheduledTrigger(d, project)

	if err != nil {
		return diag.FromErr(err)
	}

	scheduledTrigger, err := triggers.Add(client, expandedScheduledTrigger)

	if err != nil {
		return diag.FromErr(err)
	}

	if isEmpty(scheduledTrigger.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(scheduledTrigger.GetID())
	}

	return nil
}

func resourceProjectScheduledTriggerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	projectId := d.Get("project_id").(string)
	spaceId := d.Get("space_id").(string)
	project, err := projects.GetByID(client, spaceId, projectId)
	if err != nil {
		return diag.FromErr(err)
	}

	expandedScheduledTrigger, err := expandProjectScheduledTrigger(d, project)

	if err != nil {
		return diag.FromErr(err)
	}

	expandedScheduledTrigger.ID = d.Id()

	if err != nil {
		return diag.FromErr(err)
	}

	scheduledTrigger, err := triggers.Update(client, expandedScheduledTrigger)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(scheduledTrigger.GetID())

	return nil
}

func resourceProjectScheduledTriggerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	spaceId := d.Get("space_id").(string)
	err := triggers.DeleteById(client, spaceId, d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
