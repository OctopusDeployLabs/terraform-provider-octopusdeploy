package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePollingTentacleDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePollingTentacleDeploymentTargetCreate,
		DeleteContext: resourcePollingTentacleDeploymentTargetDelete,
		Description:   "This resource manages polling tentacle deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourcePollingTentacleDeploymentTargetRead,
		Schema:        getPollingTentacleDeploymentTargetSchema(),
		UpdateContext: resourcePollingTentacleDeploymentTargetUpdate,
	}
}

func resourcePollingTentacleDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandPollingTentacleDeploymentTarget(d)

	log.Printf("[INFO] creating polling tentacle deployment target: %#v", deploymentTarget)

	client := m.(*octopusdeploy.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setPollingTentacleDeploymentTarget(ctx, d, createdDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())

	log.Printf("[INFO] polling tentacle deployment target created (%s)", d.Id())
	return nil
}

func resourcePollingTentacleDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting polling tentacle deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.Machines.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] polling tentacle deployment target deleted")
	return nil
}

func resourcePollingTentacleDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading polling tentacle deployment target (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	deploymentTarget, err := client.Machines.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			log.Printf("[INFO] polling tentacle deployment target (%s) not found; deleting from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := setPollingTentacleDeploymentTarget(ctx, d, deploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] polling tentacle deployment target read (%s)", d.Id())
	return nil
}

func resourcePollingTentacleDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating polling tentacle deployment target (%s)", d.Id())

	deploymentTarget := expandPollingTentacleDeploymentTarget(d)
	client := m.(*octopusdeploy.Client)
	updatedDeploymentTarget, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setPollingTentacleDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] polling tentacle deployment target updated (%s)", d.Id())
	return nil
}
