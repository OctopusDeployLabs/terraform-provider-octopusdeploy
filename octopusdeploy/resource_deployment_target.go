package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeploymentTargetCreate,
		DeleteContext: resourceDeploymentTargetDelete,
		Importer:      getImporter(),
		ReadContext:   resourceDeploymentTargetRead,
		Schema:        getDeploymentTargetSchema(),
		UpdateContext: resourceDeploymentTargetUpdate,
	}
}

func resourceDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandDeploymentTarget(d)
	deploymentTarget.Status = "Unknown"

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	setDeploymentTarget(ctx, d, createdDeploymentTarget)
	return nil
}

func resourceDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Machines.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	deploymentTarget, err := client.Machines.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	setDeploymentTarget(ctx, d, deploymentTarget)
	return nil
}

func resourceDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	updatedDeploymentTarget, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	setDeploymentTarget(ctx, d, updatedDeploymentTarget)
	return nil
}
