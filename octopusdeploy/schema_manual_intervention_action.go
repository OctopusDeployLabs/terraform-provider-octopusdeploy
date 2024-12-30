package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenManualIntervention(actionMap map[string]interface{}, properties map[string]core.PropertyValue) {
	for propertyName, propertyValue := range properties {
		switch propertyName {
		case "Octopus.Action.Manual.Instructions":
			actionMap["instructions"] = propertyValue.Value
		case "Octopus.Action.Manual.ResponsibleTeamIds":
			actionMap["responsible_teams"] = propertyValue.Value
		case "Octopus.Action.Manual.BlockConcurrentDeployments":
			actionMap["block_deployments"] = propertyValue.Value
		}
	}
}

func flattenManualInterventionAction(action *deployments.DeploymentAction) map[string]interface{} {
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

	element.Schema["block_deployments"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Should other deployments be blocked while this manual intervention is awaiting action.",
		Optional:    true,
	}

	return actionSchema
}

func expandManualInterventionAction(tfAction map[string]interface{}) *deployments.DeploymentAction {
	resource := expandAction(tfAction)
	resource.ActionType = "Octopus.Manual"
	resource.Properties["Octopus.Action.Manual.Instructions"] = core.NewPropertyValue(tfAction["instructions"].(string), false)

	responsibleTeams := tfAction["responsible_teams"]
	if responsibleTeams != nil {
		resource.Properties["Octopus.Action.Manual.ResponsibleTeamIds"] = core.NewPropertyValue(responsibleTeams.(string), false)
	}

	if blockDeployments, ok := tfAction["block_deployments"]; ok {
		value := formatAsBoolOrBoundedValueForActionProperty(blockDeployments.(string))
		resource.Properties["Octopus.Action.Manual.BlockConcurrentDeployments"] = core.NewPropertyValue(value, false)
	}

	return resource
}
