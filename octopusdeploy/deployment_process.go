package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDeploymentProcess() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeploymentProcessCreate,
		Read:   resourceDeploymentProcessRead,
		Update: resourceDeploymentProcessUpdate,
		Delete: resourceDeploymentProcessDelete,

		Schema: map[string]*schema.Schema{
			constProjectID: {
				Type:     schema.TypeString,
				Required: true,
			},
			"step": getDeploymentStepSchema(),
		},
	}
}

func resourceDeploymentProcessCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	newDeploymentProcess := buildDeploymentProcessResource(d)

	project, err := apiClient.Projects.GetByID(newDeploymentProcess.ProjectID)
	if err != nil {
		return fmt.Errorf("error getting project %s: %s", project.Name, err.Error())
	}

	current, err := apiClient.DeploymentProcesses.GetByID(project.DeploymentProcessID)
	if err != nil {
		return fmt.Errorf("error getting deployment process for %s: %s", project.Name, err.Error())
	}

	newDeploymentProcess.ID = current.ID
	newDeploymentProcess.Version = current.Version
	createdDeploymentProcess, err := apiClient.DeploymentProcesses.Update(*newDeploymentProcess)

	if err != nil {
		return fmt.Errorf("error creating deployment process: %s", err.Error())
	}

	d.SetId(createdDeploymentProcess.ID)

	return nil
}

func buildDeploymentProcessResource(d *schema.ResourceData) *model.DeploymentProcess {
	deploymentProcess := &model.DeploymentProcess{
		ProjectID: d.Get(constProjectID).(string),
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
	apiClient := m.(*client.Client)

	deploymentProcessID := d.Id()

	resource, err := apiClient.DeploymentProcesses.GetByID(deploymentProcessID)
	if err != nil {
		return fmt.Errorf("error reading deployment process id %s: %s", deploymentProcessID, err.Error())
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	log.Printf("[DEBUG] deploymentProcess: %v", m)

	return nil
}

func resourceDeploymentProcessUpdate(d *schema.ResourceData, m interface{}) error {
	deploymentProcess := buildDeploymentProcessResource(d)
	deploymentProcess.ID = d.Id() // set deploymentProcess struct ID so octopus knows which deploymentProcess to update

	apiClient := m.(*client.Client)

	current, err := apiClient.DeploymentProcesses.GetByID(deploymentProcess.ID)
	if err != nil {
		return fmt.Errorf("error getting deployment process %s: %s", deploymentProcess.ID, err.Error())
	}

	deploymentProcess.Version = current.Version
	deploymentProcess, err = apiClient.DeploymentProcesses.Update(*deploymentProcess)

	if err != nil {
		return fmt.Errorf("error updating deployment process id %s: %s", d.Id(), err.Error())
	}

	d.SetId(deploymentProcess.ID)

	return nil
}

func resourceDeploymentProcessDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	current, err := apiClient.DeploymentProcesses.GetByID(d.Id())

	if err != nil {
		return fmt.Errorf("error getting deployment process with id %s: %s", d.Id(), err.Error())
	}

	deploymentProcess := &model.DeploymentProcess{
		Version: current.Version,
	}
	deploymentProcess.ID = d.Id()

	deploymentProcess, err = apiClient.DeploymentProcesses.Update(*deploymentProcess)

	if err != nil {
		return fmt.Errorf("error deleting deployment process with id %s: %s", deploymentProcess.ID, err.Error())
	}

	d.SetId(constEmptyString)
	return nil
}
