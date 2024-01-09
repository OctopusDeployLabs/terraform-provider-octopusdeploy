package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbookprocess"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandRunbookProcess(ctx context.Context, d *schema.ResourceData, client *client.Client) *runbookprocess.RunbookProcess {
	runbookProcess := runbookprocess.NewRunbookProcess()
	runbookProcess.ID = d.Id()

	if v, ok := d.GetOk("project_id"); ok {
		runbookProcess.ProjectID = v.(string)
	}

	if v, ok := d.GetOk("runbook_id"); ok {
		runbookProcess.RunbookID = v.(string)
	}

	if v, ok := d.GetOk("last_snapshot_id"); ok {
		runbookProcess.LastSnapshotID = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		runbookProcess.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("version"); ok {
		version := int32(v.(int))
		runbookProcess.Version = &version
	}

	if v, ok := d.GetOk("step"); ok {
		steps := v.([]interface{})
		for _, step := range steps {
			deploymentStep := expandDeploymentStep(ctx, step.(map[string]interface{}))
			runbookProcess.Steps = append(runbookProcess.Steps, deploymentStep)
		}
	}

	return runbookProcess
}

func setRunbookProcess(ctx context.Context, d *schema.ResourceData, RunbookProcess *runbookprocess.RunbookProcess) error {
	d.Set("last_snapshot_id", RunbookProcess.LastSnapshotID)
	d.Set("project_id", RunbookProcess.ProjectID)
	d.Set("runbook_id", RunbookProcess.RunbookID)
	d.Set("space_id", RunbookProcess.SpaceID)
	d.Set("version", RunbookProcess.Version)

	if err := d.Set("step", flattenDeploymentSteps(RunbookProcess.Steps)); err != nil {
		return fmt.Errorf("error setting step: %s", err)
	}

	return nil
}
