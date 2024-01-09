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

	client := m.(*client.Client)
	createdDeploymentTarget, err := machines.Add(client, deploymentTarget)
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

	client := m.(*client.Client)
	if err := machines.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] Azure service fabric cluster deployment target deleted")
	return nil
}

func resourceAzureServiceFabricClusterDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading Azure service fabric cluster deployment target (%s)", d.Id())

	client := m.(*client.Client)
	deploymentTarget, err := machines.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Azure service fabric cluster deployment target")
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
	client := m.(*client.Client)
	updatedDeploymentTarget, err := machines.Update(client, deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureServiceFabricClusterDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure service fabric cluster deployment target updated (%s)", d.Id())
	return nil
}
