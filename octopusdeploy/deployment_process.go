package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDeploymentProcess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeploymentProcessCreate,
		ReadContext:   resourceDeploymentProcessRead,
		UpdateContext: resourceDeploymentProcessUpdate,
		DeleteContext: resourceDeploymentProcessDelete,

		Schema: map[string]*schema.Schema{
			constProjectID: {
				Type:     schema.TypeString,
				Required: true,
			},
			constStep: getDeploymentStepSchema(),
		},
	}
}

func resourceDeploymentProcessCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentProcess := buildDeploymentProcessResource(d)
	diagValidate()

	apiClient := m.(*client.Client)
	project, err := apiClient.Projects.GetByID(deploymentProcess.ProjectID)
	if err != nil {
		return diag.FromErr(createResourceOperationError(errorReadingProject, project.Name, err))
	}

	current, err := apiClient.DeploymentProcesses.GetByID(project.DeploymentProcessID)
	if err != nil {
		return diag.FromErr(createResourceOperationError(errorReadingDeploymentProcess, project.DeploymentProcessID, err))
	}

	deploymentProcess.ID = current.ID
	deploymentProcess.Version = current.Version

	resource, err := apiClient.DeploymentProcesses.Update(*deploymentProcess)
	if err != nil {
		return diag.FromErr(createResourceOperationError(errorCreatingDeploymentProcess, deploymentProcess.ID, err))
	}

	if isEmpty(resource.ID) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.ID)
	}

	return nil
}

func buildDeploymentProcessResource(d *schema.ResourceData) *model.DeploymentProcess {
	deploymentProcess, err := model.NewDeploymentProcess(d.Get(constProjectID).(string))
	if err != nil {
		return nil
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

func resourceDeploymentProcessRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	diagValidate()

	apiClient := m.(*client.Client)
	resource, err := apiClient.DeploymentProcesses.GetByID(id)
	if err != nil {
		return diag.FromErr(createResourceOperationError(errorReadingDeploymentProcess, id, err))
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constDeploymentProcess, m)

	return nil
}

func resourceDeploymentProcessUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentProcess := buildDeploymentProcessResource(d)
	deploymentProcess.ID = d.Id() // set ID so Octopus API knows which deployment process to update

	diagValidate()

	apiClient := m.(*client.Client)
	current, err := apiClient.DeploymentProcesses.GetByID(deploymentProcess.ID)
	if err != nil {
		return diag.FromErr(createResourceOperationError(errorReadingDeploymentProcess, deploymentProcess.ID, err))
	}

	deploymentProcess.Version = current.Version
	resource, err := apiClient.DeploymentProcesses.Update(*deploymentProcess)
	if err != nil {
		return diag.FromErr(createResourceOperationError(errorUpdatingDeploymentProcess, d.Id(), err))
	}

	d.SetId(resource.ID)

	return nil
}

func resourceDeploymentProcessDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	diagValidate()

	current, err := apiClient.DeploymentProcesses.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(createResourceOperationError(errorReadingDeploymentProcess, d.Id(), err))
	}

	deploymentProcess := &model.DeploymentProcess{
		Version: current.Version,
	}
	deploymentProcess.ID = d.Id()

	deploymentProcess, err = apiClient.DeploymentProcesses.Update(*deploymentProcess)
	if err != nil {
		return diag.FromErr(createResourceOperationError(errorDeletingDeploymentProcess, deploymentProcess.ID, err))
	}

	d.SetId(constEmptyString)
	return nil
}
