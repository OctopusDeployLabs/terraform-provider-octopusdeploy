package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			constStep: getDeploymentStepSchema(),
		},
	}
}

func resourceDeploymentProcessCreate(d *schema.ResourceData, m interface{}) error {
	deploymentProcess := buildDeploymentProcessResource(d)

	client := m.(*octopusdeploy.Client)
	project, err := client.Projects.GetByID(deploymentProcess.ProjectID)
	if err != nil {
		return createResourceOperationError(errorReadingProject, project.Name, err)
	}

	current, err := client.DeploymentProcesses.GetByID(project.DeploymentProcessID)
	if err != nil {
		return createResourceOperationError(errorReadingDeploymentProcess, project.DeploymentProcessID, err)
	}

	deploymentProcess.ID = current.ID
	deploymentProcess.Version = current.Version

	resource, err := client.DeploymentProcesses.Update(*deploymentProcess)
	if err != nil {
		return createResourceOperationError(errorCreatingDeploymentProcess, deploymentProcess.ID, err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func buildDeploymentProcessResource(d *schema.ResourceData) *octopusdeploy.DeploymentProcess {
	deploymentProcess := octopusdeploy.NewDeploymentProcess(d.Get(constProjectID).(string))

	if attr, ok := d.GetOk(constStep); ok {
		tfSteps := attr.([]interface{})

		for _, tfStep := range tfSteps {
			step := buildDeploymentStepResource(tfStep.(map[string]interface{}))
			deploymentProcess.Steps = append(deploymentProcess.Steps, step)
		}
	}

	return deploymentProcess
}

func resourceDeploymentProcessRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	resource, err := client.DeploymentProcesses.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingDeploymentProcess, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constDeploymentProcess, m)

	return nil
}

func resourceDeploymentProcessUpdate(d *schema.ResourceData, m interface{}) error {
	deploymentProcess := buildDeploymentProcessResource(d)
	deploymentProcess.ID = d.Id() // set ID so Octopus API knows which deployment process to update

	client := m.(*octopusdeploy.Client)
	current, err := client.DeploymentProcesses.GetByID(deploymentProcess.ID)
	if err != nil {
		return createResourceOperationError(errorReadingDeploymentProcess, deploymentProcess.ID, err)
	}

	deploymentProcess.Version = current.Version
	resource, err := client.DeploymentProcesses.Update(*deploymentProcess)
	if err != nil {
		return createResourceOperationError(errorUpdatingDeploymentProcess, d.Id(), err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceDeploymentProcessDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)
	current, err := client.DeploymentProcesses.GetByID(d.Id())
	if err != nil {
		return createResourceOperationError(errorReadingDeploymentProcess, d.Id(), err)
	}

	deploymentProcess := &octopusdeploy.DeploymentProcess{
		Version: current.Version,
	}
	deploymentProcess.ID = d.Id()

	deploymentProcess, err = client.DeploymentProcesses.Update(*deploymentProcess)
	if err != nil {
		return createResourceOperationError(errorDeletingDeploymentProcess, deploymentProcess.ID, err)
	}

	d.SetId(constEmptyString)
	return nil
}
