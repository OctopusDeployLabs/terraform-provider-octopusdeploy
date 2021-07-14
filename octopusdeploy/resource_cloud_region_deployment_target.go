package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudRegionDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudRegionDeploymentTargetCreate,
		DeleteContext: resourceCloudRegionDeploymentTargetDelete,
		Description:   "This resource manages cloud region deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceCloudRegionDeploymentTargetRead,
		Schema:        getCloudRegionDeploymentTargetSchema(),
		UpdateContext: resourceCloudRegionDeploymentTargetUpdate,
	}
}

func resourceCloudRegionDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandCloudRegionDeploymentTarget(d)

	log.Printf("[INFO] creating cloud region deployment target: %#v", deploymentTarget)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setCloudRegionDeploymentTarget(ctx, d, createdDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())

	log.Printf("[INFO] cloud region deployment target created (%s)", d.Id())
	return nil
}

func resourceCloudRegionDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting cloud region deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Machines.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] cloud region deployment target deleted")
	return nil
}

func resourceCloudRegionDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading cloud region deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	deploymentTarget, err := client.Machines.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] cloud region deployment target (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if err := setCloudRegionDeploymentTarget(ctx, d, deploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] cloud region deployment target read (%s)", d.Id())
	return nil
}

func resourceCloudRegionDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating cloud region deployment target (%s)", d.Id())

	deploymentTarget := expandCloudRegionDeploymentTarget(d)
	client := m.(*octopusdeploy.Client)
	updatedDeploymentTarget, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setCloudRegionDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] cloud region deployment target updated (%s)", d.Id())
	return nil
}
