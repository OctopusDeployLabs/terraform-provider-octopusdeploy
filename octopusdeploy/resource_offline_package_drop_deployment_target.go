package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOfflinePackageDropDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOfflinePackageDropDeploymentTargetCreate,
		DeleteContext: resourceOfflinePackageDropDeploymentTargetDelete,
		Importer:      getImporter(),
		ReadContext:   resourceOfflinePackageDropDeploymentTargetRead,
		Schema:        getOfflinePackageDropDeploymentTargetSchema(),
		UpdateContext: resourceOfflinePackageDropDeploymentTargetUpdate,
	}
}

func resourceOfflinePackageDropDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandOfflinePackageDropDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())
	return resourceOfflinePackageDropDeploymentTargetRead(ctx, d, m)
}

func resourceOfflinePackageDropDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Machines.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceOfflinePackageDropDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	setOfflinePackageDropDeploymentTarget(ctx, d, deploymentTarget)
	return nil
}

func resourceOfflinePackageDropDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandOfflinePackageDropDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOfflinePackageDropDeploymentTargetRead(ctx, d, m)
}
