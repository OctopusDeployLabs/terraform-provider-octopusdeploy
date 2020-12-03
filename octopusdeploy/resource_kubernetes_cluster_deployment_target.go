package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKubernetesClusterDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesClusterDeploymentTargetCreate,
		DeleteContext: resourceKubernetesClusterDeploymentTargetDelete,
		Importer:      getImporter(),
		ReadContext:   resourceKubernetesClusterDeploymentTargetRead,
		Schema:        getKubernetesClusterDeploymentTargetSchema(),
		UpdateContext: resourceKubernetesClusterDeploymentTargetUpdate,
	}
}

func resourceKubernetesClusterDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandKubernetesClusterDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())
	return resourceKubernetesClusterDeploymentTargetRead(ctx, d, m)
}

func resourceKubernetesClusterDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Machines.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceKubernetesClusterDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	deploymentTarget, err := client.Machines.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	setKubernetesClusterDeploymentTarget(ctx, d, deploymentTarget)
	return nil
}

func resourceKubernetesClusterDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandKubernetesClusterDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceKubernetesClusterDeploymentTargetRead(ctx, d, m)
}
