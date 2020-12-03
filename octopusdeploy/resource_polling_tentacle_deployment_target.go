package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePollingTentacleDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePollingTentacleDeploymentTargetCreate,
		DeleteContext: resourcePollingTentacleDeploymentTargetDelete,
		Importer:      getImporter(),
		ReadContext:   resourcePollingTentacleDeploymentTargetRead,
		Schema:        getPollingTentacleDeploymentTargetSchema(),
		UpdateContext: resourcePollingTentacleDeploymentTargetUpdate,
	}
}

func resourcePollingTentacleDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandPollingTentacleDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())
	return resourcePollingTentacleDeploymentTargetRead(ctx, d, m)
}

func resourcePollingTentacleDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Machines.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourcePollingTentacleDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	setPollingTentacleDeploymentTarget(ctx, d, deploymentTarget)
	return nil
}

func resourcePollingTentacleDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandPollingTentacleDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePollingTentacleDeploymentTargetRead(ctx, d, m)
}
