package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSSHConnectionDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSHConnectionDeploymentTargetCreate,
		DeleteContext: resourceSSHConnectionDeploymentTargetDelete,
		Importer:      getImporter(),
		ReadContext:   resourceSSHConnectionDeploymentTargetRead,
		Schema:        getSSHConnectionDeploymentTargetSchema(),
		UpdateContext: resourceSSHConnectionDeploymentTargetUpdate,
	}
}

func resourceSSHConnectionDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandSSHConnectionDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())
	return resourceSSHConnectionDeploymentTargetRead(ctx, d, m)
}

func resourceSSHConnectionDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Machines.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceSSHConnectionDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	setSSHConnectionDeploymentTarget(ctx, d, deploymentTarget)
	return nil
}

func resourceSSHConnectionDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandSSHConnectionDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSSHConnectionDeploymentTargetRead(ctx, d, m)
}
