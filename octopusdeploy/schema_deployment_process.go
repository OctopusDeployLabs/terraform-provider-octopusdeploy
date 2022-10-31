package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandDeploymentProcess(d *schema.ResourceData, client *client.Client) *deployments.DeploymentProcess {
	projectID := d.Get("project_id").(string)
	deploymentProcess := deployments.NewDeploymentProcess(projectID)
	deploymentProcess.ID = d.Id()

	if v, ok := d.GetOk("branch"); ok {
		deploymentProcess.Branch = v.(string)
	} else {
		project, err := client.Projects.GetByID(projectID)
		if err != nil {
			return nil
		}

		if project.PersistenceSettings != nil && project.PersistenceSettings.GetType() == projects.PersistenceSettingsTypeVersionControlled {
			deploymentProcess.Branch = project.PersistenceSettings.(projects.GitPersistenceSettings).GetDefaultBranch()
		}
	}

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

func setDeploymentProcess(ctx context.Context, d *schema.ResourceData, deploymentProcess *deployments.DeploymentProcess) error {
	d.Set("branch", deploymentProcess.Branch)
	d.Set("last_snapshot_id", deploymentProcess.LastSnapshotID)
	d.Set("project_id", deploymentProcess.ProjectID)
	d.Set("space_id", deploymentProcess.SpaceID)
	d.Set("version", deploymentProcess.Version)

	if err := d.Set("step", flattenDeploymentSteps(deploymentProcess.Steps)); err != nil {
		return fmt.Errorf("error setting step: %s", err)
	}

	return nil
}
