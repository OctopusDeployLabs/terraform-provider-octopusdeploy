package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProjectGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectGroupCreate,
		Read:   resourceProjectGroupRead,
		Update: resourceProjectGroupUpdate,
		Delete: resourceProjectGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceProjectGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newProjectGroup := buildProjectGroupResource(d)

	createdProjectGroup, err := client.ProjectGroup.Add(newProjectGroup)

	if err != nil {
		return fmt.Errorf("error creating projectgroup: %s", err.Error())
	}

	d.SetId(createdProjectGroup.ID)
	return nil
}

func resourceProjectGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	projectGroupID := d.Id()

	projectGroup, err := client.ProjectGroup.Get(projectGroupID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading projectgroup id %s: %s", projectGroup.ID, err.Error())
	}

	log.Printf("[DEBUG] projectgroup: %v", m)
	d.Set("name", projectGroup.Name)
	d.Set("description", projectGroup.Description)
	return nil
}

func resourceProjectGroupUpdate(d *schema.ResourceData, m interface{}) error {
	projectGroup := buildProjectGroupResource(d)
	projectGroup.ID = d.Id() // set projectgroup struct ID so octopus knows which  to update

	client := m.(*octopusdeploy.Client)

	updatedProject, err := client.ProjectGroup.Update(projectGroup)

	if err != nil {
		return fmt.Errorf("error updating projectgroup id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedProject.ID)
	return nil
}

func resourceProjectGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	projectGroupID := d.Id()

	err := client.ProjectGroup.Delete(projectGroupID)

	if err != nil {
		return fmt.Errorf("error deleting projectgroup id %s: %s", projectGroupID, err.Error())
	}

	d.SetId("")
	return nil
}
