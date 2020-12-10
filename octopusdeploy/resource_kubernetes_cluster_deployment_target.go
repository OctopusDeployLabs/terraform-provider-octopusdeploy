package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKubernetesClusterDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesClusterDeploymentTargetCreate,
		DeleteContext: resourceKubernetesClusterDeploymentTargetDelete,
		Description:   "This resource manages Kubernets cluster deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceKubernetesClusterDeploymentTargetRead,
		Schema:        getKubernetesClusterDeploymentTargetSchema(),
		UpdateContext: resourceKubernetesClusterDeploymentTargetUpdate,
	}
}

func resourceKubernetesClusterDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandKubernetesClusterDeploymentTarget(d)

	log.Printf("[INFO] creating Kubernetes cluster deployment target: %#v", deploymentTarget)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setKubernetesClusterDeploymentTarget(ctx, d, createdDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())

	log.Printf("[INFO] Kubernetes cluster deployment target created (%s)", d.Id())
	return nil
}

func resourceKubernetesClusterDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Kubernetes cluster deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Machines.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] Kubernetes cluster deployment target deleted")
	return nil
}

func resourceKubernetesClusterDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading Kubernetes cluster deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	deploymentTarget, err := client.Machines.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			log.Printf("[INFO] Kubernetes cluster deployment target (%s) not found; deleting from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := setKubernetesClusterDeploymentTarget(ctx, d, deploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Kubernetes cluster deployment target read (%s)", d.Id())
	return nil
}

func resourceKubernetesClusterDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating Kubernetes cluster deployment target (%s)", d.Id())

	deploymentTarget := expandKubernetesClusterDeploymentTarget(d)
	client := m.(*octopusdeploy.Client)
	updatedDeploymentTarget, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setKubernetesClusterDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Kubernetes cluster deployment target updated (%s)", d.Id())
	return nil
}
