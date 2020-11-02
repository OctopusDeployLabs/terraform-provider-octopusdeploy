package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectGroupCreate,
		Read:   resourceProjectGroupRead,
		Update: resourceProjectGroupUpdate,
		Delete: resourceProjectGroupDelete,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildProjectGroupResource(d *schema.ResourceData) *octopusdeploy.ProjectGroup {
	name := d.Get(constName).(string)

	projectGroup := octopusdeploy.NewProjectGroup(name)

	if attr, ok := d.GetOk(constDescription); ok {
		projectGroup.Description = attr.(string)
	}

	return projectGroup
}

func resourceProjectGroupCreate(d *schema.ResourceData, m interface{}) error {
	projectGroup := buildProjectGroupResource(d)

	client := m.(*octopusdeploy.Client)
	resource, err := client.ProjectGroups.Add(projectGroup)
	if err != nil {
		return createResourceOperationError(errorCreatingProjectGroup, projectGroup.ID, err)
		// return diag.FromErr(err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func resourceProjectGroupRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	resource, err := client.ProjectGroups.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingProjectGroup, id, err)
		// diag.Errorf(err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constProjectGroup, m)

	d.Set(constName, resource.Name)
	d.Set(constDescription, resource.Description)

	return nil
}

func resourceProjectGroupUpdate(d *schema.ResourceData, m interface{}) error {
	projectGroup := buildProjectGroupResource(d)
	projectGroup.ID = d.Id() // set ID so Octopus API knows which project group to update

	client := m.(*octopusdeploy.Client)
	resource, err := client.ProjectGroups.Update(*projectGroup)
	if err != nil {
		return createResourceOperationError(errorUpdatingProjectGroup, d.Id(), err)
		// diag.Errorf(err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceProjectGroupDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.ProjectGroups.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingProjectGroup, id, err)
		// diag.Errorf(err)
	}

	d.SetId(constEmptyString)

	return nil
}
