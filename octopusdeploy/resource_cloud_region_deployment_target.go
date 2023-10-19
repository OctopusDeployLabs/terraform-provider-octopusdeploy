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

	client := m.(*client.Client)
	createdDeploymentTarget, err := machines.Add(client, deploymentTarget)
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

	client := m.(*client.Client)
	if err := machines.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] cloud region deployment target deleted")
	return nil
}

func resourceCloudRegionDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading cloud region deployment target (%s)", d.Id())

	client := m.(*client.Client)
	deploymentTarget, err := machines.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "cloud region deployment target")
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
	client := m.(*client.Client)
	updatedDeploymentTarget, err := machines.Update(client, deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setCloudRegionDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] cloud region deployment target updated (%s)", d.Id())
	return nil
}
