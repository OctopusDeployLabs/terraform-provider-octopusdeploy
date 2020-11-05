package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSpace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSpaceCreate,
		DeleteContext: resourceSpaceDelete,
		ReadContext:   resourceSpaceRead,
		UpdateContext: resourceSpaceUpdate,

		Schema: map[string]*schema.Schema{
			constDescription: {
				Optional: true,
				Type:     schema.TypeString,
			},
			constID: {
				Computed: true,
				Optional: true,
				Type:     schema.TypeString,
			},
			constIsDefault: {
				Optional: true,
				Type:     schema.TypeBool,
			},
			constName: {
				Required: true,
				Type:     schema.TypeString,
			},
			constSpaceManagersTeamMembers: {
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Type:     schema.TypeList,
			},
			constSpaceManagersTeams: {
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				Type:     schema.TypeList,
			},
			constTaskQueueStopped: {
				Optional: true,
				Type:     schema.TypeBool,
			},
		},
	}
}

func buildSpace(d *schema.ResourceData) *octopusdeploy.Space {
	var name string
	if v, ok := d.GetOk(constName); ok {
		name = v.(string)
	}

	space := octopusdeploy.NewSpace(name)

	if v, ok := d.GetOk(constDescription); ok {
		space.Description = v.(string)
	}

	if v, ok := d.GetOk(constID); ok {
		space.ID = v.(string)
	}

	if v, ok := d.GetOk(constIsDefault); ok {
		space.IsDefault = v.(bool)
	}

	if v, ok := d.GetOk(constSpaceManagersTeamMembers); ok {
		space.SpaceManagersTeamMembers = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constSpaceManagersTeams); ok {
		space.SpaceManagersTeams = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk(constTaskQueueStopped); ok {
		space.TaskQueueStopped = v.(bool)
	}

	return space
}

func resourceSpaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	space := buildSpace(d)

	client := m.(*octopusdeploy.Client)
	space, err := client.Spaces.Add(space)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(space.GetID())

	return nil
}

func resourceSpaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	space, err := client.Spaces.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(constDescription, space.Description)
	d.Set(constIsDefault, space.IsDefault)
	d.Set(constName, space.Name)
	d.Set(constSpaceManagersTeamMembers, space.SpaceManagersTeamMembers)
	d.Set(constSpaceManagersTeams, space.SpaceManagersTeams)
	d.Set(constTaskQueueStopped, space.TaskQueueStopped)

	return nil
}

func resourceSpaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	space := buildSpace(d)
	space.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	updatedSpace, err := client.Spaces.Update(space)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(updatedSpace.GetID())

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

	d.SetId(constEmptyString)

	return nil
}
