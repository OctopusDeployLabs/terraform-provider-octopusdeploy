package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandDeploymentProcess(d *schema.ResourceData) *octopusdeploy.DeploymentProcess {
	deploymentProcess := octopusdeploy.NewDeploymentProcess(d.Get("project_id").(string))
	deploymentProcess.ID = d.Id()

	if v, ok := d.GetOk("last_snapshot_id"); ok {
		deploymentProcess.LastSnapshotID = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		deploymentProcess.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("version"); ok {
		deploymentProcess.Version = int32(v.(int))
	}

	if v, ok := d.GetOk("step"); ok {
		steps := v.([]interface{})
		for _, step := range steps {
			deploymentStep := expandDeploymentStep(step.(map[string]interface{}))
			deploymentProcess.Steps = append(deploymentProcess.Steps, deploymentStep)
		}
	}

	return deploymentProcess
}

func flattenDeploymentProcess(deploymentProcess *octopusdeploy.DeploymentProcess) []interface{} {
	if deploymentProcess == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"id":               deploymentProcess.ID,
		"last_snapshot_id": deploymentProcess.LastSnapshotID,
		"project_id":       deploymentProcess.ProjectID,
		"space_id":         deploymentProcess.SpaceID,
		"step":             flattenDeploymentSteps(deploymentProcess.Steps),
		"version":          deploymentProcess.Version,
	}}
}

func getDeploymentProcessSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": getIDSchema(),
		"last_snapshot_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"project_id": {
			Description: "The project ID associated with this d eployment process.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"space_id": getSpaceIDSchema(),
		"step":     getDeploymentStepSchema(),
		"version": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeInt,
		},
	}
}

func setDeploymentProcess(ctx context.Context, d *schema.ResourceData, deploymentProcess *octopusdeploy.DeploymentProcess) error {
	d.Set("last_snapshot_id", deploymentProcess.LastSnapshotID)
	d.Set("project_id", deploymentProcess.ProjectID)
	d.Set("space_id", deploymentProcess.SpaceID)
	d.Set("version", deploymentProcess.Version)

	if err := d.Set("step", flattenDeploymentSteps(deploymentProcess.Steps)); err != nil {
		return fmt.Errorf("error setting step: %s", err)
	}

	d.SetId(deploymentProcess.GetID())

	return nil
}
