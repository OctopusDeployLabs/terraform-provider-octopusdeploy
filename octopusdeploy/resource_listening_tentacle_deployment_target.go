package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceListeningTentacleDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceListeningTentacleDeploymentTargetCreate,
		DeleteContext: resourceListeningTentacleDeploymentTargetDelete,
		Description:   "This resource manages listening tentacle deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceListeningTentacleDeploymentTargetRead,
		Schema:        getListeningTentacleDeploymentTargetSchema(),
		UpdateContext: resourceListeningTentacleDeploymentTargetUpdate,
	}
}

func resourceListeningTentacleDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandListeningTentacleDeploymentTarget(d)

	log.Printf("[INFO] creating listening tentacle deployment target: %#v", deploymentTarget)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setListeningTentacleDeploymentTarget(ctx, d, createdDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())

	log.Printf("[INFO] listening tentacle deployment target created (%s)", d.Id())
	return nil
}

func resourceListeningTentacleDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting listening tentacle deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Machines.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] listening tentacle deployment target deleted")
	return nil
}

func resourceListeningTentacleDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading listening tentacle deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	deploymentTarget, err := client.Machines.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] listening tentacle deployment target (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if err := setListeningTentacleDeploymentTarget(ctx, d, deploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] listening tentacle deployment target read (%s)", d.Id())
	return nil
}

func resourceListeningTentacleDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating listening tentacle deployment target (%s)", d.Id())

	deploymentTarget := expandListeningTentacleDeploymentTarget(d)
	client := m.(*octopusdeploy.Client)
	updatedDeploymentTarget, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setListeningTentacleDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] listening tentacle deployment target updated (%s)", d.Id())
	return nil
}
