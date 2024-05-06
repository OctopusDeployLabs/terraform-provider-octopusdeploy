package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKubernetesAgentDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesAgentDeploymentTargetCreate,
		DeleteContext: resourceKubernetesAgentDeploymentTargetDelete,
		Description:   "This resource manages Kubernetes agent deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceKubernetesAgentDeploymentTargetRead,
		Schema:        getKubernetesAgentDeploymentTargetSchema(),
		UpdateContext: resourceKubernetesAgentDeploymentTargetUpdate,
	}
}

func resourceKubernetesAgentDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandKubernetesAgentDeploymentTarget(d)
	client := m.(*client.Client)
	createdDeploymentTarget, err := machines.Add(client, deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())
	return nil
}

func resourceKubernetesAgentDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	deploymentTarget, err := machines.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "kubernetes tentacle deployment target")
	}

	flattenedKubernetesAgentDeploymentTarget := flattenKubernetesAgentDeploymentTarget(deploymentTarget)
	for key, value := range flattenedKubernetesAgentDeploymentTarget {
		if key != "id" {
			err := d.Set(key, value)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func resourceKubernetesAgentDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	if err := machines.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func resourceKubernetesAgentDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandKubernetesAgentDeploymentTarget(d)
	client := m.(*client.Client)

	deploymentTarget.ID = d.Id()

	updatedDeploymentTarget, err := machines.Update(client, deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(updatedDeploymentTarget.GetID())

	return nil
}
