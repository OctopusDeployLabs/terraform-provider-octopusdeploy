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

func resourceOfflinePackageDropDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOfflinePackageDropDeploymentTargetCreate,
		DeleteContext: resourceOfflinePackageDropDeploymentTargetDelete,
		Description:   "This resource manages offline package drop deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceOfflinePackageDropDeploymentTargetRead,
		Schema:        getOfflinePackageDropDeploymentTargetSchema(),
		UpdateContext: resourceOfflinePackageDropDeploymentTargetUpdate,
	}
}

func resourceOfflinePackageDropDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandOfflinePackageDropDeploymentTarget(d)

	log.Printf("[INFO] creating offline package drop deployment target: %#v", deploymentTarget)

	client := m.(*client.Client)
	createdDeploymentTarget, err := machines.Add(client, deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setOfflinePackageDropDeploymentTarget(ctx, d, createdDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())

	log.Printf("[INFO] offline package drop deployment target created (%s)", d.Id())
	return nil
}

func resourceOfflinePackageDropDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting offline package drop deployment target (%s)", d.Id())

	client := m.(*client.Client)
	if err := machines.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] offline package drop deployment target deleted")
	return nil
}

func resourceOfflinePackageDropDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading offline package drop deployment target (%s)", d.Id())

	client := m.(*client.Client)
	deploymentTarget, err := machines.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "offline package drop deployment target")
	}

	if err := setOfflinePackageDropDeploymentTarget(ctx, d, deploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] offline package drop deployment target read (%s)", d.Id())
	return nil
}

func resourceOfflinePackageDropDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating offline package drop deployment target (%s)", d.Id())

	deploymentTarget := expandOfflinePackageDropDeploymentTarget(d)
	client := m.(*client.Client)
	updatedDeploymentTarget, err := machines.Update(client, deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setOfflinePackageDropDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] offline package drop deployment target updated (%s)", d.Id())
	return nil
}
