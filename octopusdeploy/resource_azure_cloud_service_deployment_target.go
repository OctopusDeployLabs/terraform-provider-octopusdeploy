package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAzureCloudServiceDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureCloudServiceDeploymentTargetCreate,
		DeleteContext: resourceAzureCloudServiceDeploymentTargetDelete,
		Description:   "This resource manages Azure cloud service deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceAzureCloudServiceDeploymentTargetRead,
		Schema:        getAzureCloudServiceDeploymentTargetSchema(),
		UpdateContext: resourceAzureCloudServiceDeploymentTargetUpdate,
	}
}

func resourceAzureCloudServiceDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandAzureCloudServiceDeploymentTarget(d)

	log.Printf("[INFO] creating Azure cloud service deployment target: %#v", deploymentTarget)

	client := m.(*client.Client)
	createdDeploymentTarget, err := client.Machines.Add(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureCloudServiceDeploymentTarget(ctx, d, createdDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())

	log.Printf("[INFO] Azure cloud service deployment target created (%s)", d.Id())
	return nil
}

func resourceAzureCloudServiceDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Azure cloud service deployment target (%s)", d.Id())

	client := m.(*client.Client)
	if err := client.Machines.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] Azure cloud service deployment target deleted")
	return nil
}

func resourceAzureCloudServiceDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading Azure cloud service deployment target (%s)", d.Id())

	client := m.(*client.Client)
	deploymentTarget, err := client.Machines.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Azure cloud service deployment target")
	}

	if err := setAzureCloudServiceDeploymentTarget(ctx, d, deploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure cloud service deployment target read (%s)", d.Id())
	return nil
}

func resourceAzureCloudServiceDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating Azure cloud service deployment target (%s)", d.Id())

	deploymentTarget := expandAzureCloudServiceDeploymentTarget(d)
	client := m.(*client.Client)
	updatedDeploymentTarget, err := client.Machines.Update(deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureCloudServiceDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure cloud service deployment target updated (%s)", d.Id())
	return nil
}
