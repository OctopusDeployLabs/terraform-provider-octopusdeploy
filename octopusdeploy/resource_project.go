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
			"lifecycle_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"project_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_failure_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateValueFunc([]string{
					"EnvironmentDefault",
					"Off",
					"On",
				}),
			},
			"skip_machine_behavior": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateValueFunc([]string{
					"SkipUnavailableMachines",
					"None",
				}),
			},
		},
	}
}

func buildProjectResource(d *schema.ResourceData) *octopusdeploy.Project {
	name := d.Get("name").(string)
	lifecycleID := d.Get("lifecycle_id").(string)
	projectGroupID := d.Get("project_group_id").(string)

	project := octopusdeploy.NewProject(name, lifecycleID, projectGroupID)

	if attr, ok := d.GetOk("description"); ok {
		project.Description = attr.(string)
	}

	if attr, ok := d.GetOk("default_failure_mode"); ok {
		project.DefaultGuidedFailureMode = attr.(string)
	}

	if attr, ok := d.GetOk("skip_machine_behavior"); ok {
		project.ProjectConnectivityPolicy.SkipMachineBehavior = attr.(string)
	}

	return project
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newProject := buildProjectResource(d)

	createdProject, err := client.Project.Add(newProject)

	if err != nil {
		return fmt.Errorf("error creating project: %s", err.Error())
	}

	d.SetId(createdProject.ID)
	return nil
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	projectID := d.Id()

	project, err := client.Project.Get(projectID)

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
	project := buildProjectResource(d)
	project.ID = d.Id() // set project struct ID so octopus knows which project to update

	client := m.(*octopusdeploy.Client)

	project, err := client.Project.Update(project)

	if err != nil {
		return fmt.Errorf("error updating project id %s: %s", d.Id(), err.Error())
	}

	d.SetId(project.ID)
	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	projectID := d.Id()

	err := client.Project.Delete(projectID)

	if err != nil {
		return fmt.Errorf("error deleting project id %s: %s", projectID, err.Error())
	}

	d.SetId("")
	return nil
}
