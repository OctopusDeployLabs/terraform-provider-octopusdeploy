package octopusdeploy

import (
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
			constStep: getDeploymentStepSchema(),
		},
	}
}

func resourceDeploymentProcessCreate(d *schema.ResourceData, m interface{}) error {
	deploymentProcess := buildDeploymentProcessResource(d)

	apiClient := m.(*client.Client)
	project, err := apiClient.Projects.GetByID(deploymentProcess.ProjectID)
	if err != nil {
		return createResourceOperationError(errorReadingProject, project.Name, err)
	}

	current, err := apiClient.DeploymentProcesses.GetByID(project.DeploymentProcessID)
	if err != nil {
		return createResourceOperationError(errorReadingDeploymentProcess, project.DeploymentProcessID, err)
	}

	deploymentProcess.ID = current.ID
	deploymentProcess.Version = current.Version

	resource, err := apiClient.DeploymentProcesses.Update(*deploymentProcess)
	if err != nil {
		return createResourceOperationError(errorCreatingDeploymentProcess, deploymentProcess.ID, err)
	}

	if isEmpty(resource.ID) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.ID)
	}

	return nil
}

func buildDeploymentProcessResource(d *schema.ResourceData) *model.DeploymentProcess {
	deploymentProcess := &model.DeploymentProcess{
		ProjectID: d.Get(constProjectID).(string),
	}

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

	apiClient := m.(*client.Client)
	resource, err := apiClient.DeploymentProcesses.GetByID(id)
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
	deploymentProcess.ID = d.Id() // set deploymentProcess struct ID so octopus knows which deploymentProcess to update

	apiClient := m.(*client.Client)

	current, err := apiClient.DeploymentProcesses.GetByID(deploymentProcess.ID)
	if err != nil {
		return createResourceOperationError(errorReadingDeploymentProcess, deploymentProcess.ID, err)
	}

	deploymentProcess.Version = current.Version
	deploymentProcess, err = apiClient.DeploymentProcesses.Update(*deploymentProcess)
	if err != nil {
		return createResourceOperationError(errorUpdatingDeploymentProcess, d.Id(), err)
	}

	d.SetId(deploymentProcess.ID)

	return nil
}

func resourceDeploymentProcessDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	current, err := apiClient.DeploymentProcesses.GetByID(d.Id())
	if err != nil {
		return createResourceOperationError(errorReadingDeploymentProcess, d.Id(), err)
	}

	deploymentProcess := &model.DeploymentProcess{
		Version: current.Version,
	}
	deploymentProcess.ID = d.Id()

	deploymentProcess, err = apiClient.DeploymentProcesses.Update(*deploymentProcess)
	if err != nil {
		return createResourceOperationError(errorDeletingDeploymentProcess, deploymentProcess.ID, err)
	}

	d.SetId(constEmptyString)
	return nil
}
