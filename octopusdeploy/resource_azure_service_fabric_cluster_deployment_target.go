package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzureServiceFabricClusterDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureServiceFabricClusterDeploymentTargetCreate,
		DeleteContext: resourceAzureServiceFabricClusterDeploymentTargetDelete,
		Description:   "This resource manages Azure service fabric cluster deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceAzureServiceFabricClusterDeploymentTargetRead,
		Schema:        getAzureServiceFabricClusterDeploymentTargetSchema(),
		UpdateContext: resourceAzureServiceFabricClusterDeploymentTargetUpdate,
	}
}

func resourceAzureServiceFabricClusterDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandAzureServiceFabricClusterDeploymentTarget(d)

	log.Printf("[INFO] creating Azure service fabric cluster deployment target: %#v", deploymentTarget)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureServiceFabricClusterDeploymentTarget(ctx, d, createdDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())

	log.Printf("[INFO] Azure service fabric cluster deployment target created (%s)", d.Id())
	return nil
}

func resourceAzureServiceFabricClusterDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Azure service fabric cluster deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Machines.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] Azure service fabric cluster deployment target deleted")
	return nil
}

func resourceAzureServiceFabricClusterDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading Azure service fabric cluster deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	deploymentTarget, err := client.Machines.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] Azure service fabric cluster deployment target (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if err := setAzureServiceFabricClusterDeploymentTarget(ctx, d, deploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure service fabric cluster deployment target read (%s)", d.Id())
	return nil
}

func resourceAzureServiceFabricClusterDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating Azure service fabric cluster deployment target (%s)", d.Id())

	deploymentTarget := expandAzureServiceFabricClusterDeploymentTarget(d)
	client := m.(*octopusdeploy.Client)
	updatedDeploymentTarget, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureServiceFabricClusterDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure service fabric cluster deployment target updated (%s)", d.Id())
	return nil
}
