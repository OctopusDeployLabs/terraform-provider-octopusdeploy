package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceDeploymentProcess() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeploymentProcessCreate,
		Read:   resourceDeploymentProcessRead,
		Update: resourceDeploymentProcessUpdate,
		Delete: resourceDeploymentProcessDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"step": getDeploymentStepSchema(),
		},
	}
}

func resourceDeploymentProcessCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newDeploymentProcess := buildDeploymentProcessResource(d)

	project, err := client.Project.Get(newDeploymentProcess.ProjectID)
	if err != nil {
		return fmt.Errorf("error getting project %s: %s", project.Name, err.Error())
	}

	current, err := client.DeploymentProcess.Get(project.DeploymentProcessID)
	if err != nil {
		return fmt.Errorf("error getting deployment process for %s: %s", project.Name, err.Error())
	}

	newDeploymentProcess.ID = current.ID
	newDeploymentProcess.Version = current.Version
	createdDeploymentProcess, err := client.DeploymentProcess.Update(newDeploymentProcess)

	if err != nil {
		return fmt.Errorf("error creating deployment process: %s", err.Error())
	}

	d.SetId(createdDeploymentProcess.ID)

	return nil
}

func buildDeploymentProcessResource(d *schema.ResourceData) *octopusdeploy.DeploymentProcess {
	deploymentProcess := &octopusdeploy.DeploymentProcess{
		ProjectID: d.Get("project_id").(string),
	}

	if attr, ok := d.GetOk("step"); ok {
		tfSteps := attr.([]interface{})

		for _, tfStep := range tfSteps {
			step := buildDeploymentStepResource(tfStep.(map[string]interface{}))
			deploymentProcess.Steps = append(deploymentProcess.Steps, step)
		}
	}

	return deploymentProcess
}

func resourceDeploymentProcessRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	deploymentProcessID := d.Id()

	_, err := client.DeploymentProcess.Get(deploymentProcessID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading deployment process id %s: %s", deploymentProcessID, err.Error())
	}

	log.Printf("[DEBUG] deploymentProcess: %v", m)

	return nil
}

func resourceDeploymentProcessUpdate(d *schema.ResourceData, m interface{}) error {
	deploymentProcess := buildDeploymentProcessResource(d)
	deploymentProcess.ID = d.Id() // set deploymentProcess struct ID so octopus knows which deploymentProcess to update

	client := m.(*octopusdeploy.Client)

	current, err := client.DeploymentProcess.Get(deploymentProcess.ID)
	if err != nil {
		return fmt.Errorf("error getting deployment process %s: %s", deploymentProcess.ID, err.Error())
	}

	deploymentProcess.Version = current.Version
	deploymentProcess, err = client.DeploymentProcess.Update(deploymentProcess)

	if err != nil {
		return fmt.Errorf("error updating deployment process id %s: %s", d.Id(), err.Error())
	}

	d.SetId(deploymentProcess.ID)

	return nil
}

func resourceDeploymentProcessDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)
	current, err := client.DeploymentProcess.Get(d.Id())

	if err != nil {
		return fmt.Errorf("error getting deployment process with id %s: %s", d.Id(), err.Error())
	}

	deploymentProcess := &octopusdeploy.DeploymentProcess{
		ID:      d.Id(),
		Version: current.Version,
	}

	deploymentProcess, err = client.DeploymentProcess.Update(deploymentProcess)

	if err != nil {
		return fmt.Errorf("error deleting deployment process with id %s: %s", deploymentProcess.ID, err.Error())
	}

	d.SetId("")
	return nil
}
