package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectGroupCreate,
		ReadContext:   resourceProjectGroupRead,
		UpdateContext: resourceProjectGroupUpdate,
		DeleteContext: resourceProjectGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"description": {
				Optional: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func buildProjectGroupResource(d *schema.ResourceData) *octopusdeploy.ProjectGroup {
	name := d.Get("name").(string)

	projectGroup := octopusdeploy.NewProjectGroup(name)

	if attr, ok := d.GetOk("description"); ok {
		projectGroup.Description = attr.(string)
	}

	return projectGroup
}

func resourceProjectGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectGroup := buildProjectGroupResource(d)

	client := m.(*octopusdeploy.Client)
	resource, err := client.ProjectGroups.Add(projectGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func resourceProjectGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	resource, err := client.ProjectGroups.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	if resource == nil {
		d.SetId("")
		return nil
	}

	logResource(constProjectGroup, m)

	d.Set("name", resource.Name)
	d.Set("description", resource.Description)

	return nil
}

func resourceProjectGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectGroup := buildProjectGroupResource(d)
	projectGroup.ID = d.Id() // set ID so Octopus API knows which project group to update

	client := m.(*octopusdeploy.Client)
	resource, err := client.ProjectGroups.Update(*projectGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceProjectGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.ProjectGroups.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
