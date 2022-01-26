package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func expandDestroyTerraformTemplateAction(flattenedAction map[string]interface{}) *octopusdeploy.DeploymentAction {
	action := expandAction(flattenedAction)
	if isPlan, ok := flattenedAction["is_plan"].(bool); ok && isPlan {
		action.ActionType = "Octopus.TerraformPlanDestroy"
	} else {
		action.ActionType = "Octopus.TerraformDestroy"
	}
	expandTerraformTemplateAction(flattenedAction, action)

	return action
}
