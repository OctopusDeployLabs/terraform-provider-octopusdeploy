package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceListeningTentacleDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceListeningTentacleDeploymentTargetCreate,
		DeleteContext: resourceListeningTentacleDeploymentTargetDelete,
		Importer:      getImporter(),
		ReadContext:   resourceListeningTentacleDeploymentTargetRead,
		Schema:        getListeningTentacleDeploymentTargetSchema(),
		UpdateContext: resourceListeningTentacleDeploymentTargetUpdate,
	}
}

func resourceListeningTentacleDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandListeningTentacleDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())
	return resourceListeningTentacleDeploymentTargetRead(ctx, d, m)
}

func resourceListeningTentacleDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Machines.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceListeningTentacleDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	setListeningTentacleDeploymentTarget(ctx, d, deploymentTarget)
	return nil
}

func resourceListeningTentacleDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandListeningTentacleDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceListeningTentacleDeploymentTargetRead(ctx, d, m)
}
