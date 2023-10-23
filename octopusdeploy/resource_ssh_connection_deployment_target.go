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

func resourceSSHConnectionDeploymentTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSSHConnectionDeploymentTargetCreate,
		DeleteContext: resourceSSHConnectionDeploymentTargetDelete,
		Description:   "This resource manages SSH connection deployment targets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceSSHConnectionDeploymentTargetRead,
		Schema:        getSSHConnectionDeploymentTargetSchema(),
		UpdateContext: resourceSSHConnectionDeploymentTargetUpdate,
	}
}

func resourceSSHConnectionDeploymentTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	deploymentTarget := expandSSHConnectionDeploymentTarget(d)

	log.Printf("[INFO] creating SSH connection deployment target: %#v", deploymentTarget)

	client := m.(*client.Client)
	createdDeploymentTarget, err := machines.Add(client, deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setSSHConnectionDeploymentTarget(ctx, d, createdDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDeploymentTarget.GetID())

	log.Printf("[INFO] SSH connection deployment target created (%s)", d.Id())
	return nil
}

func resourceSSHConnectionDeploymentTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting SSH connection deployment target (%s)", d.Id())

	client := m.(*client.Client)
	if err := machines.DeleteByID(client, d.Get("space_id").(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] SSH connection deployment target deleted")
	return nil
}

func resourceSSHConnectionDeploymentTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading SSH connection deployment target (%s)", d.Id())

	client := m.(*client.Client)
	deploymentTarget, err := machines.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "SSH connection deployment target")
	}

	if err := setSSHConnectionDeploymentTarget(ctx, d, deploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] SSH connection deployment target read (%s)", d.Id())
	return nil
}

func resourceSSHConnectionDeploymentTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating SSH connection deployment target (%s)", d.Id())

	deploymentTarget := expandSSHConnectionDeploymentTarget(d)
	client := m.(*client.Client)
	updatedDeploymentTarget, err := machines.Update(client, deploymentTarget)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setSSHConnectionDeploymentTarget(ctx, d, updatedDeploymentTarget); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] SSH connection deployment target updated (%s)", d.Id())
	return nil
}
