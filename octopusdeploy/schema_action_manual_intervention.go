package octopusdeploy

import (
	"strconv"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenManualIntervention(actionMap map[string]interface{}, properties map[string]octopusdeploy.PropertyValue) {
	for propertyName, propertyValue := range properties {
		switch propertyName {
		case "Octopus.Action.Manual.BlockConcurrentDeployments":
			actionMap["block_deployments"], _ = strconv.ParseBool(propertyValue.Value)
		case "Octopus.Action.Manual.Instructions":
			actionMap["instructions"] = propertyValue.Value
		case "Octopus.Action.Manual.ResponsibleTeamIds":
			actionMap["responsible_teams"] = strings.Split(propertyValue.Value, ",")
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

	element.Schema["block_deployments"] = &schema.Schema{
		Type:        schema.TypeBool,
		Description: "Should other deployments be blocked while this manual intervention is awaiting action?",
		Optional:    true,
		Default:     false,
	}

	element.Schema["instructions"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "The instructions for the user to follow",
		Required:    true,
	}

	element.Schema["responsible_teams"] = &schema.Schema{
		Type:        schema.TypeList,
		Description: "The teams responsible to resolve this step. If no teams are specified, all users who have permission to deploy the project can resolve it.",
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	}

	return actionSchema
}

func expandManualInterventionAction(tfAction map[string]interface{}) *octopusdeploy.DeploymentAction {
	resource := expandAction(tfAction)
	resource.ActionType = "Octopus.Manual"

	if blockDeployments, ok := tfAction["block_deployments"]; ok {
		resource.Properties["Octopus.Action.Manual.BlockConcurrentDeployments"] = octopusdeploy.NewPropertyValue(strings.Title(strconv.FormatBool(blockDeployments.(bool))), false)
	}

	resource.Properties["Octopus.Action.Manual.Instructions"] = octopusdeploy.NewPropertyValue(tfAction["instructions"].(string), false)

	if responsibleTeams, ok := tfAction["responsible_teams"]; ok {
		resource.Properties["Octopus.Action.Manual.ResponsibleTeamIds"] = octopusdeploy.NewPropertyValue(strings.Join(getSliceFromTerraformTypeList(responsibleTeams), ","), false)
	}

	return resource
}
