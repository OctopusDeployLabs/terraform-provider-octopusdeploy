package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzureServiceFabricClusterDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureServiceFabricClusterDeploymentTargetCreate,
		DeleteContext: resourceAzureServiceFabricClusterDeploymentTargetDelete,
		Description:   "This resource manages Azure service fabric cluster deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceAzureServiceFabricClusterDeploymentTargetRead,
		Schema:        getAzureServiceFabricClusterDeploymentTargetSchema(),
		UpdateContext: resourceAzureServiceFabricClusterDeploymentTargetUpdate,
	}
}

func resourceAzureServiceFabricClusterDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandAzureServiceFabricClusterDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())
	return resourceAzureServiceFabricClusterDeploymentTargetRead(ctx, d, m)
}

func resourceAzureServiceFabricClusterDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Machines.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceAzureServiceFabricClusterDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	setAzureServiceFabricClusterDeploymentTarget(ctx, d, deploymentTarget)
	return nil
}

func resourceAzureServiceFabricClusterDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandAzureServiceFabricClusterDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAzureServiceFabricClusterDeploymentTargetRead(ctx, d, m)
}
