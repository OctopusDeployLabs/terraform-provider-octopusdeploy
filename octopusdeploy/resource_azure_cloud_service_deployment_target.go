package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzureCloudServiceDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureCloudServiceDeploymentTargetCreate,
		DeleteContext: resourceAzureCloudServiceDeploymentTargetDelete,
		Description:   "This resource manages Azure cloud service deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceAzureCloudServiceDeploymentTargetRead,
		Schema:        getAzureCloudServiceDeploymentTargetSchema(),
		UpdateContext: resourceAzureCloudServiceDeploymentTargetUpdate,
	}
}

func resourceAzureCloudServiceDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandAzureCloudServiceDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())
	return resourceAzureCloudServiceDeploymentTargetRead(ctx, d, m)
}

func resourceAzureCloudServiceDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Machines.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceAzureCloudServiceDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	setAzureCloudServiceDeploymentTarget(ctx, d, deploymentTarget)
	return nil
}

func resourceAzureCloudServiceDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandAzureCloudServiceDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAzureCloudServiceDeploymentTargetRead(ctx, d, m)
}
