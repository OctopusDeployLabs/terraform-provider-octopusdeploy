package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzureWebAppDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureWebAppDeploymentTargetCreate,
		DeleteContext: resourceAzureWebAppDeploymentTargetDelete,
		Importer:      getImporter(),
		ReadContext:   resourceAzureWebAppDeploymentTargetRead,
		Schema:        getAzureWebAppDeploymentTargetSchema(),
		UpdateContext: resourceAzureWebAppDeploymentTargetUpdate,
	}
}

func resourceAzureWebAppDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandAzureWebAppDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())
	return resourceAzureWebAppDeploymentTargetRead(ctx, d, m)
}

func resourceAzureWebAppDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Machines.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceAzureWebAppDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	setAzureWebAppDeploymentTarget(ctx, d, deploymentTarget)
	return nil
}

func resourceAzureWebAppDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandAzureWebAppDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAzureWebAppDeploymentTargetRead(ctx, d, m)
}
