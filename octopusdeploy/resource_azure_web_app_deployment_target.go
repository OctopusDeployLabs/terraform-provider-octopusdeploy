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

func resourceAzureWebAppDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureWebAppDeploymentTargetCreate,
		DeleteContext: resourceAzureWebAppDeploymentTargetDelete,
		Description:   "This resource manages Azure web app deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceAzureWebAppDeploymentTargetRead,
		Schema:        getAzureWebAppDeploymentTargetSchema(),
		UpdateContext: resourceAzureWebAppDeploymentTargetUpdate,
	}
}

func resourceAzureWebAppDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandAzureWebAppDeploymentTarget(d)

	log.Printf("[INFO] creating Azure web app deployment target: %#v", deploymentTarget)

	client := m.(*client.Client)
	createdDeploymentTarget, err := machines.Add(client, deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureWebAppDeploymentTarget(ctx, d, createdDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())

	log.Printf("[INFO] Azure web app deployment target created (%s)", d.Id())
	return nil
}

func resourceAzureWebAppDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Azure web app deployment target (%s)", d.Id())

	client := m.(*client.Client)
	if err := machines.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] Azure web app deployment target deleted")
	return nil
}

func resourceAzureWebAppDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading Azure web app deployment target (%s)", d.Id())

	client := m.(*client.Client)
	deploymentTarget, err := machines.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Azure web app deployment target")
	}

	if err := setAzureWebAppDeploymentTarget(ctx, d, deploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure web app deployment target read (%s)", d.Id())
	return nil
}

func resourceAzureWebAppDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating Azure web app deployment target (%s)", d.Id())

	deploymentTarget := expandAzureWebAppDeploymentTarget(d)
	client := m.(*client.Client)
	updatedDeploymentTarget, err := machines.Update(client, deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAzureWebAppDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Azure web app deployment target updated (%s)", d.Id())
	return nil
}
