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

	client := m.(*client.Client)
	createdDeploymentTarget, err := machines.Add(client, deploymentTarget)
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

	client := m.(*client.Client)
	if err := machines.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] listening tentacle deployment target deleted")
	return nil
}

func resourceListeningTentacleDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading listening tentacle deployment target (%s)", d.Id())

	client := m.(*client.Client)
	deploymentTarget, err := machines.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "listening tentacle deployment target")
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
	client := m.(*client.Client)
	updatedDeploymentTarget, err := machines.Update(client, deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setListeningTentacleDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] listening tentacle deployment target updated (%s)", d.Id())
	return nil
}
