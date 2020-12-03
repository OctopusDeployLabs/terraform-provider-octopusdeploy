package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudRegionDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudRegionDeploymentTargetCreate,
		DeleteContext: resourceCloudRegionDeploymentTargetDelete,
		Importer:      getImporter(),
		ReadContext:   resourceCloudRegionDeploymentTargetRead,
		Schema:        getCloudRegionDeploymentTargetSchema(),
		UpdateContext: resourceCloudRegionDeploymentTargetUpdate,
	}
}

func resourceCloudRegionDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandCloudRegionDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())
	return resourceCloudRegionDeploymentTargetRead(ctx, d, m)
}

func resourceCloudRegionDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Machines.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceCloudRegionDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	setCloudRegionDeploymentTarget(ctx, d, deploymentTarget)
	return nil
}

func resourceCloudRegionDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandCloudRegionDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCloudRegionDeploymentTargetRead(ctx, d, m)
}
