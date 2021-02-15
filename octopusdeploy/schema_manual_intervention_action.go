package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenManualIntervention(actionMap map[string]interface{}, properties map[string]string) {
	for propertyName, propertyValue := range properties {
		switch propertyName {
		case "Octopus.Action.Manual.Instructions":
			actionMap["instructions"] = propertyValue
		case "Octopus.Action.Manual.ResponsibleTeamIds":
			actionMap["responsible_teams"] = propertyValue
		}
	}
}

func flattenManualInterventionAction(deploymentAction octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedWindowsServiceAction := map[string]interface{}{
		"can_be_used_for_project_versioning": deploymentAction.CanBeUsedForProjectVersioning,
		"channels":                           deploymentAction.Channels,
		"condition":                          deploymentAction.Condition,
		"container":                          flattenDeploymentActionContainer(deploymentAction.Container),
		"environments":                       deploymentAction.Environments,
		"excluded_environments":              deploymentAction.ExcludedEnvironments,
		"id":                                 deploymentAction.ID,
		"is_disabled":                        deploymentAction.IsDisabled,
		"is_required":                        deploymentAction.IsRequired,
		"name":                               deploymentAction.Name,
		"notes":                              deploymentAction.Notes,
		"properties":                         deploymentAction.Properties,
		"tenant_tags":                        deploymentAction.TenantTags,
	}

	flattenManualIntervention(flattenedWindowsServiceAction, deploymentAction.Properties)

	return flattenedWindowsServiceAction
}

func getManualInterventionActionSchema() *schema.Schema {
	actionSchema, element := getCommonDeploymentActionSchema()

	element.Schema["instructions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The instructions for the user to follow",
		Required:    true,
	}

	element.Schema["responsible_teams"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The teams responsible to resolve this step. If no teams are specified, all users who have permission to deploy the project can resolve it.",
		Optional:    true,
	}

	return actionSchema
}

func expandManualInterventionAction(tfAction map[string]interface{}) octopusdeploy.DeploymentAction {
	resource := expandDeploymentAction(tfAction)
	resource.ActionType = "Octopus.Manual"
	resource.Properties["Octopus.Action.Manual.Instructions"] = tfAction["instructions"].(string)

	responsibleTeams := tfAction["responsible_teams"]
	if responsibleTeams != nil {
		resource.Properties["Octopus.Action.Manual.ResponsibleTeamIds"] = responsibleTeams.(string)
	}

	return resource
}
