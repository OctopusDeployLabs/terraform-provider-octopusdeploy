package octopusdeploy

import (
	"context"
	"log"

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

	log.Printf("[INFO] creating deployment target: %#v", deploymentTarget)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setDeploymentTarget(ctx, d, createdDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] deployment target created (%s)", d.Id())
	return nil
}

func resourceDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	err := client.Machines.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] deployment target deleted")
	return nil
}

func resourceDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	deploymentTarget, err := client.Machines.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	setDeploymentTarget(ctx, d, deploymentTarget)

	log.Printf("[INFO] deployment target read (%s)", d.Id())
	return nil
}

func resourceDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating deployment target (%s)", d.Id())

	deploymentTarget := expandDeploymentTarget(d)

	client := m.(*octopusdeploy.Client)
	updatedDeploymentTarget, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	setDeploymentTarget(ctx, d, updatedDeploymentTarget)

	log.Printf("[INFO] deployment target updated (%s)", d.Id())
	return nil
}
