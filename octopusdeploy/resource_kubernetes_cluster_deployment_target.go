package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/machines"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceKubernetesClusterDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKubernetesClusterDeploymentTargetCreate,
		DeleteContext: resourceKubernetesClusterDeploymentTargetDelete,
		Description:   "This resource manages Kubernetes cluster deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceKubernetesClusterDeploymentTargetRead,
		Schema:        getKubernetesClusterDeploymentTargetSchema(),
		UpdateContext: resourceKubernetesClusterDeploymentTargetUpdate,
	}
}

func resourceKubernetesClusterDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandKubernetesClusterDeploymentTarget(d)

	log.Printf("[INFO] creating Kubernetes cluster deployment target: %#v", deploymentTarget)

	client := m.(*client.Client)
	createdDeploymentTarget, err := machines.Add(client, deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setKubernetesClusterDeploymentTarget(ctx, d, createdDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())

	log.Printf("[INFO] Kubernetes cluster deployment target created (%s)", d.Id())
	return nil
}

func resourceKubernetesClusterDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Kubernetes cluster deployment target (%s)", d.Id())

	client := m.(*client.Client)
	if err := machines.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] Kubernetes cluster deployment target deleted")
	return nil
}

func resourceKubernetesClusterDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading Kubernetes cluster deployment target (%s)", d.Id())

	client := m.(*client.Client)
	deploymentTarget, err := machines.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Kubernetes cluster deployment target")
	}

	if err := setKubernetesClusterDeploymentTarget(ctx, d, deploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Kubernetes cluster deployment target read (%s)", d.Id())
	return nil
}

func resourceKubernetesClusterDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating Kubernetes cluster deployment target (%s)", d.Id())

	deploymentTarget := expandKubernetesClusterDeploymentTarget(d)
	client := m.(*client.Client)
	updatedDeploymentTarget, err := machines.Update(client, deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setKubernetesClusterDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Kubernetes cluster deployment target updated (%s)", d.Id())
	return nil
}
