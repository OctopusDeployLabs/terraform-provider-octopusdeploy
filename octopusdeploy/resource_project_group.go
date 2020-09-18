package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func buildProjectGroupResource(d *schema.ResourceData) *model.ProjectGroup {
	name := d.Get("name").(string)

	projectGroup := model.NewProjectGroup(name)

	if attr, ok := d.GetOk("description"); ok {
		projectGroup.Description = attr.(string)
	}

	return projectGroup
}

func resourceProjectGroupCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	newProjectGroup := buildProjectGroupResource(d)

	createdProjectGroup, err := apiClient.ProjectGroups.Add(newProjectGroup)

	if err != nil {
		return fmt.Errorf("error creating projectgroup: %s", err.Error())
	}

	d.SetId(createdProjectGroup.ID)
	return nil
}

func resourceProjectGroupRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	projectGroupID := d.Id()

	projectGroup, err := apiClient.ProjectGroups.Get(projectGroupID)

	if err == client.ErrItemNotFound {
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

	apiClient := m.(*client.Client)

	updatedProject, err := apiClient.ProjectGroups.Update(projectGroup)

	if err != nil {
		return fmt.Errorf("error updating projectgroup id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedProject.ID)
	return nil
}

func resourceProjectGroupDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	projectGroupID := d.Id()

	err := apiClient.ProjectGroups.Delete(projectGroupID)

	if err != nil {
		return fmt.Errorf("error deleting projectgroup id %s: %s", projectGroupID, err.Error())
	}

	d.SetId("")
	return nil
}
