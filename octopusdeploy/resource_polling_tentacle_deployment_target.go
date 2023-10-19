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

	client := m.(*client.Client)
	createdDeploymentTarget, err := machines.Add(client, deploymentTarget)
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

	client := m.(*client.Client)
	if err := machines.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] polling tentacle deployment target deleted")
	return nil
}

func resourcePollingTentacleDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading polling tentacle deployment target (%s)", d.Id())

	client := m.(*client.Client)
	deploymentTarget, err := machines.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "polling tentacle deployment target")
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
	client := m.(*client.Client)
	updatedDeploymentTarget, err := machines.Update(client, deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setPollingTentacleDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] polling tentacle deployment target updated (%s)", d.Id())
	return nil
}
