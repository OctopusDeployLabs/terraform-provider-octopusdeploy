package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDeploymentProcess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeploymentProcessCreate,
		DeleteContext: resourceDeploymentProcessDelete,
		Importer:      getImporter(),
		ReadContext:   resourceDeploymentProcessRead,
		Schema:        getDeploymentProcessSchema(),
		UpdateContext: resourceDeploymentProcessUpdate,
	}
}

func resourceDeploymentProcessCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentProcess := expandDeploymentProcess(d)

	client := m.(*octopusdeploy.Client)
	project, err := client.Projects.GetByID(deploymentProcess.ProjectID)
	if err != nil {
		diag.FromErr(err)
	}

	current, err := client.DeploymentProcesses.GetByID(project.DeploymentProcessID)
	if err != nil {
		diag.FromErr(err)
	}

	deploymentProcess.ID = current.ID
	deploymentProcess.Version = current.Version

	resource, err := client.DeploymentProcesses.Update(deploymentProcess)
	if err != nil {
		diag.FromErr(err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func resourceDeploymentProcessRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	resource, err := client.DeploymentProcesses.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if resource == nil {
		d.SetId("")
		return nil
	}

	logResource(constDeploymentProcess, m)

	return nil
}

func resourceDeploymentProcessUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentProcess := expandDeploymentProcess(d)

	client := m.(*octopusdeploy.Client)
	current, err := client.DeploymentProcesses.GetByID(deploymentProcess.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	deploymentProcess.Version = current.Version
	resource, err := client.DeploymentProcesses.Update(deploymentProcess)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceDeploymentProcessDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	current, err := client.DeploymentProcesses.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deploymentProcess := &octopusdeploy.DeploymentProcess{
		Version: current.Version,
	}
	deploymentProcess.ID = d.Id()

	deploymentProcess, err = client.DeploymentProcesses.Update(deploymentProcess)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
