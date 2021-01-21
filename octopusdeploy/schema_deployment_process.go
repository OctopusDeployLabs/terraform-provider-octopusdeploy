package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandDeploymentProcess(d *schema.ResourceData) *octopusdeploy.DeploymentProcess {
	deploymentProcess := octopusdeploy.NewDeploymentProcess(d.Get("project_id").(string))
	deploymentProcess.ID = d.Id()

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
		"last_snapshot_id": deploymentProcess.LastSnapshotID,
		"project_id":       deploymentProcess.ProjectID,
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
			Required: true,
			Type:     schema.TypeString,
		},
		"step": getDeploymentStepSchema(),
		"version": {
			Optional: true,
			Type:     schema.TypeInt,
		},
	}
}
