package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenManualIntervention(actionMap map[string]interface{}, properties map[string]octopusdeploy.PropertyValue) {
	for propertyName, propertyValue := range properties {
		switch propertyName {
		case "Octopus.Action.Manual.Instructions":
			actionMap["instructions"] = propertyValue.Value
		case "Octopus.Action.Manual.ResponsibleTeamIds":
			actionMap["responsible_teams"] = propertyValue.Value
		}
	}
}

func flattenManualInterventionAction(action *octopusdeploy.DeploymentAction) map[string]interface{} {
	flattenedAction := flattenAction(action)
	flattenManualIntervention(flattenedAction, action.Properties)

	return flattenedAction
}

func getManualInterventionActionSchema() *schema.Schema {
	actionSchema, element := getActionSchema()

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

func expandManualInterventionAction(tfAction map[string]interface{}) *octopusdeploy.DeploymentAction {
	resource := expandAction(tfAction)
	resource.ActionType = "Octopus.Manual"
	resource.Properties["Octopus.Action.Manual.Instructions"] = octopusdeploy.NewPropertyValue(tfAction["instructions"].(string), false)

	responsibleTeams := tfAction["responsible_teams"]
	if responsibleTeams != nil {
		resource.Properties["Octopus.Action.Manual.ResponsibleTeamIds"] = octopusdeploy.NewPropertyValue(responsibleTeams.(string), false)
	}

	return resource
}
