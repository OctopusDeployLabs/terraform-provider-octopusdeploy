package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/MattHodge/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"lifecycleid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"projectgroupid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)
	name := d.Get("name").(string)
	lifecycleID := d.Get("lifecycleid").(string)
	projectGroupID := d.Get("projectgroupid").(string)

	p := octopusdeploy.NewProject(name, lifecycleID, projectGroupID)

	p.Description = d.Get("description").(string)

	createdProject, err := client.Projects.Add(p)

	if err != nil {
		return fmt.Errorf("error creating project: %s", err.Error())
	}

	d.SetId(createdProject.ID)
	return nil
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	projectID := d.Id()

	project, err := client.Projects.Get(projectID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading project id %s: %s", projectID, err.Error())
	}

	log.Printf("[DEBUG] project: %v", m)
	d.Set("name", project.Name)
	d.Set("description", project.Description)
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	name := d.Get("name").(string)
	lifecycleID := d.Get("lifecycleid").(string)
	projectGroupID := d.Get("projectgroupid").(string)
	p := octopusdeploy.NewProject(name, lifecycleID, projectGroupID)

	p.ID = d.Id() // set project struct ID so octopus knows which project to update

	if attr, ok := d.GetOk("description"); ok {
		p.Description = attr.(string)
	}

	project, err := client.Projects.Update(*p)

	if err != nil {
		return fmt.Errorf("error updating project id %s: %s", d.Id(), err.Error())
	}

	d.SetId(project.ID)
	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	projectID := d.Id()

	err := client.Projects.Delete(projectID)

	if err != nil {
		return fmt.Errorf("error deleting project id %s: %s", projectID, err.Error())
	}

	d.SetId("")
	return nil
}
