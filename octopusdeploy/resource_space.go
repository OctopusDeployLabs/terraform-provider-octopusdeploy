package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSpace() *schema.Resource {
	resourceSpaceImporter := &schema.ResourceImporter{
		StateContext: schema.ImportStatePassthroughContext,
	}
	resourceSpaceSchema := map[string]*schema.Schema{
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"is_default": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"space_managers_team_members": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		"space_managers_teams": {
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
			Type:     schema.TypeList,
		},
		"task_queue_stopped": {
			Optional: true,
			Type:     schema.TypeBool,
		},
	}

	return &schema.Resource{
		CreateContext: resourceSpaceCreate,
		DeleteContext: resourceSpaceDelete,
		Importer:      resourceSpaceImporter,
		ReadContext:   resourceSpaceRead,
		Schema:        resourceSpaceSchema,
		UpdateContext: resourceSpaceUpdate,
	}
}

func buildSpace(d *schema.ResourceData) *octopusdeploy.Space {
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	space := octopusdeploy.NewSpace(name)
	space.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		space.Description = v.(string)
	}

	if v, ok := d.GetOk("is_default"); ok {
		space.IsDefault = v.(bool)
	}

	if v, ok := d.GetOk("space_managers_team_members"); ok {
		space.SpaceManagersTeamMembers = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("space_managers_teams"); ok {
		space.SpaceManagersTeams = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("task_queue_stopped"); ok {
		space.TaskQueueStopped = v.(bool)
	}

	return space
}

func resourceSpaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	space := buildSpace(d)

	client := m.(*octopusdeploy.Client)
	createdSpace, err := client.Spaces.Add(space)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenSpace(ctx, d, createdSpace)
	return nil
}

func resourceSpaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	space, err := client.Spaces.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	flattenSpace(ctx, d, space)
	return nil
}

func resourceSpaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	space := buildSpace(d)

	client := m.(*octopusdeploy.Client)
	updatedSpace, err := client.Spaces.Update(space)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenSpace(ctx, d, updatedSpace)
	return nil
}

func resourceSpaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	space := buildSpace(d)
	space.TaskQueueStopped = true

	client := m.(*octopusdeploy.Client)
	updatedSpace, err := client.Spaces.Update(space)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.Spaces.DeleteByID(updatedSpace.GetID())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func flattenSpace(ctx context.Context, d *schema.ResourceData, space *octopusdeploy.Space) {
	d.Set("description", space.Description)
	d.Set("is_default", space.IsDefault)
	d.Set("name", space.Name)
	d.Set("space_managers_team_members", space.SpaceManagersTeamMembers)
	d.Set("space_managers_teams", space.SpaceManagersTeams)
	d.Set("task_queue_stopped", space.TaskQueueStopped)

	d.SetId(space.GetID())
}
