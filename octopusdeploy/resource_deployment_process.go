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
		Description:   "This resource manages deployment processes in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceDeploymentProcessRead,
		Schema:        getDeploymentProcessSchema(),
		UpdateContext: resourceDeploymentProcessUpdate,
	}
}

func resourceDeploymentProcessCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentProcess := expandDeploymentProcess(d)

	log.Printf("[INFO] creating deployment process: %#v", deploymentProcess)

	client := m.(*octopusdeploy.Client)
	project, err := client.Projects.GetByID(deploymentProcess.ProjectID)
	if err != nil {
		return diag.FromErr(err)
	}

	current, err := client.DeploymentProcesses.GetByID(project.DeploymentProcessID)
	if err != nil {
		return diag.FromErr(err)
	}

	deploymentProcess.ID = current.ID
	deploymentProcess.Version = current.Version

	resource, err := client.DeploymentProcesses.Update(deploymentProcess)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setDeploymentProcess(ctx, d, resource); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.GetID())

	log.Printf("[INFO] deployment process created (%s)", d.Id())
	return nil
}

func resourceDeploymentProcessDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting deployment process (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	current, err := client.DeploymentProcesses.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deploymentProcess := &octopusdeploy.DeploymentProcess{
		Version: current.Version,
	}
	deploymentProcess.ID = d.Id()

	_, err = client.DeploymentProcesses.Update(deploymentProcess)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] deployment process deleted")
	return nil
}

func resourceDeploymentProcessRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading deployment process (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	deploymentProcess, err := client.DeploymentProcesses.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] deployment process (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if err := setDeploymentProcess(ctx, d, deploymentProcess); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] deployment process read (%s)", d.Id())
	return nil
}

func resourceDeploymentProcessUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating deployment process (%s)", d.Id())

	deploymentProcess := expandDeploymentProcess(d)
	client := m.(*octopusdeploy.Client)
	current, err := client.DeploymentProcesses.GetByID(deploymentProcess.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	deploymentProcess.Version = current.Version
	updatedDeploymentProcess, err := client.DeploymentProcesses.Update(deploymentProcess)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setDeploymentProcess(ctx, d, updatedDeploymentProcess); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] deployment process updated (%s)", d.Id())
	return nil
}
